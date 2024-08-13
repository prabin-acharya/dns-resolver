// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1

package dns

import (
	"bytes"
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
