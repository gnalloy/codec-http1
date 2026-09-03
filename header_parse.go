package http1

import (
	"strconv"
	"strings"

	"gnalloy.org/gnalloy/codec"
)

// parseRequestHeader 解析完整 HTTP/1 请求头，不为行列表创建临时切片。
func parseRequestHeader(src string) (Request, error) {
	return parseRequestHeaderInto(src, nil)
}

func parseRequestHeaderInto(src string, headers Headers) (Request, error) {
	req, _, err := parseRequestHeaderWithFramingInto(src, headers)
	return req, err
}

func parseRequestHeaderWithFramingInto(src string, headers Headers) (Request, messageFraming, error) {
	line, next, ok := nextHeaderLine(src, 0)
	if !ok {
		return Request{}, messageFraming{}, codec.ErrInvalidFrameLength
	}
	method, uri, version, ok := splitRequestLine(line)
	if !ok {
		return Request{}, messageFraming{}, codec.ErrInvalidFrameLength
	}
	headers, framing, err := parseHeaderFieldsWithFramingInto(src, next, headers)
	if err != nil {
		return Request{}, messageFraming{}, err
	}
	return Request{
		Method:          method,
		URI:             uri,
		Version:         version,
		Headers:         headers,
		framingKnown:    true,
		contentExpected: framing.chunked || framing.contentLength > 0,
	}, framing, nil
}

// parseResponseHeader 解析完整 HTTP/1 响应头，不为行列表创建临时切片。
func parseResponseHeader(src string) (Response, error) {
	return parseResponseHeaderInto(src, nil)
}

func parseResponseHeaderInto(src string, headers Headers) (Response, error) {
	resp, _, err := parseResponseHeaderWithFramingInto(src, headers)
	return resp, err
}

func parseResponseHeaderWithFramingInto(src string, headers Headers) (Response, messageFraming, error) {
	line, next, ok := nextHeaderLine(src, 0)
	if !ok {
		return Response{}, messageFraming{}, codec.ErrInvalidFrameLength
	}
	version, statusText, reason, ok := splitResponseLine(line)
	if !ok {
		return Response{}, messageFraming{}, codec.ErrInvalidFrameLength
	}
	statusCode, err := strconv.Atoi(statusText)
	if err != nil {
		return Response{}, messageFraming{}, codec.ErrInvalidFrameLength
	}
	headers, framing, err := parseHeaderFieldsWithFramingInto(src, next, headers)
	if err != nil {
		return Response{}, messageFraming{}, err
	}
	return Response{Version: version, StatusCode: statusCode, Reason: reason, Headers: headers}, framing, nil
}

func parseTrailerHeaders(src string) (Headers, error) {
	return parseHeaderFields(src, 0)
}

func contentLength(headers Headers) int {
	value, ok := headers["Content-Length"]
	if !ok {
		for key, candidate := range headers {
			if strings.EqualFold(key, "Content-Length") {
				value = candidate
				ok = true
				break
			}
		}
		if !ok {
			return 0
		}
	}
	return parseContentLengthValue(value)
}

func parseContentLengthValue(value string) int {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return -1
	}
	return n
}

func splitRequestLine(line string) (string, string, string, bool) {
	first := strings.IndexByte(line, ' ')
	if first < 0 {
		return "", "", "", false
	}
	second := strings.IndexByte(line[first+1:], ' ')
	if second < 0 {
		return "", "", "", false
	}
	second += first + 1
	return line[:first], line[first+1 : second], line[second+1:], true
}

func splitResponseLine(line string) (string, string, string, bool) {
	first := strings.IndexByte(line, ' ')
	if first < 0 {
		return "", "", "", false
	}
	second := strings.IndexByte(line[first+1:], ' ')
	if second < 0 {
		return line[:first], line[first+1:], "", true
	}
	second += first + 1
	return line[:first], line[first+1 : second], line[second+1:], true
}

func parseHeaderFields(src string, start int) (Headers, error) {
	return parseHeaderFieldsInto(src, start, nil)
}

func parseHeaderFieldsInto(src string, start int, headers Headers) (Headers, error) {
	headers, _, err := parseHeaderFieldsWithFramingInto(src, start, headers)
	return headers, err
}

type messageFraming struct {
	contentLength          int
	chunked                bool
	canonicalContentLength bool
}

func parseHeaderFieldsWithFramingInto(src string, start int, headers Headers) (Headers, messageFraming, error) {
	if headers == nil {
		headers = make(Headers, 4)
	}
	framing := messageFraming{}
	for start < len(src) {
		line, next, ok := nextHeaderLine(src, start)
		if !ok {
			return nil, messageFraming{}, codec.ErrInvalidFrameLength
		}
		start = next
		if line == "" {
			return headers, framing, nil
		}
		colon := strings.IndexByte(line, ':')
		if colon < 0 {
			return nil, messageFraming{}, codec.ErrInvalidFrameLength
		}
		name := strings.TrimSpace(line[:colon])
		value := strings.TrimSpace(line[colon+1:])
		headers[name] = value
		framing.observe(name, value)
	}
	return headers, framing, nil
}

func (f *messageFraming) observe(name string, value string) {
	if name == "Content-Length" {
		f.contentLength = parseContentLengthValue(value)
		f.canonicalContentLength = true
	} else if !f.canonicalContentLength && strings.EqualFold(name, "Content-Length") {
		f.contentLength = parseContentLengthValue(value)
	}
	if strings.EqualFold(name, "Transfer-Encoding") && containsHeaderToken(value, "chunked") {
		f.chunked = true
	}
}

func nextHeaderLine(src string, start int) (string, int, bool) {
	if start > len(src) {
		return "", start, false
	}
	end := strings.Index(src[start:], "\r\n")
	if end < 0 {
		return "", start, false
	}
	end += start
	return src[start:end], end + 2, true
}
