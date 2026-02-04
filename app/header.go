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

func ParseHeader(data []byte) *Header {
	h := Header{}
	h.ID = binary.BigEndian.Uint16(data[0:2])

	flags := binary.BigEndian.Uint16(data[2:4])
	h.QR = (flags & 0x8000) != 0
	// why >> 11? because OC is in bits 11 to 14,
	// so we need to shift right by 11 to get the value in the lowest 4 bits
	// so we can get OC Value by masking with 0x0F (00001111) to ensure we only get the last 4 bits
	h.OC = uint8(((flags & 0x7800) >> 11) & 0x0F)
	h.AA = (flags & 0x0400) != 0
	h.TC = (flags & 0x0200) != 0
	h.RD = (flags & 0x0100) != 0
	h.RA = (flags & 0x0080) != 0

	// Z is in bits 4 to 6, so we need to shift right by 4 to get the value in the lowest 3 bits
	// Shift right by 4 to get Z in the lowest 3 bits,
	// then mask with 0x07 (00000111) to ensure we only get the last 3 bits
	h.Z = uint8(((flags & 0x0070) >> 4) & 0x07)
	h.RC = uint8(flags & 0x000F)

	h.QDCount = binary.BigEndian.Uint16(data[4:6])
	h.ANCount = binary.BigEndian.Uint16(data[6:8])
	h.NSCount = binary.BigEndian.Uint16(data[8:10])
	h.ARCount = binary.BigEndian.Uint16(data[10:12])
	return &h
}
