# codec-http1

HTTP/1.x codecs, object bridge, cookies, multipart parsing, content coding, and upgrade helpers for Gnalloy.

This repository is part of the Gnalloy modular networking stack. The default branch is `dev`; no release tag is created during bootstrap.

## Install

```bash
go get gnalloy.org/codec-http1@dev
```

## Module Boundary

- Module path: `gnalloy.org/codec-http1`
- Responsibility: HTTP/1.x codecs, object bridge, cookies, multipart parsing, content coding, and upgrade helpers for Gnalloy
- Core dependency: `gnalloy.org/gnalloy` when this module uses Gnalloy buffers, channels, event loops, or bootstrap contracts.

## Gnalloy Dependencies

- `gnalloy.org/gnalloy`
- `gnalloy.org/codec-compression`

## Development

```bash
go test ./... -count=1
go vet ./...
go test ./... -run '^$' -bench . -benchmem -benchtime=100ms -count=1
```

For multi-repository development, use the workspace at `G:\opensource\gnalloy\go.work`. For standalone verification, set `GOWORK=off`.

## License

Apache-2.0.
