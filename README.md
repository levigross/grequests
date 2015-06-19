# GRequests
A Go "clone" of the great and famous Requests library

[![Build Status](https://travis-ci.org/levigross/grequests.svg?branch=master)](https://travis-ci.org/levigross/grequests) [![GoDoc](https://godoc.org/github.com/levigross/grequests?status.svg)](https://godoc.org/github.com/levigross/grequests) [![Coverage Status](https://coveralls.io/repos/levigross/grequests/badge.svg)](https://coveralls.io/r/levigross/grequests)

License
======

GRequests is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text

Basic Example
=========
```go
resp := <-Get("http://httpbin.org/get", nil)

resp.String()
// {
//   "args": {},
//   "headers": {
//     "Accept": "*/*",
//     "Host": "httpbin.org",
```






