# Configuration

[简体中文](configuration.zh-CN.md) | [Docs Index](README.md)

Configuration in Gnalloy modules is explicit. Prefer constructor arguments, option structs, and application-owned configuration files over package-level mutable state.

## Primary Configuration Points
- Frame, header, body, and decoded-content limits must be selected from the trusted boundary of the service.
- Streaming or chunked modes should be used for large payloads instead of materializing unbounded bodies.
- Compression modules must set decoded-size limits to defend against expansion attacks.
- ByteBuf ownership follows Gnalloy message rules: release only after the current component consumes the message.

## Recommended Defaults

- Start with bounded sizes and short integration-test timeouts.
- Increase limits only after measuring realistic payloads and peer behavior.
- Keep security-sensitive defaults closed or conservative.
- Document every production override in the owning service, not in this library module.

## Environment Variables

This library module does not require repository-specific environment variables for normal unit tests. Applications, examples, benchmarks, and CI jobs may define their own variables around it.

## Local Workspace Development

Use a local `go.work` file only as a developer convenience. Published module metadata should remain portable and must not contain machine-specific paths.
