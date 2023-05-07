# timeago - A time formatting package

## Install

```bash
go get github.com/xeonx/timeago
```

## Docs

You can see the docs page [here](http://godoc.org/github.com/xeonx/timeago)

## Use

```go
package main

import (
 "time"

 "github.com/xeonx/timeago"
)

func main() {
 t := time.Now().Add(42 * time.Second)

 // 's' will contain "less than a minute ago"
 s := timeago.English.Format(t)

 //...
}

```

## Tests

`go test` is used for testing.
