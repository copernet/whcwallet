package api

import (
	"encoding/json"

	"github.com/copernet/whcwallet/api/errs"
	"github.com/gin-gonic/gin"

	"log"
	"testing"

	"fmt"

	"github.com/copernet/whcwallet/model/view"
	"github.com/satori/go.uuid"
)

func TestChallenge(t *testing.T) {

	u1, err := uuid.NewV4()
	log.Println(u1.String())
	vo, err := accountService.ChallengeLogic("570ef0d2-8e35-46b3-e382-6cd809a4a986")

	data, err := json.Marshal(vo)
	if err != nil {
		t.Error("error:" + err.Error())
		return
	}
	log.Println(string(data))
}

func TestCreate(t *testing.T) {
	walletId := "c2e69794-857b-4934-c09d-95574f75839a"
	//accountService.ChallengeLogic(walletId)
	param := view.WalletCreateParam{"38366472456@qq.com", "1801", "-----BEGIN PUBLIC KEY-----\r\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDianj3GLLiRnaIRRQB13hdE4mk\r\nULjs8eG3vGF983WkW5ZmLVbPYryK/OUawjWTV7N7AoqoOjgtlq19wW6du+wz77cV\r\nixQy7WnW2MFs2hgoKWlWsheK88gSJbgYDjBTKWQfk9PcGckZp/jywNP1ii3CbMc8\r\nRIrO9lHrQql9gM6NSwIDAQAB\r\n-----END PUBLIC KEY-----\r\n", walletId, "d09a641fdf12d09b1a8b9655fea26440dd3f1cf34bc1a643c0f62e0743da09e3eeff2783126cf24d96fb4aebac75a19d94df33bc9f7f46e407d35d1435272418d61e6cb40b1abe2c93c04fd6a19f8c0e88f8c6ae244749c3f18192f2bbc2f493", ""}
	err := accountService.CreateWallet(&gin.Context{}, &param)
	fmt.Println(err.Error())
}

func TestUpdate(t *testing.T) {
	//accountService.ChallengeLogic(walletId)
	var param view.WalletUpdateParam
	param.Uuid = "c2e69794-857b-4934-c09d-95574f75839a"
	param.Wallet = ""
	param.Signature = "9da0d999011d35fad2f68116e2639fdb5124b88683e907c443213da61b579d3a3eb63497b2cafb206390e35dbe05990822923689f623731031a88b2f7eddef7348a64bd719db56d89987811986d8cd063ab46f8cd63f1c136cb2194368420e9a8c2f1ff4da991e14e346521891e4d58aa113a038c6aef291c40fd740927dd5bb"
	param.Email = "hd012@163.com"
	param.Answer = "Hu"
	param.Question = "Hu"
	param.MfaAction = "add"
	param.MfaSecret = "2YRJIU34DV2N4BCY"
	param.MfaToken = ""

	accountService.UpdateWallet(nil, param)
}

func TestLogin(t *testing.T) {
	walletId := "c2e69794-857b-4934-c09d-95574f75839a"
	//accountService.ChallengeLogic(walletId)

	param := view.WalletLoginParam{"087718", "xs", walletId, "", ""}
	res, err := accountService.LoginWallet(&gin.Context{}, param)
	if err != nil {
		code := ErrLogin
		if err.Error() == errs.MfaTokenNil {
			code = ErrMfaToken
		}

		fmt.Println(code)
		return
	}
	fmt.Println(res)
}

func TestVerify(t *testing.T) {
	param := view.WalletVerifyParam{"xx@qq.com", ""}
	accountService.VerifyWallet(param)
}

func TestNewMfa(t *testing.T) {
	vo := accountService.NewMfa("f9008f63-345e-46fe-e999-cd14bff0d519")
	fmt.Println(vo)
}
