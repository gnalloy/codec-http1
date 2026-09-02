package http1

import (
	"unsafe"

	"gnalloy.org/gnalloy/buffer"
)

// retainedHeaderString 为连续请求头建立只读字符串视图，并通过 owner 约束底层内存生命周期。
// 多组件请求头使用复制回退，避免把分散存储伪装成连续内存。
func retainedHeaderString(in *buffer.CompositeByteBuf, index int, length int) (string, buffer.ByteBuf, error) {
	data, contiguous := in.ReadableSpan(index, length)
	if !contiguous {
		header, err := stringSlice(in, index, length)
		return header, nil, err
	}
	owner, err := in.Slice(index, length)
	if err != nil {
		return "", nil, err
	}
	return unsafe.String(unsafe.SliceData(data), len(data)), owner, nil
}
