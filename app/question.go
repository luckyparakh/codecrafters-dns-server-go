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

	// Encode the domain name according to DNS format
	domain, err := encodeDomainName(q.DomainName)
	if err != nil {
		// If domain name is invalid, return empty buffer
		return buf
	}
	buf = append(buf, domain...)

	// Encode Type (2 bytes, big endian)
	typeBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBuf, q.Type)
	buf = append(buf, typeBuf...)

	// Encode Class (2 bytes, big endian)
	classBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(classBuf, q.Class)
	buf = append(buf, classBuf...)

	return buf
}

// encodeDomainName encodes a domain name into DNS wire format.
// e.g., "www.example.com" -> [3]www[7]example[3]com[0]
func encodeDomainName(domain string) ([]byte, error) {
	if len(domain) == 0 {
		return []byte{0}, nil // root domain
	}

	labels := strings.Split(domain, ".")
	// +1 for null byte; rest GO will increase the slice as needed
	buf := make([]byte, 0, len(domain)+1)

	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return nil, fmt.Errorf("invalid label length: %q", label)
		}
		// Append length byte it will 1 byte
		buf = append(buf, byte(len(label)))
		buf = append(buf, label...)
	}
	buf = append(buf, 0) // Null byte to terminate the domain name
	return buf, nil
}
