// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1

package dns

import "encoding/binary"

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
