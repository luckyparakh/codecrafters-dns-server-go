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
		data := fmt.Sprintf("%d%s", len(v), v)
		buf = append(buf, []byte(data)...)
	}
	binary.BigEndian.PutUint16(buf[:len(buf)+2], q.Type)
	binary.BigEndian.PutUint16(buf[:len(buf)+2], q.Class)
	return buf
}
