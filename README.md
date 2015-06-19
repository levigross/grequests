# GRequests
A Go "clone" of the great and famous Requests library

[![Build Status](https://travis-ci.org/levigross/grequests.svg?branch=master)](https://travis-ci.org/levigross/grequests) [![GoDoc](https://godoc.org/github.com/levigross/grequests?status.svg)](https://godoc.org/github.com/levigross/grequests) [![Coverage Status](https://coveralls.io/repos/levigross/grequests/badge.svg)](https://coveralls.io/r/levigross/grequests)

License
======

GRequests is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text

Features
========

- Every request runs in it's own goroutine
- Responses are asynchronous by nature (and you can easily make them synchronous)
- Doesn't depend on external libraries (functionality is designed to compliment `net/http`)
- Works with every version of Go from 1.3
- Responses can be serialized into JSON and XML
- Easy file uploads
- Easy file downloads
- Support for the following HTTP verbs `GET, HEAD, POST, PUT, DELETE, PATCH, OPTIONS`

Basic Example
=========
Basic GET request:

```go
resp := <-Get("http://httpbin.org/get", nil) // You can modify the request by passing an optional RequestOptions struct

fmt.Println(resp.String())
// {
//   "args": {},
//   "headers": {
//     "Accept": "*/*",
//     "Host": "httpbin.org",
```

As you can see – every request returns a channel
