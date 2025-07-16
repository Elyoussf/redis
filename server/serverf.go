package server

import (
	"fmt"
	"io"
	"log"
	"net"
	resp "redis/RESP"
)

func Server() {
	l, err := net.Listen("tcp", ":6379")

	if err != nil {
		log.Fatal("Error creating a tcp listener")
		return
	}
	fmt.Println("Server listening on port 6379 ")
	conn, err := l.Accept() // Blocking
	fmt.Println("Start accepting")
	if err != nil {
		log.Fatal("Cannot accept connections")
		return
	}

	defer conn.Close()

	for {

		resp := resp.Newresp(conn)

		value, err := resp.Read()

		if err != nil {

			if err == io.EOF {
				break
			}
			log.Fatal("Error occured while reading!")
		}
		fmt.Println(value)

		conn.Write([]byte("+OK\r\n"))
	}
}
