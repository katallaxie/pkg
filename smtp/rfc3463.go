package smtp

import "fmt"

// EnhancedStatusCodeClass represents a class of enhanced status codes.
type EnhancedStatusCodeClass int

// EnhancedStatusCodeSubject represents a classification of the status codes.
type EnhancedStatusCodeSubject int

// EnhancedStatusCodeDetail represents the detailed status.
type EnhancedStatusCodeDetail int

// EnhancedStatusCode is a data structure to contain enhanced
// mail system status codes from RFC 3463 (https://datatracker.ietf.org/doc/html/rfc3463).
type EnhancedStatusCode [3]int

// String returns the string representation of the enhanced status code.
func (e EnhancedStatusCode) String() string {
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
