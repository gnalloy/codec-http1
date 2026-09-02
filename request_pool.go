package http1

import "sync"

var decodedRequestPool sync.Pool

func acquireDecodedRequest(request Request) *Request {
	value := decodedRequestPool.Get()
	var decoded *Request
	if value == nil {
		decoded = new(Request)
	} else {
		decoded = value.(*Request)
	}
	*decoded = request
	decoded.pooled = true
	return decoded
}

func recycleDecodedRequest(request *Request) {
	if request == nil {
		return
	}
	pooled := request.pooled
	*request = Request{}
	if pooled {
		decodedRequestPool.Put(request)
	}
}

func releaseDecodedRequestEnvelope(request *Request) {
	if request == nil || !request.pooled {
		return
	}
	if request.recycleHeaders {
		releaseDecodedHeaders(request.Headers)
	}
	recycleDecodedRequest(request)
}
