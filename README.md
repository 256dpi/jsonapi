# jsonapi

[![Build Status](https://travis-ci.org/gonfire/jsonapi.svg?branch=master)](https://travis-ci.org/gonfire/jsonapi)
[![Coverage Status](https://coveralls.io/repos/github/gonfire/jsonapi/badge.svg?branch=master)](https://coveralls.io/github/gonfire/jsonapi?branch=master)
[![GoDoc](https://godoc.org/github.com/gonfire/jsonapi?status.svg)](http://godoc.org/github.com/gonfire/jsonapi)
[![Release](https://img.shields.io/github/release/gonfire/jsonapi.svg)](https://github.com/gonfire/jsonapi/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/gonfire/jsonapi)](http://goreportcard.com/report/gonfire/jsonapi)

**An extensible JSON API implementation for Go.**

Package [`jsonapi`](http://godoc.org/github.com/gonfire/jsonapi) provides structures and methods to implement JSON API compatible APIs. The library can be used with any framework or http library and ships with a built-in bridging for the native http standard library. The [`adapter`](http://godoc.org/github.com/gonfire/jsonapi/adapter) package provides additional bridging for the [echo](http://github.com/labstack/echo) framework.

# Usage

The following examples show the usage of this package:

- The [native](https://github.com/gonfire/jsonapi/blob/master/examples/native/main.go) example implements a basic API using the standard HTTP package.
- The [echo](https://github.com/gonfire/jsonapi/blob/master/examples/echo/main.go) example implements a basic API using the echo framework.
- The [client](https://github.com/gonfire/jsonapi/blob/master/examples/client/main.go) example uses the client to query the example APIs.

# Installation

Get the package using the go tool:

```bash
$ go get -u github.com/gonfire/jsonapi
```

# License

The MIT License (MIT)

Copyright (c) 2016 Joël Gähwiler
