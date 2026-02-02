package main

import (
	"encoding/binary"
	"strconv"
	"strings"
)

type Question struct {
	DomainName string
	Type       uint16
	Class      uint16
}

func (q *Question) Encode() []byte {
	buf := make([]byte, 0)
	for v := range strings.SplitSeq(q.DomainName, ".") {
		buf = append(buf, []byte(strconv.Itoa(len(v)))...)
		buf = append(buf, []byte(v)...)
	}
	buf = append(buf, 0) // Null byte to terminate the domain name

	typeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBuf, q.Type)
	buf = append(buf, typeBuf...)
	
	binary.BigEndian.PutUint16(typeBuf, q.Class)
	buf = append(buf, typeBuf...)
	return buf
}
