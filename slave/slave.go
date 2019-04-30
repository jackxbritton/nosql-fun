package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	//"time"
)

func main() {
	//Global for now
	m := make(map[string]string)

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
	fmt.Printf("TCP connection established with %s\n", os.Args[1])

	reader := bufio.NewReader(conn)
	//Grab first line (GET or SET)
	for true {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if isPrefix == true {
			fmt.Println("Error: Line too long!")
		}
		var op = string(line)
		//fmt.Printf("Operation: %s\n", op)

		//Grab key value
		line, isPrefix, err = reader.ReadLine()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if isPrefix == true {
			fmt.Println("Error: Line too long!")
		}

		var key = string(line)
		//fmt.Printf("Key %s\n", key)

		//GET and SET routines
		if op == "get" {
			fmt.Printf("Received GET request for KEY:%s\n", key)
			ret := m[key]
			var strlen = len(ret)
			fmt.Fprintf(conn, "%d\n", strlen)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			_, err = conn.Write([]byte(ret))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			fmt.Printf("Returning data of LENGTH:%d, VALUE:%s for KEY:%s to master\n\n", strlen, ret, key)

		} else if op == "set" {
			//Read data length
			fmt.Printf("Received SET request for KEY:%s\n", key)
			line, isPrefix, err = reader.ReadLine()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			if isPrefix == true {
				fmt.Println("Error: Line too long!")
			}
			var datalen = string(line)
			//fmt.Printf("Data length: %s\n", datalen)

			//Read in data
			datalenint, _ := strconv.Atoi(datalen)
			data := make([]byte, datalenint)
			_, _ = reader.Read(data)

			s := string(data)
			fmt.Printf("SET KEY:%s to VALUE:%s\n\n", key, s)
			m[key] = s

		} else {
			fmt.Println("Error: Unexpected operation. Expect \"get\" or \"set\"")
			return
		}
		//time.Sleep(1 * time.Second)
	}
}
