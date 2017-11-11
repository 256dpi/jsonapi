# jsonapi

[![Build Status](https://travis-ci.org/256dpi/jsonapi.svg?branch=master)](https://travis-ci.org/256dpi/jsonapi)
[![Coverage Status](https://coveralls.io/repos/github/256dpi/jsonapi/badge.svg?branch=master)](https://coveralls.io/github/256dpi/jsonapi?branch=master)
[![GoDoc](https://godoc.org/github.com/256dpi/jsonapi?status.svg)](http://godoc.org/github.com/256dpi/jsonapi)
[![Release](https://img.shields.io/github/release/256dpi/jsonapi.svg)](https://github.com/256dpi/jsonapi/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/256dpi/jsonapi)](http://goreportcard.com/report/256dpi/jsonapi)

**A fundamental and extendable JSON API library for Go.**

Package [`jsonapi`](http://godoc.org/github.com/256dpi/jsonapi) provides structures and functions to implement [JSON API](http://jsonapi.org) compatible APIs. The library can be used with any framework and is built on top of the standard Go http library.

## Examples

The following examples show the usage of this package:

- The [server](https://github.com/256dpi/jsonapi/blob/master/examples/server/main.go) example implements a basic API server using the standard HTTP package.
- The [client](https://github.com/256dpi/jsonapi/blob/master/examples/client/main.go) example uses the client to query the example API.

## Installation

Get the package using the go tool:

```bash
$ go get -u github.com/256dpi/jsonapi
```

## License

The MIT License (MIT)

Copyright (c) 2016 Joël Gähwiler
