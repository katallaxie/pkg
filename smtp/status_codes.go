package smtp

// EnhancedMailSystemStatusCode ...
// https://datatracker.ietf.org/doc/html/rfc3463
type EnhancedMailSystemStatusCode [3]int

// StatusCode ...
type StatusCode struct {
	replyCode    int
	enhancedCode EnhancedStatusCode
	message      string
}

// ReplyCode ...
func (s *StatusCode) ReplyCode() int {
	return s.replyCode
}

// EnhancedStatusCode ...
func (s *StatusCode) EnhancedStatusCode() EnhancedStatusCode {
	return s.enhancedCode
}

// Message ...
func (s *StatusCode) Message() string {
	return s.message
}

// Error ...
func (s *StatusCode) Error() string {
	return s.message
}

// Error ...
type Error struct {
	statusCode StatusCode
}

// Error ...
func (e *Error) Error() string {
	return e.statusCode.message
}

func (e *Error) Temporary() bool {
	return e.statusCode.replyCode/100 == 4
}

// NewStatusCode ...
func (s *StatusCode) NewStatusCode(replyCode int, enhancedCode EnhancedStatusCode, message string) *StatusCode {
	return &StatusCode{
		replyCode:    replyCode,
		enhancedCode: enhancedCode,
		message:      message,
	}
}

// ErrorFromStatus ...
func ErrorFromStatus(replyCode int, enhancedCode EnhancedStatusCode, message string) error {
	return &Error{
		statusCode: StatusCode{
			replyCode:    replyCode,
			enhancedCode: enhancedCode,
			message:      message,
		},
	}
}
