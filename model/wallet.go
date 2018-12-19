package model

import (
	"time"

	common "github.com/copernet/whccommon/model"
)

type WalletSetting struct {
	MfaSecret Secret `json:"mfa_secret"`
	Asq       Secret `json:"asq"`
}

type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Secret struct {
	Value    string    `json:"Value"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

func GetWalletById(id string) (*common.Wallet, bool) {
	var wallet common.Wallet
	exist := walletdb.Where("wallet_id=?", id).First(&wallet).RecordNotFound()
	if exist {
		return nil, false
	}

	return &wallet, true
}

func CreateWallet(model *common.Wallet) error {
	return walletdb.Create(model).Error
}

func UpdateWallet(wallet *common.Wallet) error {
	return walletdb.Save(wallet).Error
}

func VerifyWalletByEmail(email string) bool {
	var wallet common.Wallet
	return walletdb.Where("email=?", email).First(&wallet).RecordNotFound()
}
