package dns

import (
	"net"
	"time"
)

// Client represents a DNS client.
type Client struct {
	Server  string
	Timeout time.Duration
}

// NewClient creates a new DNS client for the given server with a specified timeout.
func NewClient(server string, timeout time.Duration) *Client {
	return &Client{
		Server:  server,
		Timeout: timeout,
	}
}

// SendQuery sends a DNS query and returns the response.
func (c *Client) SendQuery(query *DNSMessage) (*DNSMessage, error) {
	// Establish a UDP connection to the DNS server
	conn, err := net.Dial("udp", c.Server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Convert the query DNSMessage to bytes
	queryBytes := query.ToBytes()

	// Send the query bytes
	_, err = conn.Write(queryBytes)
	if err != nil {
		return nil, err
	}

	// Set a deadline for reading the response
	conn.SetDeadline(time.Now().Add(c.Timeout))

	// Read the response
	responseBytes := make([]byte, 512)
	n, err := conn.Read(responseBytes)
	if err != nil {
		return nil, err
	}

	// Parse the response bytes into a DNSMessage
	response := &DNSMessage{}
	err = response.FromBytes(responseBytes[:n])
	if err != nil {
		return nil, err
	}

	return response, nil
}

// // sendDNSQuery sends a DNS query to a server and returns the response.
// func sendDNSQuery(msg *DNSMessage, server string) ([]byte, error) {
// 	conn, err := net.Dial("udp", server)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer conn.Close()

// 	query := msg.ToBytes()
// 	_, err = conn.Write(query)
// 	if err != nil {
// 		return nil, err
// 	}

// 	conn.SetDeadline(time.Now().Add(5 * time.Second))
// 	response := make([]byte, 512)
// 	n, err := conn.Read(response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return response[:n], nil
// }
