package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/prabin-acharya/dns-resolver/pkg/client"
	"github.com/prabin-acharya/dns-resolver/pkg/dns"
)

func main() {
	// Define the --raw flag
	rawFlag := flag.Bool("raw", false, "Display the raw DNS response")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Error: You need to pass a domain name as argument.")
		fmt.Println("Usage: go run main.go <domain-name>")
		return
	}

	domainName := flag.Arg(0)
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
	DNSClient := client.NewDNSClient("8.8.8.8:53", 5*time.Second)

	// Send the query and get the response
	response, err := DNSClient.SendQuery(message)
	if err != nil {
		fmt.Printf("Failed to send DNS query: %v\n", err)
		return
	}

	//  if the --raw flag is set, print the raw response
	if *rawFlag {
		fmt.Printf("DNS Response: \n %+v\n", response)
		return
	}

	// else print the formatted response
	fmt.Println("DNS Response:")
	for _, rr := range response.Answers {
		fmt.Println(rr.String())
	}
}
