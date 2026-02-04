// https://github.com/EmilHernvall/dnsguide/blob/b52da3b32b27c81e5c6729ac14fe01fef8b1b593/chapter1.md
// https://en.wikipedia.org/wiki/Domain_Name_System#DNS_message_format
package main

import (
	"fmt"
	"net"
)

type Message struct {
	Header   Header
	Question Question
	Answer   Answer
}
type config struct {
	Address  string
	Port     string
	Protocol string
}

func getConnection(c config) (*net.UDPConn, error) {
	hostPort := net.JoinHostPort(c.Address, c.Port)
	udpAddr, err := net.ResolveUDPAddr(c.Protocol, hostPort)
	if err != nil {
		return nil, err
	}
	return net.ListenUDP(c.Protocol, udpAddr)
}

func main() {
	c := config{
		Address:  "127.0.0.1",
		Port:     "2053",
		Protocol: "udp",
	}

	udpConn, err := getConnection(c)
	if err != nil {
		fmt.Println("Failed to get connection:", err)
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
		parsedMessage := parseData(buf)

		var rc uint8
		if parsedMessage.Header.OC != 0 {
			rc = 4
		}

		h := Header{
			ID:      parsedMessage.Header.ID,
			QR:      true,
			OC:      parsedMessage.Header.OC,
			AA:      false,
			TC:      false,
			RD:      parsedMessage.Header.RD,
			RA:      false,
			Z:       0,
			RC:      rc,
			QDCount: parsedMessage.Header.QDCount,
			ANCount: 1,
			NSCount: parsedMessage.Header.NSCount,
			ARCount: parsedMessage.Header.ARCount,
		}

		a := Answer{
			Name:     parsedMessage.Question.DomainName,
			Type:     1,
			Class:    1,
			TTL:      60,
			RDLength: 4,
			RData: net.IPAddr{
				IP: net.IP{8, 8, 8, 8},
			},
		}

		// Create a response
		response := []byte{}
		response = append(response, h.Encode()...)
		response = append(response, parsedMessage.Question.Encode()...)
		response = append(response, a.Encode()...)
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func parseData(data []byte) Message {
	receivedHeader := ParseHeader(data[:12])
	fmt.Printf("\nReceived Header: %+v\n", receivedHeader)

	q := ParseQuestion(data[12:])
	fmt.Printf("\nQuestion: %+v\n", q)

	return Message{
		Header:   receivedHeader,
		Question: q,
	}
}
