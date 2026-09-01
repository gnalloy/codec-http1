# Testing and Performance

[简体中文](testing.zh-CN.md) | [Docs Index](README.md)

## Required Checks

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
GOWORK=off GOTOOLCHAIN=local go vet ./...
gofmt -l .
git diff --check
```

## Focused Behavior Checks

Run focused tests while working on a small behavior change:

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run 'TestName' -count=1
```

## Discovered Test Entry Points

This inventory is generated from the current `_test.go` files in this repository. It is intentionally complete so documentation review can catch stale test, benchmark, fuzz, and example coverage when code changes.

Total discovered entry points: 50.

### Tests (40)
- `TestAppendQueryStringEncodesParamsInCallOrder`
- `TestChunkedBodyEncoderWritesStreamingChunks`
- `TestContentCompressorCompressesAcceptedResponse`
- `TestContentDecompressorDecodesGzipResponse`
- `TestContentEncodingInputStreamsHTTP1GzipBody`
- `TestContinueHandlerWritesInterimResponseAndPropagatesRequest`
- `TestDecodeCookieHeaderParsesPairs`
- `TestDecodeQueryStringEnforcesParamLimit`
- `TestDecodeQueryStringPreservesOrderAndRepeatedKeys`
- `TestDecodeRejectsBodyLimits`
- `TestDecodeRejectsHeaderLimit`
- `TestDecodeRejectsPartLimit`
- `TestDecodeRequestUsesHTTPContentType`
- `TestDecodeSetCookieParsesAttributes`
- `TestDecodeSetCookieRejectsInvalidSameSite`
- `TestEncodeCookieHeaderUsesStableOrder`
- `TestEncoderDecoderRoundTrip`
- `TestEncodeRejectsInvalidName`
- `TestEncodeSetCookieWritesAttributes`
- `TestHTTPObjectAggregatorBuildsFullRequest`
- `TestHTTPObjectAggregatorBuildsFullResponse`
- `TestHTTPObjectAggregatorRejectsOversizedSingleLastContent`
- `TestParseBoundary`
- `TestQueryStringValuesIsIndependent`
- `TestRequestDecoderWithBody`
- `TestRequestDecoderWithChunkedBody`
- `TestRequestEncoder`
- `TestRequestObjectDecoderStreamsChunkedContentAndTrailers`
- `TestRequestObjectDecoderStreamsFixedContent`
- `TestResponseDecoderWithBody`
- `TestResponseDecoderWithChunkedBody`
- `TestResponseEncoder`
- `TestResponseEncoderCoalescesSmallBodyWhenEnabled`
- `TestResponseEncoderDoesNotCoalesceLargeBody`
- `TestResponseEncoderWithChunkedBody`
- `TestResponseObjectDecoderStreamsChunkedContentAndTrailers`
- `TestResponseObjectDecoderStreamsFixedContent`
- `TestStreamDecoderDrainsUnreadPart`
- `TestUpgradeHelpersBuildClientRequestAndServerResponse`
- `TestUpgradeHelpersKeepExistingConnectionTokens`

### Benchmarks (6)
- `BenchmarkFindHeaderEndFragmented`
- `BenchmarkRequestDecoderChunkedFragmentedBody`
- `BenchmarkRequestDecoderFragmentedHeader`
- `BenchmarkRequestEncoderHeader`
- `BenchmarkResponseEncoderHeader`
- `BenchmarkStringSliceFragmented`

### Fuzz Targets (4)
- `FuzzHTTP1ObjectRequestDecoder`
- `FuzzHTTP1ObjectResponseDecoder`
- `FuzzHTTP1RequestDecoder`
- `FuzzHTTP1ResponseDecoder`

### Examples (0)
- No Example functions are currently declared.

## Race Checks

```bash
GOWORK=off GOTOOLCHAIN=local go test -race ./... -count=1
```

Race checks are most valuable for core, transport, handler, resolver, observability, examples, and benchmark modules. Platform-specific transports may require native host capabilities.

## Benchmarks

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -run '^$' -bench . -benchmem -benchtime=1s -count=5
```

Report `ns/op`, `B/op`, `allocs/op`, throughput, and p99 latency separately. Include host and OS details with every result.

## Pressure Testing

Pressure tests should run against a realistic assembled stack. Use `gnalloy.org/benchmarks` for repeatable matrices and `gnalloy.org/examples` for runnable clients. Keep warmup and measurement phases separate.

## CI

The repository validation workflow runs formatting, tests, and vet on Linux, macOS, and Windows for pushes and pull requests.
