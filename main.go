package main

import (
	"fmt"
	"os"
	"time"

	"github.com/prabin-acharya/dns-resolver/dns"
)

func main() {
	fmt.Println("Hello DNS####")

	if len(os.Args) < 2 {
		fmt.Println("Error: You need to pass a domain name as argument.")
		fmt.Println("Usage: go run main.go <domain-name>")
		return
	}

	domainName := os.Args[1]
	fmt.Printf("Querying DNS for domain: %s\n", domainName)

	// Create a DNS header
	header := dns.Header{
		ID:      5578,
		QR:      false, // This is a query
		Opcode:  0,     // Standard query
		AA:      false,
		TC:      false,
		RD:      true, // Recursion desired
		RA:      false,
		Z:       0,
		Rcode:   0,
		QDCount: 1, // One question
		ANCount: 0,
		NSCount: 0,
		ARCount: 0,
	}

	// Create a DNS question
	question := dns.Question{
		Name:   domainName,
		QName:  domainName,
		QType:  1, // Type A
		QClass: 1, // Class IN, IN for Internet
	}

	// Create a new DNS message with the header and question
	message := dns.NewDNSMessage(header, []dns.Question{question}, nil, nil, nil)

	// Create a DNS client
	client := dns.NewDNSClient("8.8.8.8:53", 5*time.Second)

	// Send the query and get the response
	response, err := client.SendQuery(message)
	if err != nil {
		fmt.Printf("Failed to send DNS query: %v\n", err)
		return
	}

	// Print the response
	fmt.Printf("Response: %+v\n", response)
}

// // BuildRequest creates a DNS request message
// func BuildRequest(domain string, qType uint16) *DNSMessage {
// 	header := Header{
// 		ID:      1, // Typically, a unique identifier per request
// 		QR:      false,
// 		Opcode:  0, // Standard query
// 		AA:      false,
// 		TC:      false,
// 		RD:      true,
// 		RA:      false,
// 		Z:       0,
// 		Rcode:   0,
// 		QDCount: 1,
// 		ANCount: 0,
// 		NSCount: 0,
// 		ARCount: 0,
// 	}

// 	question := Question{
// 		Name:   domain,
// 		QType:  qType,
// 		QClass: 1, // IN for the Internet
// 	}

// 	message := NewDNSMessage(header, []Question{question}, nil, nil, nil)
// 	return message
// }
