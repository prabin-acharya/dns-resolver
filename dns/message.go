// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1

package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// DNSMessage represents a complete DNS message
type DNSMessage struct {
	Header        Header
	Questions     []Question
	Answers       []ResourceRecord
	AuthorityRRs  []ResourceRecord
	AdditionalRRs []ResourceRecord
}

// NewDNSMessage creates a new DNSMessage with the given header, questions, and resource records.
func NewDNSMessage(header Header, questions []Question, answers, authorityRRs, additionalRRs []ResourceRecord) *DNSMessage {
	return &DNSMessage{
		Header:        header,
		Questions:     questions,
		Answers:       make([]ResourceRecord, 0),
		AuthorityRRs:  make([]ResourceRecord, 0),
		AdditionalRRs: make([]ResourceRecord, 0),
	}
}

// ToBytes converts the entire DNSMessage to its byte representation.
func (m *DNSMessage) ToBytes() []byte {
	headerBytes := m.Header.ToBytes()
	var questionBytes []byte
	for _, q := range m.Questions {
		questionBytes = append(questionBytes, q.ToBytes()...)
	}
	return append(headerBytes, questionBytes...)
}

// FromBytes parses a byte slice into a DNSMessage.
func (m *DNSMessage) FromBytes(data []byte) error {
	// Parse Header
	m.Header.FromBytes(data[:12])
	offset := 12

	// Parse Questions
	m.Questions = make([]Question, m.Header.QDCount)
	for i := 0; i < int(m.Header.QDCount); i++ {
		q := &m.Questions[i]
		n, err := q.FromBytes(data[offset:])
		if err != nil {
			return err
		}
		offset += n
	}

	// Parse Answers
	m.Answers = make([]ResourceRecord, m.Header.ANCount)
	for i := 0; i < int(m.Header.ANCount); i++ {
		rr, err := ResourceRecordFromBytes(data[offset:], bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		m.Answers[i] = *rr
		offset += int(rr.RDLength) + len(rr.Name) + 10 // name length + fixed fields length
	}

	// Parse AuthorityRRS
	m.AuthorityRRs = make([]ResourceRecord, m.Header.NSCount)
	for i := 0; i < int(m.Header.NSCount); i++ {
		rr, err := ResourceRecordFromBytes(data[offset:], bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		m.AuthorityRRs[i] = *rr
		offset += int(rr.RDLength) + len(rr.Name) + 10
	}

	// Parse AdditionalRRs
	m.AdditionalRRs = make([]ResourceRecord, m.Header.ARCount)
	for i := 0; i < int(m.Header.ARCount); i++ {
		rr, err := ResourceRecordFromBytes(data[offset:], bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		m.AdditionalRRs[i] = *rr
		offset += int(rr.RDLength) + len(rr.Name) + 10
	}

	return nil

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
