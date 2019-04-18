package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s HOST:PORT\n", os.Args[0])
		return
	}

	// Open the connection with master.
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	_ = conn

}
