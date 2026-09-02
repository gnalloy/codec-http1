package http1

import (
	"testing"

	"gnalloy.org/gnalloy/buffer"
)

func TestResponseDecoderHeaderOwnership(t *testing.T) {
	tests := []struct {
		name       string
		parts      []string
		wantRefs   []int32
		wantStatus int
		wantBody   string
	}{
		{
			name:       "contiguous header",
			parts:      []string{"HTTP/1.1 204 No Content\r\nServer: gnalloy\r\n\r\n"},
			wantRefs:   []int32{1},
			wantStatus: 204,
		},
		{
			name:       "contiguous header and body",
			parts:      []string{"HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\nok"},
			wantRefs:   []int32{2},
			wantStatus: 200,
			wantBody:   "ok",
		},
		{
			name:       "fragmented header",
			parts:      []string{"HTTP/1.1 204 No Content\r\n", "Server: gnalloy\r\n\r\n"},
			wantRefs:   []int32{0, 0},
			wantStatus: 204,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder, err := NewResponseDecoder(1024, 1024)
			if err != nil {
				t.Fatal(err)
			}
			in := buffer.NewCompositeByteBuf()
			sources := make([]buffer.ByteBuf, 0, len(tt.parts))
			for _, part := range tt.parts {
				source := testBuf([]byte(part))
				sources = append(sources, source)
				in.Append(source)
			}

			message, err := decoder.Decode(nil, in)
			if err != nil {
				in.Release()
				t.Fatal(err)
			}
			resp, ok := message.(*Response)
			if !ok {
				in.Release()
				t.Fatalf("message=%T, want *Response", message)
			}
			in.DiscardReadComponents()
			for i, source := range sources {
				if got := source.RefCnt(); got != tt.wantRefs[i] {
					resp.Release()
					in.Release()
					t.Fatalf("source[%d] refCnt=%d, want %d", i, got, tt.wantRefs[i])
				}
			}
			if resp.StatusCode != tt.wantStatus || (resp.Body != nil && string(resp.Body.Bytes()) != tt.wantBody) {
				resp.Release()
				in.Release()
				t.Fatalf("response lost retained data: %+v", resp)
			}

			resp.Release()
			for i, source := range sources {
				if got := source.RefCnt(); got != 0 {
					in.Release()
					t.Fatalf("source[%d] refCnt=%d, want 0 after response release", i, got)
				}
			}
			in.Release()
		})
	}
}
