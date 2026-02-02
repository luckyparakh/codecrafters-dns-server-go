// https://github.com/EmilHernvall/dnsguide/blob/b52da3b32b27c81e5c6729ac14fe01fef8b1b593/chapter1.md
// https://en.wikipedia.org/wiki/Domain_Name_System#DNS_message_format
package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Message struct {
	Header
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		h := Header{
			ID:      binary.BigEndian.Uint16(buf[0:2]),
			QR:      true,
			OC:      0,
			AA:      false,
			TC:      false,
			RD:      false,
			RA:      false,
			Z:       0,
			RC:      0,
			QDCount: 0,
			ANCount: 0,
			NSCount: 0,
			ARCount: 0,
		}
		encodedHeader := h.Encode()

		// Create an empty response
		response := []byte{}
		response = append(response, encodedHeader...)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
