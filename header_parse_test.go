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

func TestParsedRequestContentExpectation(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want bool
	}{
		{
			name: "bodyless",
			raw:  "GET / HTTP/1.1\r\nHost: example.test\r\nExpect: 100-continue\r\n\r\n",
			want: false,
		},
		{
			name: "fixed content",
			raw:  "POST / HTTP/1.1\r\nContent-Length: 4\r\nExpect: 100-continue\r\n\r\n",
			want: true,
		},
		{
			name: "chunked content",
			raw:  "POST / HTTP/1.1\r\nTransfer-Encoding: chunked\r\nExpect: 100-continue\r\n\r\n",
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := parseRequestHeader(test.raw)
			if err != nil {
				t.Fatal(err)
			}
			if got := req.expectsContent(); got != test.want {
				t.Fatalf("expectsContent=%t, want %t", got, test.want)
			}
			if !req.ExpectsContinue() {
				t.Fatal("parsed Expect header was not preserved")
			}
		})
	}
}

func BenchmarkParsedBodylessRequestContinueCheck(b *testing.B) {
	req, err := parseRequestHeader("GET / HTTP/1.1\r\nHost: example.test\r\n\r\n")
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if req.expectsContent() && req.ExpectsContinue() {
			b.Fatal("bodyless request unexpectedly requires continue")
		}
	}
}

func BenchmarkBodylessRequestExpectHeaderScan(b *testing.B) {
	headers := Headers{"Host": "example.test"}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if headers.ContainsToken("Expect", "100-continue") {
			b.Fatal("bodyless request unexpectedly requires continue")
		}
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
