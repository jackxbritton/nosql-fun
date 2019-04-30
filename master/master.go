package main

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/mux"
)

type slave struct {
	conn *net.Conn
}

type query struct {
	key string
}

func main() {

	log.SetFlags(log.Lshortfile)

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s :UI_PORT :SLAVE_PORT\n", os.Args[0])
		return
	}

	// Listen on the TCP port.
	ln, err := net.Listen("tcp", os.Args[2])
	if err != nil {
		log.Println(err)
		return
	}

	// Array of slaves,
	// and a flag to signal when querying has begun.
	var slaves []slave
	stateFlag := int32(0)

	// Launch goroutine to accept new connections.
	go func() {

		for {

			// Accept incoming connections.
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			// If we've moved to the query state,
			// log an error and ignore the connection.
			if atomic.LoadInt32(&stateFlag) == 1 {
				conn.Close()
				log.Println("slave attempted to join after querying had begun")
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
			log.Println(err)
			return
		}
		io.Copy(w, file)
		file.Close()
	})
	router.HandleFunc("/api/{key}", func(w http.ResponseWriter, r *http.Request) {

		// GET.

		// Set the stateFlag to signal that queries have begun.
		atomic.StoreInt32(&stateFlag, 1)

		// Hash the key to find out which slave to send it to.
		key := mux.Vars(r)["key"]
		hash := md5.New().Sum([]byte(key))
		slaveIndex := binary.LittleEndian.Uint32(hash) % uint32(len(slaves))

		fmt.Printf("getting key '%s'\n", key)
		fmt.Printf("md5(%s) = %s\n", key, hex.EncodeToString(hash))
		fmt.Printf("slave index = %d\n", slaveIndex)

		// Write request to the lucky slave.
		fmt.Fprintf(*slaves[slaveIndex].conn, "get\n%s\n", key)

		// First line of the slave's response is content length.
		reader := bufio.NewReader(*slaves[slaveIndex].conn)
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if isPrefix {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("")
			return
		}
		contentLength, err := strconv.Atoi(string(line))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Copy contentLength bytes from the slave connection to the response.
		buf := make([]byte, contentLength)
		if _, err := reader.Read(buf); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		if _, err := w.Write(buf); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		fmt.Println(string(buf))

	}).Methods("GET")
	router.HandleFunc("/api/{key}", func(w http.ResponseWriter, r *http.Request) {

		// SET.

		// Set the stateFlag to signal that queries have begun.
		atomic.StoreInt32(&stateFlag, 1)

		// Hash the key to find out which slave to send it to.
		key := mux.Vars(r)["key"]
		hash := md5.New().Sum([]byte(key))
		slaveIndex := binary.LittleEndian.Uint32(hash) % uint32(len(slaves))

		fmt.Printf("setting key '%s'\n", key)
		fmt.Printf("md5(%s) = %s\n", key, string(hash))
		fmt.Printf("slave index = %d\n", slaveIndex)

		// Write request to the lucky slave.
		fmt.Fprintf(*slaves[slaveIndex].conn, "set\n%s\n%d\n", key, r.ContentLength)
		io.Copy(*slaves[slaveIndex].conn, r.Body)

	}).Methods("POST")
	log.Fatal(http.ListenAndServe(os.Args[1], router))

}
