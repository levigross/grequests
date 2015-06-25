# GRequests
A Go "clone" of the great and famous Requests library

[![Build Status](https://travis-ci.org/levigross/grequests.svg?branch=master)](https://travis-ci.org/levigross/grequests) [![GoDoc](https://godoc.org/github.com/levigross/grequests?status.svg)](https://godoc.org/github.com/levigross/grequests) [![Coverage Status](https://coveralls.io/repos/levigross/grequests/badge.svg)](https://coveralls.io/r/levigross/grequests)

License
======

GRequests is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text

Features
========

- Asynchronous and synchronous functionality built in
- Doesn't depend on external libraries (functionality is designed to compliment `net/http`)
- Works with every version of Go from 1.3
- Responses can be serialized into JSON and XML
- Easy file uploads
- Easy file downloads
- Support for the following HTTP verbs `GET, HEAD, POST, PUT, DELETE, PATCH, OPTIONS`

Install
=======
`go get -u github.com/levigross/grequests`


Basic Example
=========
Basic GET request:

```go
resp := Get("http://httpbin.org/get", nil) // You can modify the request by passing an optional RequestOptions struct

fmt.Println(resp.String())
// {
//   "args": {},
//   "headers": {
//     "Accept": "*/*",
//     "Host": "httpbin.org",
```
Because all of the HTTP methods return a channel, you can read the in a `select` statement as well.

```go
respChan := GetAsync("http://httpbin.org/get", nil)
	select {
	case resp := <-respChan:
		fmt.Println(resp.String())
    // {
    //   "args": {},
    //   "headers": {
    //     "Accept": "*/*",
    //     "Host": "httpbin.org",
	}

```

It is very important to check the `.Error` property of the `Response` e.g:

```go
resp := Get("http://httpbin.org/xml", nil)

if resp.Error != nil {
	log.Fatalln("Unable to make request", resp.Error)
}
```

If an error occurs all of the other properties and methods of a `Response` will be `nil`
