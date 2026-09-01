# Overview

[简体中文](overview.zh-CN.md) | [Docs Index](README.md)

## Purpose

HTTP/1.x codecs, object bridge, cookies, multipart parsing, content coding, and upgrade helpers for Gnalloy.

This module sits above transports and below application handlers. It translates bytes or Gnalloy messages into protocol objects, and translates outbound protocol objects back to bytes. It does not open sockets or own EventLoops.

## Repository Identity

- Module path: `gnalloy.org/codec-http1`
- GitHub repository: `github.com/gnalloy/codec-http1`
- Default branch: `dev`
- License: Apache-2.0

## Package Map
- `gnalloy.org/codec-http1` (`http1`)
- `gnalloy.org/codec-http1/cookie` (`cookie`)
- `gnalloy.org/codec-http1/multipart` (`multipart`)

## Direct Gnalloy Dependencies

- `gnalloy.org/codec-compression`
- `gnalloy.org/gnalloy`

## Direct Dependents in the Current Repository Set

- `gnalloy.org/codec-http2`
- `gnalloy.org/codec-http3`
- `gnalloy.org/codec-websocket`
- `gnalloy.org/examples`
- `gnalloy.org/handler-cors`
- `gnalloy.org/recipes`

## Architecture Position

Gnalloy keeps the core small and dependency-light. This repository is a replaceable module around one responsibility, connected through explicit Go packages instead of runtime discovery.

## Compatibility

The public import path is `gnalloy.org/codec-http1`. Until the first stable tag is published, use `@dev` or an explicit pseudo-version selected by your dependency policy.
