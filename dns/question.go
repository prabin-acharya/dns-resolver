// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2
// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2

package dns

import (
	"encoding/binary"
	"strings"
)

// question does not have  a fixed size
type Question struct {
	Name   string // Original domain name
	QName  string // Encoded domain name
	QType  uint16 // record type like A, NS, CNAME, etc.
	QClass uint16 // a two octet code that specifies the class of the query.
	//              For example, the QCLASS field is IN for the Internet.
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

// encodes the domain name to the DNS message format( specified in RFC 1035)

// func encodeDomainName(name string) []byte {
//     var encodedName []byte
//     parts := strings.Split(name, ".")
//     for _, part := range parts {
//         encodedName = append(encodedName, byte(len(part)))
//         encodedName = append(encodedName, part...)
//     }
//     encodedName = append(encodedName, 0) // terminating zero
//     return encodedName
// }

func encodeDomainName(name string) string {
	domainParts := strings.Split(name, ".")
	qname := ""
	for _, part := range domainParts {
		newDomainPart := string(byte(len(part))) + part
		qname += newDomainPart
	}
	return qname + "\x00"
}

func (q *Question) FromBytes(data []byte) (int, error) {
	var offset int
	q.Name, offset = decodeDomainName(data)
	q.QType = binary.BigEndian.Uint16(data[offset : offset+2])
	q.QClass = binary.BigEndian.Uint16(data[offset+2 : offset+4])
	return offset + 4, nil
}

func decodeDomainName(data []byte) (string, int) {
	var nameParts []string
	offset := 0

	for {
		length := int(data[offset])
		if length == 0 {
			offset++
			break
		}
		offset++
		nameParts = append(nameParts, string(data[offset:offset+length]))
		offset += length
	}

	return strings.Join(nameParts, "."), offset
}

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
