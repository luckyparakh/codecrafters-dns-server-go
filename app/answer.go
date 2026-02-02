package main

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

type Answer struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    net.IPAddr
}

func (a *Answer) Encode() []byte {
	buf := make([]byte, 0)

	// Encode the domain name according to DNS format
	domain, err := encodeDomainName(a.Name)
	if err != nil {
		// If domain name is invalid, return empty buffer
		return buf
	}
	buf = append(buf, domain...)

	// Encode Type (2 bytes, big endian)
	typeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBuf, a.Type)
	buf = append(buf, typeBuf...)

	// Encode Class (2 bytes, big endian)
	binary.BigEndian.PutUint16(typeBuf, a.Class)
	buf = append(buf, typeBuf...)

	// Encode TTL (4 bytes, big endian)
	ttlBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(ttlBuf, a.TTL)
	buf = append(buf, typeBuf...)

	// Encode RDLength (2 bytes, big endian)
	binary.BigEndian.PutUint16(typeBuf, a.RDLength)
	buf = append(buf, typeBuf...)

	// Add RData
	for i := range strings.SplitSeq(a.RData.IP.String(), ".") {
		ii, err := strconv.Atoi(i)
		if err != nil {
			return []byte{}
		}
		buf = append(buf, byte(ii))
	}
	return buf
}
