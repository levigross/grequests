# GRequests

GRequests provides a clean wrapper around Go's `net/http` package.  It mimics the convenience of the Python Requests library while keeping the power and safety of Go.

[![Join the chat at https://gitter.im/levigross/grequests](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/levigross/grequests?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Features

- Simple helpers for every HTTP verb
- Context aware request functions for easy cancellation
- RequestOptions for headers, query parameters, proxies, cookies and more
- Built in support for JSON and XML responses
- File uploads and convenient download helpers
- Session type for reusing cookies between requests

## Installation

```bash
go get github.com/levigross/grequests/v2
```

## Quick start

```go
package main

import (
    "context"
    "log"

    "github.com/levigross/grequests/v2"
)

func main() {
    resp, err := grequests.Get(context.Background(), "https://httpbin.org/get",
        grequests.UserAgent("MyAgent"))
    if err != nil {
        log.Fatal(err)
    }

    var data map[string]any
    if err := resp.JSON(&data); err != nil {
        log.Fatal(err)
    }
    log.Println(data)
}
```

### Uploading a file

```go
fd, _ := grequests.FileUploadFromDisk("testdata/file.txt")
ro := &grequests.RequestOptions{
    Files: fd,
    Data:  map[string]string{"desc": "test"},
}
resp, err := grequests.Post(context.Background(), "https://httpbin.org/post",
    grequests.FromRequestOptions(ro))
```

### Using a session

```go
sess := grequests.NewSession(nil)
_, _ = sess.Get(context.Background(), "https://httpbin.org/cookies/set?one=two", nil)
resp, _ := sess.Get(context.Background(), "https://httpbin.org/cookies", nil)
log.Println(resp.String())
```

See the documentation for a full list of available options.

## License

GRequests is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.
