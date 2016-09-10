# jsonapi

[![Build Status](https://travis-ci.org/gonfire/jsonapi.svg?branch=master)](https://travis-ci.org/gonfire/jsonapi)
[![Coverage Status](https://coveralls.io/repos/gonfire/jsonapi/badge.svg?branch=master&service=github)](https://coveralls.io/github/gonfire/jsonapi?branch=master)
[![GoDoc](https://godoc.org/github.com/gonfire/jsonapi?status.svg)](http://godoc.org/github.com/gonfire/jsonapi)
[![Release](https://img.shields.io/github/release/gonfire/jsonapi.svg)](https://github.com/gonfire/jsonapi/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/gonfire/jsonapi)](http://goreportcard.com/report/gonfire/jsonapi)

**An extensible JSON API implementation for Go.**

Package [`jsonapi`](http://godoc.org/github.com/gonfire/jsonapi) provides structures and methods to implement JSON API compatible APIs. Most methods are tailored to be used together with the echo framework, yet all of them also have a native counterpart in the [`compat`](http://godoc.org/github.com/gonfire/jsonapi/compat) sub package that allows implementing APIs using the standard HTTP library.

# Usage

The following examples show the usage of this package:

- The [native](https://github.com/gonfire/jsonapi/blob/master/examples/native/main.go) example implements a basic API using the standard HTTP package.
- The [echo](https://github.com/gonfire/jsonapi/blob/master/examples/echo/main.go) example implements a basic API using the echo framework.

# Installation

Get the package using the go tool:

```bash
$ go get -u github.com/gonfire/jsonapi
```

# License

The MIT License (MIT)

Copyright (c) 2016 Joël Gähwiler
