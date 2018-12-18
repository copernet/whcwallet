package view

/**
Define Api Parameters,use by *Param
*/

type WalletCreateParam struct {
	Email         string `form:"email" binding:"required"`
	Nonce         string `form:"nonce" binding:"required"`
	PublicKey     string `form:"public_key" binding:"required"`
	Uuid          string `form:"uuid" binding:"required"`
	Wallet        string `form:"wallet" binding:"required"`
	RecaptchaResp string `form:"recaptcha_response_field"`
}

type WalletUpdateParam struct {
	Uuid      string `form:"uuid" binding:"required"`
	Wallet    string `form:"wallet"`
	Signature string `form:"signature" binding:"required"`
	Email     string `form:"email"`

	//about MFA
	MfaSecret string `form:"mfa_secret"`
	MfaToken  string `form:"mfa_token"`
	MfaAction string `form:"mfa_action"`
	Question  string `form:"question"`
	Answer    string `form:"answer"`
}

type WalletLoginParam struct {
	Nonce         string `form:"nonce" binding:"required"`
	PublicKey     string `form:"public_key" binding:"required"`
	Uuid          string `form:"uuid" binding:"required"`
	Mfatoken      string `form:"mfatoken"`
	RecaptchaResp string `form:"recaptcha_response_field"`
}

type WalletVerifyParam struct {
	Email string `form:"email"`
	Phone string `form:"phone"`
}

/**
Define Value Object,use by *VO
*/
type ChallengeVO struct {
	Salt         string `json:"salt"`
	PowChallenge string `json:"pow_challenge"`
	Challenge    string `json:"challenge"`
}

type WalletVO struct {
	Wallet string `json:"wallet"`
	Mfa    bool   `json:"mfa"`
	Asq    string `json:"asq"`
}

type MfaVO struct {
	Secret string `json:"secret"`
	Prov   string `json:"prov"`
}

type Version struct {
	Version     string `json:"version"`
	VersionCode int    `json:"versionCode"`
	Download    string `json:"download"`
	Description string `json:"description"`
}
