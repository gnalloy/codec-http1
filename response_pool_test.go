package http1

import (
	"testing"

	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
)

func TestResponseEncoderRecyclesPooledResponse(t *testing.T) {
	ch := channel.NewLocalChannel(1, buffer.NewHeapAllocator(), releasingHTTP1Sink{})
	if err := ch.Pipeline().AddLast("encoder", NewResponseEncoder()); err != nil {
		t.Fatal(err)
	}
	resp := AcquireResponse()
	resp.StatusCode = 200
	resp.Headers = Headers{"Server": "gnalloy"}
	resp.Body = buffer.NewSharedBuffer([]byte("ok"))

	if err := ch.Write(resp); err != nil {
		t.Fatal(err)
	}
	if resp.Version != "" || resp.StatusCode != 0 || resp.Reason != "" || resp.Headers != nil || resp.Body != nil {
		t.Fatalf("response retained state after encoding: %+v", *resp)
	}
}

func TestAcquireResponseRoundTripDoesNotAllocate(t *testing.T) {
	resp := AcquireResponse()
	resp.Release()

	allocs := testing.AllocsPerRun(1000, func() {
		resp := AcquireResponse()
		resp.Release()
	})
	if allocs != 0 {
		t.Fatalf("allocs=%f, want 0", allocs)
	}
}

func BenchmarkAcquireResponseRoundTrip(b *testing.B) {
	resp := AcquireResponse()
	resp.Release()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		resp := AcquireResponse()
		resp.Release()
	}
}
