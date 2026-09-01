# codec-http1

[简体中文](README.zh-CN.md) | [Documentation](docs/README.md)

HTTP/1.x codecs, object bridge, cookies, multipart parsing, content coding, and upgrade helpers for Gnalloy.

This module sits above transports and below application handlers. It translates bytes or Gnalloy messages into protocol objects, and translates outbound protocol objects back to bytes. It does not open sockets or own EventLoops.

## Status

- Import path: `gnalloy.org/codec-http1`
- Repository: `github.com/gnalloy/codec-http1`
- Default branch: `dev`
- Preview install: `go get gnalloy.org/codec-http1@dev`
- License: Apache-2.0

## Install
```bash
go get gnalloy.org/codec-http1@dev
go doc gnalloy.org/codec-http1
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
```

## Documentation
- [Overview](docs/overview.md) ([中文](docs/overview.zh-CN.md))
- [Usage](docs/usage.md) ([中文](docs/usage.zh-CN.md))
- [Examples](docs/examples.md) ([中文](docs/examples.zh-CN.md))
- [Configuration](docs/configuration.md) ([中文](docs/configuration.zh-CN.md))
- [Testing and Performance](docs/testing.md) ([中文](docs/testing.zh-CN.md))
- [API Reference](docs/api.md) ([中文](docs/api.zh-CN.md))
- [Notes and Caveats](docs/notes.md) ([中文](docs/notes.zh-CN.md))
- [ADR-001 Module Boundary](docs/decisions/0001-module-boundary.md) ([中文](docs/decisions/0001-module-boundary.zh-CN.md))

## Module Boundary

This repository owns: HTTP/1.x codecs, object bridge, cookies, multipart parsing, content coding, and upgrade helpers for Gnalloy.

It does not absorb neighboring module responsibilities. Core primitives stay in `gnalloy.org/gnalloy`; protocol codecs, transports, handlers, resolvers, examples, and benchmarks stay in their own repositories.

## Packages
- `gnalloy.org/codec-http1` (`http1`)
- `gnalloy.org/codec-http1/cookie` (`cookie`)
- `gnalloy.org/codec-http1/multipart` (`multipart`)

## Gnalloy Dependencies
- `gnalloy.org/gnalloy`
- `gnalloy.org/codec-compression`

## Common Integration Pattern
- Frame, header, body, and decoded-content limits must be selected from the trusted boundary of the service.
- Streaming or chunked modes should be used for large payloads instead of materializing unbounded bodies.
- Compression modules must set decoded-size limits to defend against expansion attacks.
- ByteBuf ownership follows Gnalloy message rules: release only after the current component consumes the message.

## Current Public Entry Points

The generated API reference lists the full public surface. Common constructors or option types currently include:
- `const DefaultMaxDecodedContentBytes = 16 << 20`
- `func NewContentEncodingInput(input codec.ChunkedInput, coding ContentCoding, cfg ContentEncodingInputConfig) (*compression.CompressingChunkedInput, error)`
- `type ContentEncodingInputConfig struct{ ... }`
- `type ResponseEncoderOptions struct{ ... }`
- `var ErrInvalidCookie = errors.New("gnalloy/codec/http1/cookie: invalid cookie") ...`
- `var ErrInvalidContentType = errors.New("gnalloy/codec/http1/multipart: invalid content type") ...`

## Verification

```bash
GOWORK=off GOTOOLCHAIN=local go test ./... -count=1
GOWORK=off GOTOOLCHAIN=local go vet ./...
GOWORK=off GOTOOLCHAIN=local go test ./... -run '^$' -bench . -benchmem -count=1
```

For pressure tests, assemble this module with the relevant transport, codec, and handler stack and run the scenario from `gnalloy.org/benchmarks` or `gnalloy.org/examples`. Keep host, operating system, payload, concurrency, warmup, and repetitions in the report.

## Caveats
- This repository is intentionally narrow. Cross-module behavior should be assembled in applications, recipes, examples, or benchmark harnesses.
- Public APIs should remain Go-native and explicit; avoid runtime scanning, hidden global registries, and reflection-heavy behavior in hot paths.
- Treat network input as untrusted. Configure parser limits and return typed errors instead of panics.
- Keep benchmark claims tied to a concrete host, operating system, protocol, payload, concurrency, warmup, and repetition count.
- Codec modules do not provide a network server by themselves; combine them with a transport module and application handlers.
