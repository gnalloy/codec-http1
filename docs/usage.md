# Usage

[简体中文](usage.zh-CN.md) | [Docs Index](README.md)

## Requirements

- Go 1.25 or newer, matching the module `go` directive.
- A Gnalloy application, recipe, example, or benchmark harness that owns lifecycle and deployment configuration.
- Standalone module verification should set `GOWORK=off` so the module is tested through its published dependency graph.

## Install
```bash
go get gnalloy.org/codec-http1@dev
```

## Import
```go
import "gnalloy.org/codec-http1"
```

## Integration Pattern
- Frame, header, body, and decoded-content limits must be selected from the trusted boundary of the service.
- Streaming or chunked modes should be used for large payloads instead of materializing unbounded bodies.
- Compression modules must set decoded-size limits to defend against expansion attacks.
- ByteBuf ownership follows Gnalloy message rules: release only after the current component consumes the message.

## API Selection

Use the API inventory to choose the exact constructor or option type for your protocol path:

```bash
go doc gnalloy.org/codec-http1
```

Common current entry points:
- `const DefaultMaxDecodedContentBytes = 16 << 20`
- `func NewContentEncodingInput(input codec.ChunkedInput, coding ContentCoding, cfg ContentEncodingInputConfig) (*compression.CompressingChunkedInput, error)`
- `type ContentEncodingInputConfig struct{ ... }`
- `type ResponseEncoderOptions struct{ ... }`
- `var ErrInvalidCookie = errors.New("gnalloy/codec/http1/cookie: invalid cookie") ...`
- `var ErrInvalidContentType = errors.New("gnalloy/codec/http1/multipart: invalid content type") ...`

## Cross-Module Assembly

When multiple Gnalloy repositories are developed together, create a local `go.work` file in your chosen workspace. Keep application-local `replace` directives out of published library modules unless the change is intentionally temporary and never committed.

## Error Handling

Network input, peer behavior, platform capability, and timeout failures must be handled as normal errors. Do not recover protocol correctness by panicking. Return or propagate the module error and close the affected Channel when ownership requires it.
