package http1

import (
	"strconv"
	"strings"

	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/codec"
)

var (
	crlfBytes      = []byte{'\r', '\n'}
	headerEndBytes = []byte{'\r', '\n', '\r', '\n'}
)

type Headers map[string]string

func (h Headers) Get(name string) string {
	for k, v := range h {
		if strings.EqualFold(k, name) {
			return v
		}
	}
	return ""
}

func (h Headers) Set(name string, value string) {
	h[name] = value
}

func (h Headers) Del(name string) {
	if _, ok := h[name]; ok {
		delete(h, name)
		return
	}
	for k := range h {
		if strings.EqualFold(k, name) {
			delete(h, k)
			return
		}
	}
}

func (h Headers) ContainsToken(name string, token string) bool {
	value := h.Get(name)
	for part := range strings.SplitSeq(value, ",") {
		if strings.EqualFold(strings.TrimSpace(part), token) {
			return true
		}
	}
	return false
}

// Request 表示完整 HTTP/1 请求。
//
// RequestDecoder 产出池化指针；下游处理完成后必须调用 Release 一次，且之后不得继续访问。
type Request struct {
	Method  string
	URI     string
	Version string
	Headers Headers
	Body    buffer.ByteBuf

	recycleHeaders bool
	pooled         bool
}

func (r Request) KeepAlive() bool {
	if r.Headers.ContainsToken("Connection", "close") {
		return false
	}
	if r.Version == "HTTP/1.0" {
		return r.Headers.ContainsToken("Connection", "keep-alive")
	}
	return true
}

func (r Request) ExpectsContinue() bool {
	return r.Headers.ContainsToken("Expect", "100-continue")
}

// Release 释放正文、回收解码头和池化请求对象。
func (r *Request) Release() {
	if r == nil {
		return
	}
	if r.Body != nil {
		r.Body.Release()
	}
	if r.pooled {
		releaseDecodedRequestEnvelope(r)
		return
	}
	if r.recycleHeaders {
		releaseDecodedHeaders(r.Headers)
	}
}

type Response struct {
	Version    string
	StatusCode int
	Reason     string
	Headers    Headers
	Body       buffer.ByteBuf

	recycleHeaders bool
	pooled         bool
}

func (r Response) KeepAlive() bool {
	if r.Headers.ContainsToken("Connection", "close") {
		return false
	}
	if r.Version == "HTTP/1.0" {
		return r.Headers.ContainsToken("Connection", "keep-alive")
	}
	return true
}

// Release 释放正文、回收解码头和池化响应对象。
func (r *Response) Release() {
	if r == nil {
		return
	}
	if r.Body != nil {
		r.Body.Release()
	}
	if r.pooled {
		releaseDecodedResponseEnvelope(r)
		return
	}
	if r.recycleHeaders {
		releaseDecodedHeaders(r.Headers)
	}
}

type RequestDecoder struct {
	*codec.ByteToMessageDecoder
	maxHeaderBytes int
	maxBodyBytes   int
}

// NewRequestDecoder 创建产出 *Request 的流式请求解码器。
func NewRequestDecoder(maxHeaderBytes int, maxBodyBytes int) (*RequestDecoder, error) {
	if maxHeaderBytes <= 0 || maxBodyBytes < 0 {
		return nil, codec.ErrInvalidFrameLength
	}
	d := &RequestDecoder{maxHeaderBytes: maxHeaderBytes, maxBodyBytes: maxBodyBytes}
	d.ByteToMessageDecoder = codec.NewByteToMessageDecoder(d)
	return d, nil
}

