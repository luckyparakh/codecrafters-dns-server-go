package main

import "encoding/binary"

type Header struct {
	ID      uint16 // 16 bit
	QR      bool   // 1 bit
	OC      uint8  // 4 bit
	AA      bool   // 1 bit
	TC      bool   // 1 bit
	RD      bool   // 1 bit
	RA      bool   // 1 bit
	Z       uint8  // 3 bit
	RC      uint8  // 4 bit
	QDCount uint16 // 16 bit
	ANCount uint16 // 16 bit
	NSCount uint16 // 16 bit
	ARCount uint16 // 16 bit
}

func (h *Header) Encode() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:2], h.ID)

	// All flags QR, OC, AA, TC, RD, RA, Z, RC takes 16 bits
	var flags uint16

	if h.QR {
		// Set the highest bit (15th bit) for QR
		flags |= 1 << 15
	}

	// Set the next 4 bits (11th to 14th bits) for OC
	// First mask OC to ensure it's only 4 bits (00001111 -> 0x0F)
	// Then convert to uint16 and shift left by 11 bits
	// Finally, use bitwise OR to set these bits in flags
	flags |= uint16(h.OC&0x0F) << 11
	if h.AA {
		flags |= 1 << 10
	}
	if h.TC {
		flags |= 1 << 9
	}
	if h.RD {
		flags |= 1 << 8
	}
	if h.RA {
		flags |= 1 << 7
	}

	// Set the next 3 bits (4th to 6th bits) for Z
	// Mask Z to ensure it's only 3 bits (00000111 -> 0x07)
	// Then convert to uint16 and shift left by 4 bits
	// Finally, use bitwise OR to set these bits in flags
	flags |= uint16(h.Z&0x07) << 4

	// Set the last 4 bits (0th to 3rd bits) for RC
	flags |= uint16(h.RC & 0x0F)

	// Write the flags to the buffer
	binary.BigEndian.PutUint16(buf[2:4], flags)

	// Write the counts to the buffer
	binary.BigEndian.PutUint16(buf[4:6], h.QDCount)
	binary.BigEndian.PutUint16(buf[6:8], h.ANCount)
	binary.BigEndian.PutUint16(buf[8:10], h.NSCount)
	binary.BigEndian.PutUint16(buf[10:12], h.ARCount)

	return buf
}