package errs

const (
	DefaultMessage    = "there is a problem on the server, please wait a moment"
	DuplicateEmail    = "your email:%s is registered. Please Change it"
	UuidDuplicate     = "your uuid:%s is registered. Please Change it"
	UuidNotExist      = "uuid not exist in session"
	PermissionRefused = "sorry, you don't have permission to the system"
	MfaActionNil      = "action only support 'add' or del"
	MfaTokenNil      = "you account has set login with mfa_token"
	MfaTokenNotNil      = "you account hasn't set mfa_token"

)
