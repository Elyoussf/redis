package parser

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

type Value struct {
	Typ string
	Str string
	// num   int
	Bulk  string
	Array []Value
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

func (resp *resp) Read() (Value, error) {
	_type, err := resp.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return resp.readArray()
	case BULK:
		return resp.readBulk()
	default:
		return Value{}, fmt.Errorf("undefined or unexpected type")
	}

}
func (resp *resp) readArray() (Value, error) {
	v := Value{}
	v.Typ = "array"

	length, _, err := resp.ReadInteger()
	if err != nil {
		return v, err
	}
	v.Array = make([]Value, length)

	for i := 0; i < length; i++ {

		val, err := resp.Read()

		if err != nil {
			return val, err
		}

		v.Array[i] = val
	}
	return v, nil
}

func (resp *resp) readBulk() (Value, error) {
	val := Value{}
	l, _, err := resp.ReadInteger()
	if err != nil {
		return val, err
	}
	bulk := make([]byte, l)
	resp.reader.Read(bulk)
	val.Bulk = string(bulk)
	resp.ReadLine()
	return val, nil
}

func (v *Value) Marshall() (res []byte) {

	typ := v.Typ

	switch typ {
	case "array":
		return v.marshallArray()
	case "string":
		return v.marshallString()
	case "bulk":
		return v.marshallBulk()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshallError()
	default:
		return []byte{}
	}
}

func (v *Value) marshallString() (res []byte) {
	// E=resp : +OKjkdj\r\n
	res = append(res, STRING)
	res = append(res, v.Str...) // this notation used with variadic functions that accepts mulktiple variable
	// This notation sends you string as series of bytes (Character by character)
	res = append(res, '\r')
	res = append(res, '\n')
	return res
}

func (v *Value) marshallArray() (res []byte) {
	length := len(v.Array)
	res = append(res, ARRAY)
	res = append(res, strconv.Itoa(length)...)
	res = append(res, '\r')
	res = append(res, '\n')

	for i := range length {
		e_byte := v.Array[i].Marshall()
		res = append(res, e_byte...)
	}
	return res
}

func (v *Value) marshallBulk() (res []byte) {

	length := len(v.Bulk)
	res = append(res, BULK)
	res = append(res, strconv.Itoa(length)...)
	res = append(res, "\r\n"...)
	res = append(res, []byte(v.Bulk)...)
	return res
}

func (v *Value) marshallNull() (res []byte) {
	return append(res, "$-1\r\n"...)
}

func (v *Value) marshallError() (res []byte) {
	res = append(res, ERROR)
	res = append(res, []byte(v.Str)...)
	res = append(res, "\r\n"...)
	return res
}
