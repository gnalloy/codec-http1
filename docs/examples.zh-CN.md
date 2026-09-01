# 案例

[English](examples.md) | [文档索引](README.zh-CN.md)

## 案例 1：将模块加入应用

```bash
mkdir gnalloy-app && cd gnalloy-app
go mod init example.com/gnalloy-app
go get gnalloy.org/codec-http1@dev
go doc gnalloy.org/codec-http1
```

## 案例 2：查看当前包

当前源码树暴露这些 package 导入路径：
- `gnalloy.org/codec-http1`
- `gnalloy.org/codec-http1/cookie`
- `gnalloy.org/codec-http1/multipart`

按需要的行为对对应 package 执行 `go doc`：

```bash
go doc gnalloy.org/codec-http1
go doc gnalloy.org/codec-http1/cookie
go doc gnalloy.org/codec-http1/multipart
```

精选当前导出入口：
- `const DefaultMaxDecodedContentBytes = 16 << 20`
- `func AppendQueryString(dst []byte, path string, params []QueryParam) ([]byte, error)`
- `func IsSwitchingProtocolsResponse(resp Response, protocol string) bool`
- `func IsUpgradeRequest(req Request, protocol string) bool`
- `func NewContentEncodingInput(input codec.ChunkedInput, coding ContentCoding, cfg ContentEncodingInputConfig) (*compression.CompressingChunkedInput, error)`
- `type Chunk struct{ ... }`
- `var ErrInvalidCookie = errors.New("gnalloy/codec/http1/cookie: invalid cookie") ...`
- `func AppendCookieHeader(dst []byte, cookies []Cookie) ([]byte, error)`
- `func AppendSetCookie(dst []byte, c Cookie) ([]byte, error)`
- `func EncodeCookieHeader(cookies []Cookie) (string, error)`
- `func EncodeSetCookie(c Cookie) (string, error)`
- `func SameSiteString(mode SameSite) (string, bool)`
- `var ErrInvalidContentType = errors.New("gnalloy/codec/http1/multipart: invalid content type") ...`
- `func ContentType(boundary string) (string, error)`
- `func FormDataHeader(name string, fileName string, contentType string) textproto.MIMEHeader`
- `func ParseBoundary(contentType string) (string, error)`
- `func StreamRequest(req http1.Request, limits Limits, handler PartHandler) error`
- `type Decoder struct{ ... }`

## 案例 3：将可执行测试作为行为示例

仓库测试是受支持行为的可执行示例。先从下面的精选名称开始，再阅读对应 `_test.go` 文件中的完整 setup 和断言。完整发现列表见 [测试与性能](testing.zh-CN.md)。

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run Test -count=1
```

精选当前 test、benchmark、fuzz 与 example 入口：
- `BenchmarkFindHeaderEndFragmented`
- `BenchmarkRequestDecoderChunkedFragmentedBody`
- `BenchmarkRequestDecoderFragmentedHeader`
- `BenchmarkRequestEncoderHeader`
- `BenchmarkResponseEncoderHeader`
- `BenchmarkStringSliceFragmented`
- `FuzzHTTP1ObjectRequestDecoder`
- `FuzzHTTP1ObjectResponseDecoder`
- `FuzzHTTP1RequestDecoder`
- `FuzzHTTP1ResponseDecoder`
- `TestAppendQueryStringEncodesParamsInCallOrder`
- `TestChunkedBodyEncoderWritesStreamingChunks`
- `TestContentCompressorCompressesAcceptedResponse`
- `TestContentDecompressorDecodesGzipResponse`
- `TestContentEncodingInputStreamsHTTP1GzipBody`
- `TestContinueHandlerWritesInterimResponseAndPropagatesRequest`
- `TestDecodeCookieHeaderParsesPairs`
- `TestDecodeQueryStringEnforcesParamLimit`

## 案例 4：跨模块装配

本模块的直接 Gnalloy 依赖：
- `gnalloy.org/codec-compression`
- `gnalloy.org/gnalloy`

装配说明：
- codec 位于面向字节或 datagram 的 transport 之上、应用 handler 之下。
- 它负责把字节或 Gnalloy 消息转换成协议对象，并把出站协议对象转换回字节。
- 它不打开 socket，不拥有 EventLoop，也不定义应用生命周期。

## 案例 5：压测 Harness

持续负载测试时，如果该模块参与网络流量路径，将它接入 `gnalloy.org/benchmarks` 的场景，或接入 `gnalloy.org/examples` 的可运行客户端。报告中记录 host、OS、CPU、Go version、protocol、payload、concurrency、warmup、repetitions、throughput 和 p99 latency。
