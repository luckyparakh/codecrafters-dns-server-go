package main

import (
	"encoding/binary"
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
		// Length of content label 
		// e.g., for "codecrafters", length is 12
		// We append this length as a single byte
		buf = append(buf, byte(len(v)))

		// Go allows spreading a string into a byte slice
		buf = append(buf, v...)
	}
	buf = append(buf, 0) // Null byte to terminate the domain name

	typeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBuf, q.Type)
	buf = append(buf, typeBuf...)

	binary.BigEndian.PutUint16(typeBuf, q.Class)
	buf = append(buf, typeBuf...)
	return buf
}
