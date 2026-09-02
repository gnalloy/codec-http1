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

func releaseDecodedRequestEnvelope(request *Request) {
	if request == nil || !request.pooled {
		return
	}
	owner := request.headerOwner
	if request.recycleHeaders {
		releaseDecodedHeaders(request.Headers)
	}
	pooled := request.pooled
	*request = Request{}
	if owner != nil {
		owner.Release()
	}
	if pooled {
		decodedRequestPool.Put(request)
	}
}
