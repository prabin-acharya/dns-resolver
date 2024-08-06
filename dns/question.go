// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2
// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2

package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	qname := encodeDomainName(q.Name)
	// qname := "\x06google\x03com\x00"
	bytes := make([]byte, len(qname)+4)
	copy(bytes, qname)
	binary.BigEndian.PutUint16(bytes[len(qname):len(qname)+2], q.QType)
	binary.BigEndian.PutUint16(bytes[len(qname)+2:], q.QClass)
	return bytes
}

// encodes the domain name to the DNS message format( specified in RFC 1035)

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
	q.QName, offset = decodeDomainName(data)
	q.QType = binary.BigEndian.Uint16(data[offset : offset+2])
	q.QClass = binary.BigEndian.Uint16(data[offset+2 : offset+4])
	return offset + 4, nil
}

func decodeDomainName(data []byte) (string, int) {
	var nameParts []string
	offset := 0

	for {
		// here
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

func DecodeName(qname string, messageBuf *bytes.Buffer) (string, error) {
	encoded := []byte(qname)
	var result bytes.Buffer

	for i := 0; i < len(encoded); {
		length := int(encoded[i])

		// End of domain name
		if length == 0 {
			break
		}

		// Check for compression pointer
		if length&0xC0 == 0xC0 {
			if messageBuf == nil {
				return "", fmt.Errorf("compression pointer found but no message buffer provided")
			}
			offset := int(encoded[i]&0x3F)<<8 | int(encoded[i+1])
			return decompressName(offset, messageBuf, result)
		}

		// Regular label
		i++
		if i+length > len(encoded) {
			return "", fmt.Errorf("invalid encoded domain name")
		}
		if result.Len() > 0 {
			result.WriteByte('.')
		}
		result.Write(encoded[i : i+length])
		i += length
	}

	return result.String(), nil
}

// decompressName handles decompression of a compressed domain name
func decompressName(offset int, messageBuf *bytes.Buffer, result bytes.Buffer) (string, error) {
	messageBytes := messageBuf.Bytes()[offset:]
	name := appendFromBufferUntilNull(bytes.NewBuffer(messageBytes))
	decodedName, err := DecodeName(string(name), messageBuf)
	if err != nil {
		return "", err
	}
	if result.Len() > 0 {
		result.WriteByte('.')
	}
	result.WriteString(decodedName)
	return result.String(), nil
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
