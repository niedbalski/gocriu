Golang CRIU ( http://criu.org )
=================================

This package is not official, just an experimental package
for interact with the rpc/protobuf criu service.

# How to use this library

First you need a running criu server

```bash
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
	criu, _ := gocriu.CriuClient("/tmp/criu.socket", "/tmp/pid_dump", true)

    dumped, err := criu.Dump(pid)

	if err != nil {
		panic(err)
	}

	fmt.Println(dumped)

	restored, err := criu.Restore(pid) // Restore the PID

    if err != nil {
		panic(err)
	}

	fmt.Println(restored)
}

```

# Todo

* Test
* Documentation
