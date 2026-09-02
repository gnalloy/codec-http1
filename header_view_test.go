package http1

import (
	"testing"

	"gnalloy.org/gnalloy/buffer"
)

func TestRetainedHeaderStringContiguousDoesNotAllocate(t *testing.T) {
	data := []byte("GET / HTTP/1.1\r\nHost: example.test\r\n\r\n")
	want := string(data)
	in := buffer.NewCompositeByteBuf()
	in.Append(testBuf(data))
	defer in.Release()

	allocs := testing.AllocsPerRun(1000, func() {
		header, owner, err := retainedHeaderString(in, in.ReaderIndex(), len(data))
		if err != nil {
			t.Fatal(err)
		}
		if header != want {
			owner.Release()
			t.Fatalf("header=%q, want %q", header, data)
		}
		owner.Release()
	})
	if allocs != 0 {
		t.Fatalf("allocs=%f, want 0", allocs)
	}
}
