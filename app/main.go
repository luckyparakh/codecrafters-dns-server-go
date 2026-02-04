// https://github.com/EmilHernvall/dnsguide/blob/b52da3b32b27c81e5c6729ac14fe01fef8b1b593/chapter1.md
// https://en.wikipedia.org/wiki/Domain_Name_System#DNS_message_format
package main

import (
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
		receivedHeader := ParseHeader(buf[:12])
		fmt.Printf("\nReceived Header %+v\n", receivedHeader)

		var rc uint8
		if receivedHeader.RC != 0 {
			rc = 4
		}
		h := Header{
			ID:      receivedHeader.ID,
			QR:      true,
			OC:      receivedHeader.OC,
			AA:      false,
			TC:      false,
			RD:      receivedHeader.RD,
			RA:      false,
			Z:       0,
			RC:      rc,
			QDCount: receivedHeader.QDCount,
			ANCount: receivedHeader.ANCount,
			NSCount: receivedHeader.NSCount,
			ARCount: receivedHeader.ARCount,
		}

		q := Question{
			DomainName: "codecrafters.io",
			Type:       1,
			Class:      1,
		}

		a := Answer{
			Name:     "codecrafters.io",
			Type:     1,
			Class:    1,
			TTL:      60,
			RDLength: 4,
			RData: net.IPAddr{
				IP: net.IP{8, 8, 8, 8},
			},
		}

		// Create an empty response
		response := []byte{}
		response = append(response, h.Encode()...)
		response = append(response, q.Encode()...)
		response = append(response, a.Encode()...)
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
