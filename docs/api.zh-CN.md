# API 参考

[English](api.md) | [文档索引](README.zh-CN.md)

本清单由本仓库 package 的 `go doc -short` 生成，用于快速查看公共面。精确语义以源码和测试为准。

## 包

### `gnalloy.org/codec-http1`

包名：`http1`

```text
const DefaultMaxDecodedContentBytes = 16 << 20
func AppendQueryString(dst []byte, path string, params []QueryParam) ([]byte, error)
func IsSwitchingProtocolsResponse(resp Response, protocol string) bool
func IsUpgradeRequest(req Request, protocol string) bool
func NewContentEncodingInput(input codec.ChunkedInput, coding ContentCoding, cfg ContentEncodingInputConfig) (*compression.CompressingChunkedInput, error)
type Chunk struct{ ... }
type ChunkedBodyEncoder struct{}
    func NewChunkedBodyEncoder() *ChunkedBodyEncoder
type ContentCoding string
    const ContentCodingGzip ContentCoding = "gzip" ...
type ContentCompressor struct{ ... }
    func NewContentCompressor(minBytes int, codings ...ContentCoding) *ContentCompressor
type ContentDecompressor struct{ ... }
    func NewContentDecompressor(maxDecodedBytes int) *ContentDecompressor
type ContentEncodingInputConfig struct{ ... }
type ContinueHandler struct{}
    func NewContinueHandler() *ContinueHandler
type HTTPContent struct{ ... }
type HTTPObjectAggregator struct{ ... }
    func NewHTTPObjectAggregator(maxContentBytes int) *HTTPObjectAggregator
type Headers map[string]string
type LastHTTPContent struct{ ... }
type Object interface{ ... }
type QueryParam struct{ ... }
type QueryString struct{ ... }
    func DecodeQueryString(uri string, maxParams int) (QueryString, error)
type Request struct{ ... }
    func NewClientUpgradeRequest(method string, uri string, protocol string, headers Headers) Request
type RequestDecoder struct{ ... }
    func NewRequestDecoder(maxHeaderBytes int, maxBodyBytes int) (*RequestDecoder, error)
type RequestEncoder struct{}
    func NewRequestEncoder() *RequestEncoder
type RequestObjectDecoder struct{ ... }
    func NewRequestObjectDecoder(maxHeaderBytes int, maxContentBytes int) (*RequestObjectDecoder, error)
type Response struct{ ... }
    func NewSwitchingProtocolsResponse(protocol string, headers Headers) Response
type ResponseDecoder struct{ ... }
    func NewResponseDecoder(maxHeaderBytes int, maxBodyBytes int) (*ResponseDecoder, error)
type ResponseEncoder struct{ ... }
    func NewResponseEncoder() *ResponseEncoder
    func NewResponseEncoderWithOptions(options ResponseEncoderOptions) *ResponseEncoder
type ResponseEncoderOptions struct{ ... }
type ResponseObjectDecoder struct{ ... }
    func NewResponseObjectDecoder(maxHeaderBytes int, maxContentBytes int) (*ResponseObjectDecoder, error)
```

### `gnalloy.org/codec-http1/cookie`

包名：`cookie`

```text
var ErrInvalidCookie = errors.New("gnalloy/codec/http1/cookie: invalid cookie") ...
func AppendCookieHeader(dst []byte, cookies []Cookie) ([]byte, error)
func AppendSetCookie(dst []byte, c Cookie) ([]byte, error)
func EncodeCookieHeader(cookies []Cookie) (string, error)
func EncodeSetCookie(c Cookie) (string, error)
func SameSiteString(mode SameSite) (string, bool)
type Cookie struct{ ... }
    func DecodeCookieHeader(header string) ([]Cookie, error)
    func DecodeSetCookie(header string) (Cookie, error)
type SameSite uint8
    const SameSiteDefault SameSite = iota ...
```

### `gnalloy.org/codec-http1/multipart`

包名：`multipart`

```text
var ErrInvalidContentType = errors.New("gnalloy/codec/http1/multipart: invalid content type") ...
func ContentType(boundary string) (string, error)
func FormDataHeader(name string, fileName string, contentType string) textproto.MIMEHeader
func ParseBoundary(contentType string) (string, error)
func StreamRequest(req http1.Request, limits Limits, handler PartHandler) error
type Decoder struct{ ... }
    func NewDecoder(boundary string, limits Limits) (*Decoder, error)
    func NewDecoderFromContentType(contentType string, limits Limits) (*Decoder, error)
type Encoder struct{ ... }
    func NewEncoder(writer io.Writer) (*Encoder, error)
    func NewEncoderWithBoundary(writer io.Writer, boundary string) (*Encoder, error)
type Limits struct{ ... }
type Part struct{ ... }
    func DecodeRequest(req http1.Request, limits Limits) ([]Part, error)
type PartHandler interface{ ... }
type PartHandlerFunc func(info PartInfo, body io.Reader) error
type PartInfo struct{ ... }
```
