package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type slave struct {
	connected bool
	channel   chan []byte
}

type query struct {
	key string
}

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s :PORT\n", os.Args[0])
		return
	}

	// Listen on the TCP port.
	ln, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var slaves []slave

	for {

		// Accept incoming connections.
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// We have a new slave! Append it to the array of slaves.
		s := &slave{
			connected: true,
			channel:   make(chan []byte),
		}
		slaves = append(slaves, *s)

		// With the slave added, launch a goroutine to dispatch queries.
		go func() {

			for {

				// TODO Assuming that slaves never disconnect, and no slaves join after we begin receiving queries.

				// Wait for a query and write it to the slave.
				query := <-s.channel
				conn.Write(query)

				// Keys cannot contain `=` characters.
				// As a result, if the query contained an `=` character,
				// it was just setting a value and we can continue to the next query.
				// TODO Need to include content length. Maybe key\ncontent-length\nvalue?
				if strings.Index(string(query), "=") != -1 {
					continue
				}

				// TODO Read the response.

			}

		}()

	}

}
