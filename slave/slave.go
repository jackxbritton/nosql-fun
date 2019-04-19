package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
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
		fmt.Printf("Operation: %s\n", op)

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
		fmt.Printf("Key %s\n", key)

		//GET and SET routines
		if op == "get" {
			//TODO get(key)
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
			fmt.Println("GET hit")
			fmt.Printf("Length %d, Data: %s\n", strlen, ret)

		} else if op == "set" {
			//Read data length
			line, isPrefix, err = reader.ReadLine()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			if isPrefix == true {
				fmt.Println("Error: Line too long!")
			}
			var datalen = string(line)
			fmt.Printf("Data length: %s\n", datalen)

			//Read in data
			datalenint, _ := strconv.Atoi(datalen)
			data := make([]byte, datalenint)
			_, _ = reader.Read(data)

			s := string(data)
			fmt.Printf("Data: %s\n", s)
			m[key] = s

		} else {
			fmt.Println("Error: Unexpected operation. Expect \"get\" or \"set\"")
			return
		}
		time.Sleep(1 * time.Second)
	}
}

//func get() {

//}

//func set(key, data) err {

//	return
//}
