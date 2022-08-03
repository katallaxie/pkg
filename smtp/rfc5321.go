package smtp

// ReplyCode ...
type ReplyCode int

const (
	ReplyCodeSystemStatus                   = 211
	ReplyCodeHelpMessage                    = 214
	ReplyCodeServiceReady                   = 220
	ReplyCodeServiceClosing                 = 221
	ReplyCodeMailActionOkay                 = 250
	ReplyCodeMailActionCompleted            = 250
	ReplyCodeUserNotLocal                   = 251
	ReplyCodeCannotVerifyUser               = 252
	ReplyCodeStartMailInput                 = 354
	ReplyCodeServiceNotAvailable            = 421
	ReplyCodeMailboxUnavailable             = 450
	ReplyCodeLocalError                     = 451
	ReplyCodeInsufficientStorage            = 452
	ReplyCodeSyntaxError                    = 500
	ReplyCodeSyntaxErrorInParameters        = 501
	ReplyCodeCommandNotImplemented          = 502
	ReplyCodeCommandBadSequence             = 503
	ReplyCodeCommandParameterNotImplemented = 504
	ReplyCodeRequestActionNotTaken          = 550
	ReplyCodeUserNotLocalForThisHost        = 551
	ReplyCodeRequestedActionAborted         = 552
	ReplyCodeRequestedActionNotTaken        = 553
	ReplyCodeTransactionFailed              = 554
	ReplyCodeMailFromOrRcptToError          = 555
)
