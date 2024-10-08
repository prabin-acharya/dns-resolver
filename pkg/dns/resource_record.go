package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// ResourceRecord represents a DNS resource record
type ResourceRecord struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

func ResourceRecordFromBytes(data []byte, messageBuf *bytes.Buffer) (*ResourceRecord, error) {
	if len(data) < 10 {
		return nil, fmt.Errorf("insufficient data for resource record")
	}

	// Read and decode the name
	buf := bytes.NewBuffer(data)
	name := appendFromBufferUntilNull(buf)
	decodedName, err := DecodeName(string(name), messageBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the name: %v", err)
	}

	nameLength := len(name) - 1
	remaining := data[nameLength:]

	// Ensure we have enough data for the fixed-length fields
	if len(remaining) < 10 {
		return nil, fmt.Errorf("insufficient data for resource record fields")
	}

	// Read fixed-length fields
	typ := binary.BigEndian.Uint16(remaining[0:2])
	class := binary.BigEndian.Uint16(remaining[2:4])
	ttl := binary.BigEndian.Uint32(remaining[4:8])
	rdLength := binary.BigEndian.Uint16(remaining[8:10])

	// Ensure we have enough data for RDATA
	if len(remaining) < 10+int(rdLength) {
		return nil, fmt.Errorf("insufficient data for RDATA")
	}

	// Extract RDATA
	rData := remaining[10 : 10+int(rdLength)]

	return &ResourceRecord{
		Name:     decodedName,
		Type:     typ,
		Class:    class,
		TTL:      ttl,
		RDLength: rdLength,
		RData:    rData,
	}, nil
}

func appendFromBufferUntilNull(buf *bytes.Buffer) []byte {
	data := make([]byte, 0)
	for {
		b := buf.Next(1)
		data = append(data, b[0])
		if b[0] == 0 {
			break
		}
	}
	return data
}

// String formats a ResourceRecord into a human-readable string.
func (rr *ResourceRecord) String() string {
	// Map the Type to a human-readable string
	var recordType string
	switch rr.Type {
	case 1:
		recordType = "A"
	case 2:
		recordType = "NS"
	case 5:
		recordType = "CNAME"
	case 6:
		recordType = "SOA"
	case 15:
		recordType = "MX"
	case 28:
		recordType = "AAAA"
	default:
		recordType = fmt.Sprintf("Unknown Type (%d)", rr.Type)
	}

	// Convert RData to a human-readable form
	rDataStr := rr.RDataString()

	return fmt.Sprintf("Name: %s\nType: %s\nClass: %d\nTTL: %d\nData: %s\n",
		rr.Name, recordType, rr.Class, rr.TTL, rDataStr)
}

// RDataString converts RData into a human-readable string depending on the record type.
func (rr *ResourceRecord) RDataString() string {
	switch rr.Type {
	case 1: // A record (IPv4 address)
		ip := net.IP(rr.RData)
		return ip.String()
	case 28: // AAAA record (IPv6 address)
		ip := net.IP(rr.RData)
		return ip.String()
	case 5: // CNAME record
		return rr.Name
	default:
		return fmt.Sprintf("%x", rr.RData)
	}
}
