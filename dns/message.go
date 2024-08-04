// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1

package dns

import (
	"encoding/binary"
	"strings"
)

// 96 bits(12 bytes)
type Header struct {
	ID uint16 // 16 bit identification number

	// flags - 16 bit
	QR     bool  // query/response flag
	Opcode uint8 // purpose of message
	AA     bool  // authoritative answer
	TC     bool  // truncated message
	RD     bool  // recursion desired
	RA     bool  // recursion available
	Z      uint8 // reserved
	Rcode  uint8 // response code

	// 16 bit each
	QDCount uint16 // number of question entries
	ANCount uint16 // number of answer entries
	NSCount uint16 // number of authority entries
	ARCount uint16 // number of resource entries
}

// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2
// question does not have  a fixed size
type Question struct {
	QName  string
	QType  uint16 // record type like A, NS, CNAME, etc.
	QClass uint16 // a two octet code that specifies the class of the query.
	//  			For example, the QCLASS field is IN for the Internet.
}

// converts the Header to its byte representation
func (h *Header) ToBytes() []byte {
	bytes := make([]byte, 12)
	binary.BigEndian.PutUint16(bytes[0:2], h.ID)

	flags := uint16(0)
	if h.QR {
		flags |= 0x8000
	}
	flags |= uint16(h.Opcode) << 11
	if h.AA {
		flags |= 0x0400
	}
	if h.TC {
		flags |= 0x0200
	}
	if h.RD {
		flags |= 0x0100
	}
	if h.RA {
		flags |= 0x0080
	}
	flags |= uint16(h.Z) << 4
	flags |= uint16(h.Rcode)

	binary.BigEndian.PutUint16(bytes[2:4], flags)
	binary.BigEndian.PutUint16(bytes[4:6], h.QDCount)
	binary.BigEndian.PutUint16(bytes[6:8], h.ANCount)
	binary.BigEndian.PutUint16(bytes[8:10], h.NSCount)
	binary.BigEndian.PutUint16(bytes[10:12], h.ARCount)

	return bytes
}

// ResourceRecord represents a DNS resource record
type ResourceRecord struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

// DNSMessage represents a complete DNS message
type DNSMessage struct {
	Header        Header
	Questions     []Question
	Answers       []ResourceRecord
	AuthorityRRs  []ResourceRecord
	AdditionalRRs []ResourceRecord
}

// converts the Question to its byte representation
func (q *Question) ToBytes() []byte {
	qname := encodeDomainName(q.QName)
	bytes := make([]byte, len(qname)+4)
	copy(bytes, qname)
	binary.BigEndian.PutUint16(bytes[len(qname):len(qname)+2], q.QType)
	binary.BigEndian.PutUint16(bytes[len(qname)+2:], q.QClass)
	return bytes
}

// converts a domain name string to its DNS message format
func encodeDomainName(domain string) []byte {
	var encoded []byte
	labels := strings.Split(domain, ".")
	for _, label := range labels {
		encoded = append(encoded, byte(len(label)))
		encoded = append(encoded, label...)
	}
	encoded = append(encoded, 0) // root label
	return encoded
}

// 4.1.1. Header section format

//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      ID                       |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    QDCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ANCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    NSCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ARCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

// 4.1.2. Question section format

//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                                               |
//     /                     QNAME                     /
//     /                                               /
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     QTYPE                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     QCLASS                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
