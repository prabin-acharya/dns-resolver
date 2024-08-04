// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1

package dns

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
		Answers:       answers,
		AuthorityRRs:  authorityRRs,
		AdditionalRRs: additionalRRs,
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
	m.Header.FromBytes(data[:12])
	offset := 12

	m.Questions = make([]Question, m.Header.QDCount)
	for i := 0; i < int(m.Header.QDCount); i++ {
		q := &m.Questions[i]
		n, err := q.FromBytes(data[offset:])
		if err != nil {
			return err
		}
		offset += n
	}

	// Similar parsing would be needed for Answers, AuthorityRRs, and AdditionalRRs

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
