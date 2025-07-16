package resp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

const (
	STRING  = '+'
	ERROR   = '-'
	BULK    = '$'
	ARRAY   = '*'
	INTEGER = ':'
)

type value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []value
}

type resp struct {
	reader *bufio.Reader
}

func Newresp(rd io.Reader) *resp {
	return &resp{reader: bufio.NewReader(rd)}
}

func Parser() {
	input := "$5\r\nHAMZA\r\n"
	resp := Newresp(strings.NewReader(input))

	b, _ := resp.reader.ReadByte()

	if b != '$' {
		fmt.Println(string(b))
		// Something other than bul string
		log.Fatal("Expected bulk string starting with $")
		return
	}
	b, _ = resp.reader.ReadByte()

	size, err := strconv.ParseInt(string(b), 10, 64)

	if err != nil {
		log.Fatal(err)
		return
	}
	resp.reader.ReadByte()
	resp.reader.ReadByte()
	name := make([]byte, size)
	resp.reader.Read(name)
	fmt.Println(string(name))

}
