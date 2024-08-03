// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1

package dns

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
type Question struct {
	QName  string
	QType  uint16
	QClass uint16
}
