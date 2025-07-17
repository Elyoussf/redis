package server

import (
	"fmt"
	"io"
	"log"
	"net"
	resp "redis/RESP"
	writer "redis/Writer"
	handlers "redis/commandshandlers"
	"strings"
)

func Server() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal("Error creating a tcp listener")
		return
	}

	fmt.Println("Server listening on port 6379 ")

	for {
		conn, err := l.Accept() // Blocking
		if err != nil {
			log.Println("Cannot accept connections:", err)
			continue
		}
		fmt.Println("Accepted connection")

		go func(conn net.Conn) {
			defer conn.Close()
			handlers.LoadHandlers() // Load the handlers
			respObj := resp.Newresp(conn)
			for {
				value, err := respObj.Read()
				if err != nil {
					if err == io.EOF {
						fmt.Println("Client disconnected")
						return
					}
					fmt.Println(err)
					continue
				}
				if value.Typ != "array" {
					fmt.Println("Expected an array, got:", value.Typ)
					continue
				}
				if len(value.Array) == 0 {
					fmt.Println("Received an empty array")
					continue
				}

				handler_key := strings.ToUpper(value.Array[0].Bulk)
				handler, ok := handlers.Handlers[handler_key]
				if !ok {
					fmt.Println("Command Not found")
					continue
				}
				result := handler(value.Array[1:])

				writerObj := writer.NewWriter(conn)
				err = writerObj.Write(result)
				if err != nil {
					fmt.Println(err)
				}
			}
		}(conn)
	}
}