func (d *RequestDecoder) Decode(ctx *channel.HandlerContext, in *buffer.CompositeByteBuf) (any, error) {
	headerEnd, ok := findHeaderEnd(in)
	if !ok {
		if in.ReadableBytes() > d.maxHeaderBytes {
			return nil, codec.ErrFrameTooLong
		}
		return nil, nil
	}
	reader := in.ReaderIndex()
	headerBytes := headerEnd - reader + 4
	if headerBytes > d.maxHeaderBytes {
		return nil, codec.ErrFrameTooLong
	}
	header, err := stringSlice(in, reader, headerBytes)
	if err != nil {
		return nil, err
	}
	headers := acquireDecodedHeaders()
	parsed, err := parseRequestHeaderInto(header, headers)
	if err != nil {
		releaseDecodedHeaders(headers)
		return nil, err
	}
	req := acquireDecodedRequest(parsed)
	req.recycleHeaders = true
	bodyLength := contentLength(req.Headers)
	if req.Headers.ContainsToken("Transfer-Encoding", "chunked") {
		body, total, ok, err := d.decodeChunkedBody(ctx, in, reader+headerBytes)
		if err != nil || !ok {
			req.Release()
			return nil, err
		}
		req.Body = body
		req.Headers.Del("Transfer-Encoding")
		req.Headers.Set("Content-Length", strconv.Itoa(body.ReadableBytes()))
		if err := in.SkipBytes(headerBytes + total); err != nil {
			req.Release()
			return nil, err
		}
		return req, nil
	}
	if bodyLength < 0 || bodyLength > d.maxBodyBytes {
		req.Release()
		return nil, codec.ErrFrameTooLong
	}
	total := headerBytes + bodyLength
	if in.ReadableBytes() < total {
		req.Release()
		return nil, nil
	}
	if bodyLength > 0 {
		req.Body, err = in.Slice(reader+headerBytes, bodyLength)
		if err != nil {
			req.Release()
			return nil, err
		}
	}
	if err := in.SkipBytes(total); err != nil {
		if req.Body != nil {
			req.Body.Release()
		}
		return nil, err
	}
	return req, nil
}

type ResponseDecoder struct {
	*codec.ByteToMessageDecoder
	maxHeaderBytes int
	maxBodyBytes   int
}

func NewResponseDecoder(maxHeaderBytes int, maxBodyBytes int) (*ResponseDecoder, error) {
	if maxHeaderBytes <= 0 || maxBodyBytes < 0 {
		return nil, codec.ErrInvalidFrameLength
	}
	d := &ResponseDecoder{maxHeaderBytes: maxHeaderBytes, maxBodyBytes: maxBodyBytes}
	d.ByteToMessageDecoder = codec.NewByteToMessageDecoder(d)
	return d, nil
}

func (d *ResponseDecoder) Decode(ctx *channel.HandlerContext, in *buffer.CompositeByteBuf) (any, error) {
	headerEnd, ok := findHeaderEnd(in)
	if !ok {
		if in.ReadableBytes() > d.maxHeaderBytes {
			return nil, codec.ErrFrameTooLong
		}
		return nil, nil
	}
	reader := in.ReaderIndex()
	headerBytes := headerEnd - reader + 4
	if headerBytes > d.maxHeaderBytes {
		return nil, codec.ErrFrameTooLong
	}
	header, err := stringSlice(in, reader, headerBytes)
	if err != nil {
		return nil, err
	}
	headers := acquireDecodedHeaders()
	parsed, err := parseResponseHeaderInto(header, headers)
	if err != nil {
		releaseDecodedHeaders(headers)
		return nil, err
	}
	resp := acquireDecodedResponse(parsed)
	resp.recycleHeaders = true
	bodyLength := contentLength(resp.Headers)
	if resp.Headers.ContainsToken("Transfer-Encoding", "chunked") {
		body, total, ok, err := d.decodeChunkedBody(ctx, in, reader+headerBytes)
		if err != nil || !ok {
			resp.Release()
			return nil, err
		}
		resp.Body = body
		resp.Headers.Del("Transfer-Encoding")
		resp.Headers.Set("Content-Length", strconv.Itoa(body.ReadableBytes()))
		if err := in.SkipBytes(headerBytes + total); err != nil {
			resp.Release()
			return nil, err
		}
		return resp, nil
	}
	if bodyLength < 0 || bodyLength > d.maxBodyBytes {
		resp.Release()
		return nil, codec.ErrFrameTooLong
	}
	total := headerBytes + bodyLength
	if in.ReadableBytes() < total {
		resp.Release()
		return nil, nil
	}
	if bodyLength > 0 {
		resp.Body, err = in.Slice(reader+headerBytes, bodyLength)
		if err != nil {
			resp.Release()
			return nil, err
		}
	}
	if err := in.SkipBytes(total); err != nil {
		if resp.Body != nil {
			resp.Body.Release()
		}
		return nil, err
	}
	return resp, nil
}

