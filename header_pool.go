package http1

import "sync"

const maxRecycledHeaderFields = 32

var decodedHeadersPool = sync.Pool{
	New: func() any {
		return make(Headers, 4)
	},
}

func acquireDecodedHeaders() Headers {
	return decodedHeadersPool.Get().(Headers)
}

func releaseDecodedHeaders(headers Headers) {
	if headers == nil {
		return
	}
	fields := len(headers)
	clear(headers)
	// 限制池化 map 的容量增长，避免异常大请求长期占用堆内存。
	if fields <= maxRecycledHeaderFields {
		decodedHeadersPool.Put(headers)
	}
}
