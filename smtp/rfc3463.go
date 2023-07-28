package smtp

import "fmt"

var (
	// EnhancedStatusCodeUnknown is the default value for the enhanced status code.
	EnhancedStatusCodeUnknown EnhancedMailSystemStatusCode = EnhancedMailSystemStatusCode{-1, -1, -1}
)

// EnhancedStatusCodeClass represents a class of enhanced status codes.
type EnhancedStatusCodeClass int

// EnhancedStatusCodeSubject represents a classification of the status codes.
type EnhancedStatusCodeSubject int

// EnhancedStatusCodeDetail represents the detailed status.
type EnhancedStatusCodeDetail int

// EnhancedStatusCode is a data structure to contain enhanced
// mail system status codes from RFC 3463 (https://datatracker.ietf.org/doc/html/rfc3463).
type EnhancedStatusCode [3]int

// Class returns the class of the enhanced status code.
func (e EnhancedStatusCode) Class() EnhancedStatusCodeClass {
	return EnhancedStatusCodeClass(e[0])
}

// Subject returns the subject of the enhanced status code.
func (e EnhancedStatusCode) Subject() EnhancedStatusCodeSubject {
	return EnhancedStatusCodeSubject(e[1])
}

// Detail returns the detail of the enhanced status code.
func (e EnhancedStatusCode) Detail() EnhancedStatusCodeDetail {
	return EnhancedStatusCodeDetail(e[2])
}

// String returns the string representation of the enhanced status code.
func (e EnhancedMailSystemStatusCode) String() string {
	return fmt.Sprintf("%v.%v.%v", e[0], e[1], e[2])
}

const (
	// Signals a positive delivery action.
	EnhancedStatusCodeSuccess EnhancedStatusCodeClass = 2
	// Signals that there is a temporary failure in positively delivery action.
	EnhancedStatusCodePersistentTransientFailure EnhancedStatusCodeClass = 4
	// Signals that there is a permanent failure in the delivery action.
	EnhancedStatusCodePermanentFailure EnhancedStatusCodeClass = 5
)
