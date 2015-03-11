Golang CRIU ( http://criu.org )
=================================

This package is not official, just an experimental package
for interact with the rpc/protobuf criu service.

# How to use this library

First you need a running criu server

```shell
$ sudo criu service -vvv -W criu -o service.log -b -x --address /tmp/criu.socket -j --shell-job
```

Then you can write a Go client, as an example:

```go
package main

import (
	gocriu "github.com/niedbalski/gocriu"
	"os"
        "fmt"
	"strconv"
)

func main() {
	c, _ := gocriu.NewCriu("/tmp/criu.socket", "/tmp/criogenic", true)
	dumped, err := c.Dump(pid)
	if err != nil {
		panic(err)
	}

	fmt.Println(dumped)

	restored, err := c.Restore(pid)
	if err != nil {
		panic(err)
	}

	fmt.Println(restored)
}

```

# Todo

* Test
* Documentation
