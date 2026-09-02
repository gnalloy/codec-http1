package http1

import "testing"

func TestParseRequestHeaderWithFramingInto(t *testing.T) {
	req, framing, err := parseRequestHeaderWithFramingInto(
		"POST /upload HTTP/1.1\r\ncontent-length: 3\r\ntransfer-encoding: gzip, CHUNKED\r\n\r\n",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if req.Method != "POST" || req.URI != "/upload" || req.Version != "HTTP/1.1" {
		t.Fatalf("request=%+v", req)
	}
	if framing.contentLength != 3 || !framing.chunked {
		t.Fatalf("framing=%+v, want content length 3 and chunked", framing)
	}
}

func TestContentLengthSupportsCanonicalAndCaseInsensitiveNames(t *testing.T) {
	tests := []struct {
		name    string
		headers Headers
		want    int
	}{
		{name: "canonical", headers: Headers{"Content-Length": "12"}, want: 12},
		{name: "lowercase", headers: Headers{"content-length": "13"}, want: 13},
		{name: "missing", headers: Headers{"Content-Type": "text/plain"}, want: 0},
		{name: "empty", headers: Headers{"Content-Length": ""}, want: -1},
		{name: "invalid", headers: Headers{"Content-Length": "x"}, want: -1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := contentLength(test.headers); got != test.want {
				t.Fatalf("contentLength=%d, want %d", got, test.want)
			}
		})
	}
}