func (d *ResponseDecoder) decodeChunkedBody(ctx *channel.HandlerContext, in *buffer.CompositeByteBuf, index int) (buffer.ByteBuf, int, bool, error) {
	return decodeChunkedBody(ctx, in, index, d.maxBodyBytes)
}

func (d *RequestDecoder) decodeChunkedBody(ctx *channel.HandlerContext, in *buffer.CompositeByteBuf, index int) (buffer.ByteBuf, int, bool, error) {
	return decodeChunkedBody(ctx, in, index, d.maxBodyBytes)
}

func decodeChunkedBody(ctx *channel.HandlerContext, in *buffer.CompositeByteBuf, index int, maxBodyBytes int) (buffer.ByteBuf, int, bool, error) {
	total := 0
	aggregate, err := ctx.Channel().Allocator().Acquire(1)
	if err != nil {
		return nil, 0, false, err
	}
	aggregate.Clear()
	for {
		lineEnd, ok := findCRLF(in, index+total)
		if !ok {
			aggregate.Release()
			return nil, 0, false, nil
		}
		line, err := stringSlice(in, index+total, lineEnd-(index+total))
		if err != nil {
			aggregate.Release()
			return nil, 0, false, err
		}
		sizeText, _, _ := strings.Cut(line, ";")
		chunkSize64, err := strconv.ParseInt(strings.TrimSpace(sizeText), 16, 64)
		if err != nil || chunkSize64 < 0 {
			aggregate.Release()
			return nil, 0, false, codec.ErrInvalidFrameLength
		}
		chunkSize := int(chunkSize64)
		total += lineEnd - (index + total) + 2
		if chunkSize == 0 {
			if in.WriterIndex() < index+total+2 {
				aggregate.Release()
				return nil, 0, false, nil
			}
			total += 2
			return aggregate, total, true, nil
		}
		if aggregate.ReadableBytes()+chunkSize > maxBodyBytes {
			aggregate.Release()
			return nil, 0, false, codec.ErrFrameTooLong
		}
		if in.WriterIndex() < index+total+chunkSize+2 {
			aggregate.Release()
			return nil, 0, false, nil
		}
		part, err := in.Slice(index+total, chunkSize)
		if err != nil {
			aggregate.Release()
			return nil, 0, false, err
		}
		if aggregate.WritableBytes() < chunkSize {
			next, err := ctx.Channel().Allocator().Acquire(aggregate.ReadableBytes() + chunkSize)
			if err != nil {
				part.Release()
				aggregate.Release()
				return nil, 0, false, err
			}
			if err := buffer.WriteReadableBytes(next, aggregate); err != nil {
				next.Release()
				part.Release()
				aggregate.Release()
				return nil, 0, false, err
			}
			aggregate.Release()
			aggregate = next
		}
		if err := buffer.WriteReadableBytes(aggregate, part); err != nil {
			part.Release()
			aggregate.Release()
			return nil, 0, false, err
		}
		part.Release()
		total += chunkSize
		cr, _ := in.GetByte(index + total)
		lf, _ := in.GetByte(index + total + 1)
		if cr != '\r' || lf != '\n' {
			aggregate.Release()
			return nil, 0, false, codec.ErrInvalidFrameLength
		}
		total += 2
	}
}

type RequestEncoder struct{}

