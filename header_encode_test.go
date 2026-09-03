package http1

import "testing"

func TestScanHeaderFields(t *testing.T) {
	tests := []struct {
		name             string
		headers          Headers
		hasContentLength bool
		chunked          bool
	}{
		{name: "empty"},
		{
			name:             "content length canonical",
			headers:          Headers{"Content-Length": "128"},
			hasContentLength: true,
		},
		{
			name:             "content length mixed case",
			headers:          Headers{"cOnTeNt-LeNgTh": "128"},
			hasContentLength: true,
		},
		{
			name:    "chunked token",
			headers: Headers{"Transfer-Encoding": "gzip, chunked"},
			chunked: true,
		},
		{
			name:    "chunked mixed case",
			headers: Headers{"tRaNsFeR-EnCoDiNg": "gzip, ChUnKeD"},
			chunked: true,
		},
		{
			name:    "non chunked transfer encoding",
			headers: Headers{"Transfer-Encoding": "gzip"},
		},
		{
			name: "combined framing",
			headers: Headers{
				"Content-Length":    "128",
				"Transfer-Encoding": "chunked",
				"Server":            "gnalloy",
			},
			hasContentLength: true,
			chunked:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := scanHeaderFields(tt.headers)
			if metadata.hasContentLength != tt.hasContentLength {
				t.Fatalf("hasContentLength=%v, want %v", metadata.hasContentLength, tt.hasContentLength)
			}
			if metadata.chunked != tt.chunked {
				t.Fatalf("chunked=%v, want %v", metadata.chunked, tt.chunked)
			}
			if metadata.size != encodedHeaderFieldsSize(tt.headers) {
				t.Fatalf("size=%d, want %d", metadata.size, encodedHeaderFieldsSize(tt.headers))
			}
		})
	}
}

func encodedHeaderFieldsSize(headers Headers) int {
	size := 0
	for name, value := range headers {
		size += len(name) + 2 + len(value) + 2
	}
	return size
}
