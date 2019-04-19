package main

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type slave struct {
	conn *net.Conn
}

type query struct {
	key string
}

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s :UI_PORT :SLAVE_PORT\n", os.Args[0])
		return
	}

	// Listen on the TCP port.
	ln, err := net.Listen("tcp", os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var slaves []slave

	// Launch goroutine to accept new connections.
	go func() {

		for {

			// Accept incoming connections.
			conn, err := ln.Accept()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			// We have a new slave! Append it to the array of slaves.
			s := &slave{
				conn: &conn,
			}
			slaves = append(slaves, *s)

			fmt.Println("new slave!")

		}

	}()

	// Start the HTTP interface.
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("master/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err)
			return
		}
		io.Copy(w, file)
		file.Close()
	})
	router.HandleFunc("/api/{key}", func(w http.ResponseWriter, r *http.Request) {

		// GET.

		// Hash the key to find out which slave to send it to.
		key := mux.Vars(r)["key"]
		hash := md5.New().Sum([]byte(key))
		slaveIndex := binary.LittleEndian.Uint32(hash) % uint32(len(slaves))

		// Write request to the lucky slave.
		fmt.Fprintf(*slaves[slaveIndex].conn, "get\n%s\n", key)

		// First line of the slave's response is content length.
		line, isPrefix, err := bufio.NewReader(*slaves[slaveIndex].conn).ReadLine()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if isPrefix {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, "weird")
			return
		}
		contentLength, err := strconv.Atoi(string(line))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Copy contentLength bytes from the slave connection to the response.
		const bufSize int = 4096
		buf := make([]byte, bufSize)
		for contentLength > bufSize {
			if _, err := (*slaves[slaveIndex].conn).Read(buf); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(os.Stderr, err)
				return
			}
			contentLength -= bufSize
			if _, err := w.Write(buf); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(os.Stderr, err)
				return
			}
		}
		// Read the last bit.
		buf = buf[:contentLength]
		if _, err := (*slaves[slaveIndex].conn).Read(buf); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if _, err := w.Write(buf); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, err)
			return
		}

	}).Methods("GET")
	router.HandleFunc("/api/{key}", func(w http.ResponseWriter, r *http.Request) {

		// SET.

		// Hash the key to find out which slave to send it to.
		key := mux.Vars(r)["key"]
		hash := md5.New().Sum([]byte(key))
		slaveIndex := binary.LittleEndian.Uint32(hash) % uint32(len(slaves))

		// Write request to the lucky slave.
		fmt.Fprintf(*slaves[slaveIndex].conn, "set\n%s\n%d\n", key, r.ContentLength)
		io.Copy(*slaves[slaveIndex].conn, r.Body)

	}).Methods("POST")
	log.Fatal(http.ListenAndServe(os.Args[1], router))

}
