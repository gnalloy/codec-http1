# Examples

[简体中文](examples.zh-CN.md) | [Docs Index](README.md)

## Example 1: Add the Module to an Application

```bash
mkdir gnalloy-app && cd gnalloy-app
go mod init example.com/gnalloy-app
go get gnalloy.org/codec-http1@dev
go doc gnalloy.org/codec-http1
```

## Example 2: Inspect Current Packages

The current source tree exposes these package import paths:
- `gnalloy.org/codec-http1`
- `gnalloy.org/codec-http1/cookie`
- `gnalloy.org/codec-http1/multipart`

Use `go doc` against the package that matches the behavior you need:

```bash
go doc gnalloy.org/codec-http1
go doc gnalloy.org/codec-http1/cookie
go doc gnalloy.org/codec-http1/multipart
```

Selected current exported entry points:
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

## Example 3: Use Executable Tests as Behavioral Examples

Repository tests are executable examples of supported behavior. Start with the selected names below, then read the matching `_test.go` files for complete setup and assertions. See [Testing and Performance](testing.md) for the complete discovered list.

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run Test -count=1
```

Selected current test, benchmark, fuzz, and example entry points:
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

## Example 4: Cross-Module Assembly

Direct Gnalloy dependencies for this module:
- `gnalloy.org/codec-compression`
- `gnalloy.org/gnalloy`

Assembly guidance:
- Use this codec above a byte-oriented or datagram transport and below application handlers.
- The codec converts bytes or Gnalloy messages into protocol objects and converts outbound protocol objects back to bytes.
- It does not open sockets, own EventLoops, or define application lifecycle.

## Example 5: Pressure-Test Harness

For sustained load, wire this module into a scenario under `gnalloy.org/benchmarks` or a runnable client under `gnalloy.org/examples` when the module participates in network traffic. Record host, OS, CPU, Go version, protocol, payload, concurrency, warmup, repetitions, throughput, and p99 latency in the report.
