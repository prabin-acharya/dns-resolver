package dns

import (
	"fmt"
	"net"
	"time"
)

// Client represents a DNS client.
type Client struct {
	ServerIPAddress string
	Timeout         time.Duration
}

// NewClient creates a new DNS client for the given server with a specified timeout.
func NewDNSClient(serverIPAddr string, timeout time.Duration) *Client {
	return &Client{
		ServerIPAddress: serverIPAddr,
		Timeout:         timeout,
	}
}

// SendQuery sends a DNS query and returns the response.
func (c *Client) SendQuery(msg *DNSMessage) (*DNSMessage, error) {

	// // Resolve the UDP address
	// udpAddr, err := net.ResolveUDPAddr("udp", c.ServerIPAddress)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to resolve UDP address: %v", err)
	// }

	// Establish a UDP connection to the DNS server
	conn, err := net.Dial("udp", c.ServerIPAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the DNS server: %v", err)
	}
	defer conn.Close()

	// Convert the query DNSMessage to bytes
	msgBytes := msg.ToBytes()

	// Send the DNS query
	_, err = conn.Write(msgBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to send DNS query: %v", err)
	}

	// deadline for reading the response
	conn.SetDeadline(time.Now().Add(c.Timeout))

	// Read the response
	responseBytes := make([]byte, 512) // typical DNS response size
	n, err := conn.Read(responseBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read DNS response: %v", err)
	}

	// Parse the response bytes into a DNSMessage
	response := &DNSMessage{}
	err = response.FromBytes(responseBytes[:n])
	if err != nil {
		return nil, fmt.Errorf("failed to parse DNS response: %v", err)
	}

	// TODO: check if response ID matches request ID

	return response, nil
}
