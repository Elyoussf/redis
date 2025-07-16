package resp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
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

func (resp *resp) ReadLine() (line []byte, n int, err error) {
	for {
		b, err := resp.reader.ReadByte()

		if err != nil {
			return []byte{}, 0, err
		}

		n += 1
		line = append(line, b)

		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}

	}
	return line[:len(line)-2], n, nil
}

func (resp *resp) ReadInteger() (x int, n int, err error) {

	line, _, err := resp.ReadLine()
	if err != nil {
		log.Fatal("ReadLine failed within ReadInteger")
		return 0, n, nil
	}
	res, err := strconv.ParseInt(string(line), 10, 64)

	if err != nil {
		log.Fatal("Failed to parse integer")
		return 0, 0, nil
	}
	return int(res), n, nil

}

func (resp *resp) Read() (value, error) {
	_type, err := resp.reader.ReadByte()

	if err != nil {
		return value{}, err
	}

	switch _type {
	case ARRAY:
		return resp.readArray()
	case BULK:
		return resp.readBulk()
	default:
		return value{}, fmt.Errorf("undefined or unexpected type")
	}

}
func (resp *resp) readArray() (value, error) {
	v := value{}
	v.typ = "array"

	length, _, err := resp.ReadInteger()
	if err != nil {
		return v, err
	}
	v.array = make([]value, length)

	for i := 0; i < length; i++ {

		val, err := resp.Read()

		if err != nil {
			return val, err
		}

		v.array[i] = val
	}
	return v, nil
}

func (resp *resp) readBulk() (value, error) {
	val := value{}
	l, _, err := resp.ReadInteger()
	if err != nil {
		return val, err
	}
	bulk := make([]byte, l)
	resp.reader.Read(bulk)

	val.bulk = string(bulk)

	resp.ReadLine()

	return val, nil
}
