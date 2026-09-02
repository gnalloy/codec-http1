package http1

import "sync"

var decodedResponsePool sync.Pool

// AcquireResponse 获取一个空的池化响应对象。
//
// 调用方必须把对象交给 ResponseEncoder，或在放弃写出时调用 Release；此后不得访问。
func AcquireResponse() *Response {
	value := decodedResponsePool.Get()
	var response *Response
	if value == nil {
		response = new(Response)
	} else {
		response = value.(*Response)
	}
	response.pooled = true
	return response
}

func acquireDecodedResponse(response Response) *Response {
	decoded := AcquireResponse()
	*decoded = response
	decoded.pooled = true
	return decoded
}

func releaseDecodedResponseEnvelope(response *Response) {
	if response == nil || !response.pooled {
		return
	}
	owner := response.headerOwner
	if response.recycleHeaders {
		releaseDecodedHeaders(response.Headers)
	}
	*response = Response{}
	if owner != nil {
		owner.Release()
	}
	decodedResponsePool.Put(response)
}
