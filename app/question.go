package main

import (
	"encoding/binary"
	"fmt"
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
		v = strings.TrimSpace(v)
		data := fmt.Sprintf("%x%s", len(v), v)
		buf = append(buf, []byte(data)...)
	}
	buf = append(buf, 0) // Null byte to terminate the domain name
	binary.BigEndian.PutUint16(buf[len(buf):len(buf)+2], q.Type)
	binary.BigEndian.PutUint16(buf[len(buf):len(buf)+2], q.Class)
	return buf
}
