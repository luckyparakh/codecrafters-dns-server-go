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

func parseName(data []byte, offset int) (string, int, error) {
	if offset > len(data) {
		return "", 0, fmt.Errorf("Offset out of range")
	}

	domain := make([]byte, 0, 64)
	curOffset := offset

	for {
		if curOffset >= len(data) {
			return "", 0, fmt.Errorf("name parse: truncated")
		}

		curByte := data[curOffset]

		if curByte == 0x00 {
			curOffset++
			break
		}

		nameLen := int(curByte)
		curOffset++

		if curOffset+nameLen > len(data) {
			return "", 0, fmt.Errorf("name parse: truncated label")
		}

		if len(domain) > 0 {
			domain = append(domain, '.')
		}
		domain = append(domain, data[curOffset:curOffset+nameLen]...)
		curOffset += nameLen
	}

	return string(domain), curOffset, nil
}

func ParseQuestion1(data []byte) Question {
	domain := make([]byte, 0, 64)
	curOffset := 0
	for {
		if curOffset >= len(data) {
			break
		}
		curByte := data[curOffset]
		if curByte == 0x00 {
			curOffset++
			break
		}
		nameLen := int(curByte)
		curOffset++
		if curOffset+nameLen > len(data) {
			break
		}
		if len(domain) > 0 {
			domain = append(domain, '.')
		}
		domain = append(domain, data[curOffset:curOffset+nameLen]...)
		curOffset += nameLen
	}

	if curOffset+4 > len(data) {
		// Not enough data for Type and Class
		return Question{}
	}

	qType := binary.BigEndian.Uint16(data[curOffset : curOffset+2])
	curOffset += 2
	qClass := binary.BigEndian.Uint16(data[curOffset : curOffset+2])
	return Question{
		DomainName: string(domain),
		Type:       qType,
		Class:      qClass,
	}
}

func ParseQuestion(data []byte) Question {
	curr := 0
	var sb strings.Builder
	for {
		// Safety check to avoid out-of-bounds slice
		if curr >= len(data) {
			break
		}

		currByte := data[curr]
		if currByte == 0x00 {
			curr++ // Move past the null byte
			break
		}

		lengthOfLabel := int(data[curr])
		curr++

		// End of label index
		eol := curr + lengthOfLabel

		// Safety check to avoid out-of-bounds slice
		if eol > len(data) {
			fmt.Println("Invalid label length, out of bounds")
			break
		}

		// Get the label and append to the domain name
		label := data[curr:eol]
		if len(label) > 0 {
			sb.Write(label)
			sb.WriteByte('.')
		}
		curr += lengthOfLabel
	}

	if curr+4 > len(data) {
		// Not enough data for Type and Class
		return Question{}
	}

	qType := binary.BigEndian.Uint16(data[curr : curr+2])
	curr += 2
	qClass := binary.BigEndian.Uint16(data[curr : curr+2])
	return Question{
		DomainName: sb.String(),
		Type:       qType,
		Class:      qClass,
	}
}

// 0(4),1d,2,3,4d,5(5),6d,7,8,9,10,11d,12E