func NewRequestEncoder() *RequestEncoder {
	return &RequestEncoder{}
}

func (e *RequestEncoder) Write(ctx *channel.HandlerContext, msg any) error {
	var req Request
	var pooled *Request
	switch value := msg.(type) {
	case Request:
		req = value
	case *Request:
		if value == nil {
			return ctx.Write(msg)
		}
		req = *value
		if value.pooled {
			pooled = value
		}
	default:
		return ctx.Write(msg)
	}
	if pooled != nil {
		defer releaseDecodedRequestEnvelope(pooled)
	}
	chunked := requestChunked(req)
	out, err := encodeRequestHead(ctx, req, chunked)
	if err != nil {
		if req.Body != nil {
			req.Body.Release()
		}
		return err
	}
	if err := ctx.Write(out); err != nil {
		out.Release()
		if req.Body != nil {
			req.Body.Release()
		}
		return err
	}
	if req.Body != nil {
		if chunked {
			if err := writeChunkedData(ctx, req.Body); err != nil {
				return err
			}
			return writeLastChunk(ctx, nil)
		}
		return codec.WriteOutboundBuffer(ctx, req.Body)
	}
	return nil
}

// ResponseEncoderOptions 描述 HTTP/1 响应编码器的可选热路径策略。
type ResponseEncoderOptions struct {
	// CoalesceBodyBytes 大于 0 时，小于等于该阈值的非 chunked 响应会合并头部和正文。
	//
	// 该选项适合 TLS 下减少 record 数和 goroutine 边界往返；明文高吞吐场景可保持
	// 默认 0，继续依赖 writev/批量写出避免额外拷贝。
	CoalesceBodyBytes int
}

type ResponseEncoder struct {
	options ResponseEncoderOptions
}

func NewResponseEncoder() *ResponseEncoder {
	return &ResponseEncoder{}
}

func NewResponseEncoderWithOptions(options ResponseEncoderOptions) *ResponseEncoder {
	if options.CoalesceBodyBytes < 0 {
		options.CoalesceBodyBytes = 0
	}
	return &ResponseEncoder{options: options}
}

func (e *ResponseEncoder) Write(ctx *channel.HandlerContext, msg any) error {
	var resp Response
	var pooled *Response
	switch value := msg.(type) {
	case Response:
		resp = value
	case *Response:
		if value == nil {
			return ctx.Write(msg)
		}
		resp = *value
		if value.pooled {
			pooled = value
		}
	default:
		return ctx.Write(msg)
	}
	if pooled != nil {
		defer releaseDecodedResponseEnvelope(pooled)
	}
	chunked := responseChunked(resp)
	if resp.Body != nil && !chunked && e.options.CoalesceBodyBytes > 0 {
		return e.writeCoalesced(ctx, resp)
	}
	return e.writeSplit(ctx, resp, chunked)
}

func (e *ResponseEncoder) writeCoalesced(ctx *channel.HandlerContext, resp Response) error {
	bodyBytes := resp.Body.ReadableBytes()
	if bodyBytes == 0 || bodyBytes > e.options.CoalesceBodyBytes {
		return e.writeSplit(ctx, resp, false)
	}
	out, err := encodeResponse(ctx, resp, bodyBytes)
	if err != nil {
		resp.Body.Release()
		return err
	}
	resp.Body.Release()
	return codec.WriteOutboundBuffer(ctx, out)
}

func (e *ResponseEncoder) writeSplit(ctx *channel.HandlerContext, resp Response, chunked bool) error {
	out, err := encodeResponseHead(ctx, resp, chunked)
	if err != nil {
		if resp.Body != nil {
			resp.Body.Release()
		}
		return err
	}
	if err := ctx.Write(out); err != nil {
		out.Release()
		if resp.Body != nil {
			resp.Body.Release()
		}
		return err
	}
	if resp.Body != nil {
		if chunked {
			if err := writeChunkedData(ctx, resp.Body); err != nil {
				return err
			}
			return writeLastChunk(ctx, nil)
		}
		return codec.WriteOutboundBuffer(ctx, resp.Body)
	}
	return nil
}

type ContinueHandler struct{}

func NewContinueHandler() *ContinueHandler {
	return &ContinueHandler{}
}

func (h *ContinueHandler) ChannelRead(ctx *channel.HandlerContext, msg any) {
	var req *Request
	switch value := msg.(type) {
	case Request:
		req = &value
	case *Request:
		req = value
	}
	if req != nil && req.ExpectsContinue() {
		if err := ctx.Channel().WriteAndFlush(Response{StatusCode: 100}); err != nil {
			req.Release()
			ctx.FireExceptionCaught(err)
			return
		}
	}
	ctx.FireChannelRead(msg)
}

type Chunk struct {
	Data     buffer.ByteBuf
	Last     bool
	Trailers Headers
}

func (c Chunk) Release() {
	if c.Data != nil {
		c.Data.Release()
	}
}

type ChunkedBodyEncoder struct{}

func NewChunkedBodyEncoder() *ChunkedBodyEncoder {
	return &ChunkedBodyEncoder{}
}

func (e *ChunkedBodyEncoder) Write(ctx *channel.HandlerContext, msg any) error {
	chunk, ok := msg.(Chunk)
	if !ok {
		return ctx.Write(msg)
	}
	if chunk.Data != nil {
		if err := writeChunkedData(ctx, chunk.Data); err != nil {
			return err
		}
		chunk.Data = nil
	}
	if chunk.Last {
		return writeLastChunk(ctx, chunk.Trailers)
	}
	return nil
}

func writeChunkedData(ctx *channel.HandlerContext, body buffer.ByteBuf) error {
	prefix := strconv.FormatInt(int64(body.ReadableBytes()), 16) + "\r\n"
	head, err := ctx.Channel().Allocator().Acquire(len(prefix))
	if err != nil {
		body.Release()
		return err
	}
	if _, err := head.WriteBytes([]byte(prefix)); err != nil {
		head.Release()
		body.Release()
		return err
	}
	tail, err := ctx.Channel().Allocator().Acquire(2)
	if err != nil {
		head.Release()
		body.Release()
		return err
	}
	if _, err := tail.WriteBytes([]byte("\r\n")); err != nil {
		head.Release()
		tail.Release()
		body.Release()
		return err
	}
	if err := ctx.Write(head); err != nil {
		head.Release()
		tail.Release()
		body.Release()
		return err
	}
	if err := ctx.Write(body); err != nil {
		body.Release()
		tail.Release()
		return err
	}
	return codec.WriteOutboundBuffer(ctx, tail)
}

func writeLastChunk(ctx *channel.HandlerContext, trailers Headers) error {
	var builder strings.Builder
	builder.WriteString("0\r\n")
	for k, v := range trailers {
		builder.WriteString(k)
		builder.WriteString(": ")
		builder.WriteString(v)
		builder.WriteString("\r\n")
	}
	builder.WriteString("\r\n")
	data := builder.String()
	out, err := ctx.Channel().Allocator().Acquire(len(data))
	if err != nil {
		return err
	}
	if _, err := out.WriteBytes([]byte(data)); err != nil {
		out.Release()
		return err
	}
	return codec.WriteOutboundBuffer(ctx, out)
}

func findCRLF(in *buffer.CompositeByteBuf, start int) (int, bool) {
	return in.Index(start, crlfBytes)
}

func findHeaderEnd(in *buffer.CompositeByteBuf) (int, bool) {
	return findHeaderEndFrom(in, in.ReaderIndex())
}

func stringSlice(in *buffer.CompositeByteBuf, index int, length int) (string, error) {
	return buffer.ReadableStringAt(in, index, length)
}

func defaultReason(code int) string {
	switch code {
	case 100:
		return "Continue"
	case 101:
		return "Switching Protocols"
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Status"
	}
}
