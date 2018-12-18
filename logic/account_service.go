package logic

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/copernet/whccommon/log"
	common "github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/api/errs"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/model/view"
	"github.com/copernet/whcwallet/util"
	"github.com/gin-gonic/gin"
	"github.com/hgfischer/go-otp"
	"github.com/tidwall/gjson"
)

type AccountService struct{}

func (a *AccountService) ChallengeLogic(uuid string) (*view.ChallengeVO, error) {
	sessionId := util.CryptoSha256(config.GetConf().Private.SessionSecret + uuid)
	salt := util.CryptoSha256(config.GetConf().Private.ServerSecret + uuid)

	powChallenge := util.GenSalt(16)
	challenge := util.GenSalt(16)

	session, exist := model.GetSessionById(sessionId)
	if !exist {
		session = &common.Session{
			SessionID:  sessionId,
			Challenge:  &challenge,
			PChallenge: &powChallenge,
		}

		err := model.CreateSession(session)
		if err != nil {
			return nil, process(err)
		}
	} else {
		session.Challenge, session.PChallenge = &challenge, &powChallenge
		err := model.UpdateSessionById(session)
		if err != nil {
			return nil, process(err)
		}
	}

	return &view.ChallengeVO{
		Salt:         salt,
		PowChallenge: powChallenge,
		Challenge:    challenge,
	}, nil
}

/**
Create wallet && update Session PubKey
*/
func (a *AccountService) CreateWallet(c *gin.Context, param *view.WalletCreateParam) error {
	sessionId := util.CryptoSha256(config.GetConf().Private.SessionSecret + param.Uuid)
	session, exist := model.GetSessionById(sessionId)
	if !exist {
		return errors.New(errs.UuidNotExist)
	}

	if failedChallenge(*session.PChallenge, param.Nonce) || failedRecaptcha(c, param.RecaptchaResp) {
		return errors.New(errs.PermissionRefused)
	}

	wallet, exist := model.GetWalletById(param.Uuid)
	if exist {
		return errors.New(errs.UuidDuplicate)
	}

	now := time.Now()
	wallet = &common.Wallet{
		WalletID:   param.Uuid,
		LastLogin:  now,
		LastBackup: now,
		WalletBlob: param.Wallet,
		Email:      param.Email,
	}

	// CreateWallet
	err := model.CreateWallet(wallet)
	if err != nil {
		return process(err, param.Email)
	}
	//rsync send mail
	go util.Mail(param,c)

	// Reset session
	session.PubKey, session.PChallenge = param.PublicKey, nil
	//Update PublicKey
	return process(model.UpdateSessionById(session))
}

func failedRecaptcha(c *gin.Context, recaptchaResp string) bool {

	if recaptchaResp != "" {
		data := url.Values{
			"secret":   {config.GetConf().Private.RecaptchaPrivate},
			"response": {recaptchaResp},
		}

		uri := "https://www.google.com/recaptcha/api/siteverify"
		res, err := http.PostForm(uri, data)
		if err != nil {
			log.WithCtx(c).Info("failedPost:" + err.Error())
			return true
		}

		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.WithCtx(c).Error("request recapthca validation failed from google")
			return true
		}

		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.WithCtx(c).Info("failedRecaptcha:" + err.Error())
			return true
		}

		// here using gjson package for analyzing the request result
		if gjson.Get(string(content), "success").Bool() {
			// render a error page if failed
			return false
		}
	}

	return false
}

func failedChallenge(pcChallenge string, nonce string) bool {
	res := util.CryptoSha256(pcChallenge + nonce)
	ret := res[len(res)-len(config.GetConf().Private.LoginDifficulty):] !=
		config.GetConf().Private.LoginDifficulty

	return ret
}

/**
update Wallet info
*/
func (a *AccountService) UpdateWallet(c *gin.Context, param view.WalletUpdateParam) error {
	sessionId := util.CryptoSha256(config.GetConf().Private.SessionSecret + param.Uuid)
	session, exist := model.GetSessionById(sessionId)
	if !exist {
		return errors.New(errs.UuidNotExist)
	}

	log.WithCtx(c).Infof("pemStr%s,challenge:%s,signature:%s", session.PubKey, *session.Challenge, param.Signature)
	if !verifyRsaSignature(session.PubKey, *session.Challenge, param.Signature) {
		return errors.New(errs.PermissionRefused)
	}

	if param.MfaAction != "" && param.MfaAction != Add && param.MfaAction != Del {
		return errors.New(errs.MfaActionNil)
	}

	wallet, exist := model.GetWalletById(param.Uuid)
	if !exist {
		return errors.New(errs.UuidNotExist)
	}

	//verify mfa
	if param.MfaToken != "" && param.MfaAction != "" && !UpdateMfa(wallet, param) {
		log.WithCtx(c).Errorf("VerifyMfa Fail,token:%s,secret:%s", param.MfaToken, param.MfaSecret)
		return errors.New(errs.PermissionRefused)
	}

	//Reset session
	session.Challenge = nil
	model.UpdateSessionById(session)

	if param.Wallet != ""{
		wallet.WalletBlob = param.Wallet
	}

	if param.Email != ""{
		wallet.Email = param.Email
	}

	//update wallet
	return process(model.UpdateWallet(wallet),param.Email)
}

func verifyRsaSignature(pemStr string, challenge string, signature string) bool {
	block, _ := pem.Decode([]byte(pemStr))
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}

	hash := sha1.New()
	hash.Write([]byte(challenge))
	sign, _ := hex.DecodeString(signature)

	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
	if err != nil {
		return false
	}

	return true
}

func (a *AccountService) LoginWallet(c *gin.Context, param view.WalletLoginParam) (*view.WalletVO, error) {
	sessionId := util.CryptoSha256(config.GetConf().Private.SessionSecret + param.Uuid)
	session, exist := model.GetSessionById(sessionId)
	if !exist {
		return nil, errors.New(errs.UuidNotExist)
	}

	if failedChallenge(*session.PChallenge, param.Nonce) || failedRecaptcha(c, param.RecaptchaResp) {
		return nil, errors.New(errs.PermissionRefused)
	}

	wallet, exist := model.GetWalletById(param.Uuid)
	if !exist {
		return nil, errors.New(errs.UuidNotExist)
	}

	res := &view.WalletVO{Wallet: wallet.WalletBlob, Mfa: false}
	if wallet.Settings != "" {
		var setting model.WalletSetting
		err := json.Unmarshal([]byte(wallet.Settings), &setting)
		if err != nil {
			log.WithCtx(c).Errorf("Unmarshal Settings error:%s", err.Error())
			return nil, errors.New(errs.PermissionRefused)
		}

		if setting.MfaSecret.Value != "" {
			if param.Mfatoken == "" {
				return nil, errors.New(errs.MfaTokenNil)
			}

			secret, _ := util.AesDecrypt(setting.MfaSecret.Value, []byte(config.GetConf().Private.AesSecret))
			res.Mfa = VerifyMfa(param.Mfatoken, secret)
			if !res.Mfa {
				log.WithCtx(c).Info("VerifyMfa failed.")
				return nil, errors.New(errs.PermissionRefused)
			}

			//AesDecrypt asq
			asq := setting.Asq
			asq.Value, err = util.AesDecrypt(asq.Value, []byte(config.GetConf().Private.AesSecret))
			if err != nil {
				log.WithCtx(c).Errorf("AesDecrypt error:%s", err.Error())
				return nil, errors.New(errs.PermissionRefused)
			}

			var qa model.QuestionAnswer
			err = json.Unmarshal([]byte(asq.Value), &qa)
			if err != nil {
				log.WithCtx(c).Errorf("Unmarshal asq error:%s", err.Error())
				return nil, errors.New(errs.PermissionRefused)
			}
			res.Asq = qa.Question
		}
	} else if param.Mfatoken != "" {
		return nil, errors.New(errs.MfaTokenNotNil)
	}

	//Reset session
	session.PubKey, session.PChallenge = param.PublicKey, nil
	err := model.UpdateSessionById(session)
	if err != nil {
		return nil, err
	}
	//Update Last Login
	wallet.LastLogin = time.Now()
	err = model.UpdateWallet(wallet)
	if err != nil {
		return nil, process(err)
	}

	return res, nil
}

func (a *AccountService) VerifyWallet(param view.WalletVerifyParam) error {

	if param.Email != "" {
		noExist := model.VerifyWalletByEmail(param.Email)
		if !noExist {
			return errors.New(fmt.Sprintf(errs.DuplicateEmail, param.Email))
		}
	}

	return nil
}

func (a *AccountService) NewMfa(uuid string) *view.MfaVO {
	return NewMfa(uuid)
}

const (
	Add string = "add"
	Del string = "del"
)

func ProvisionUri(secret string, user string, issuer string) string {
	auth := "totp/"
	q := make(url.Values)
	q.Add("secret", secret)
	if issuer != "" {
		//q.Add("issuer", issuer)
		auth += issuer + ":"
	}

	return "otpauth://" + auth + user + "?" + q.Encode()
}

func VerifyMfa(token string, secret string) bool {
	totp := &otp.TOTP{
		Secret:         secret,
		Length:         uint8(otp.DefaultLength),
		Period:         uint8(otp.DefaultPeriod),
		IsBase32Secret: true,
	}
	return totp.Verify(token)
}

func NewMfa(uuid string) *view.MfaVO {
	secret := util.RandomBase32(16)
	prov := ProvisionUri(secret, uuid, "whcwallet")
	return &view.MfaVO{Secret: secret, Prov: prov}
}

func UpdateMfa(wallet *common.Wallet, param view.WalletUpdateParam) bool {
	var setting model.WalletSetting
	if param.MfaAction == Add {
		if !VerifyMfa(param.MfaToken, param.MfaSecret) {
			return false
		}

		sec, _ := util.AesEncrypt([]byte(param.MfaSecret), []byte(config.GetConf().Private.AesSecret))
		setting.MfaSecret = model.Secret{Value: sec, CreateAt: time.Now(), UpdateAt: time.Now()}

		asq := &model.QuestionAnswer{Question: param.Question, Answer: param.Answer}
		data, _ := json.Marshal(asq)
		asqEnc, _ := util.AesEncrypt(data, []byte(config.GetConf().Private.AesSecret))
		setting.Asq = model.Secret{Value: asqEnc, CreateAt: time.Now(), UpdateAt: time.Now()}
	} else {
		err := json.Unmarshal([]byte(wallet.Settings), &setting)
		if err != nil {
			return false
		}

		secret, _ := util.AesDecrypt(setting.MfaSecret.Value, []byte(config.GetConf().Private.AesSecret))
		if !VerifyMfa(param.MfaToken, secret) {
			return false
		}

		setting.MfaSecret.Value = ""
		setting.MfaSecret.UpdateAt = time.Now()
	}

	data, _ := json.Marshal(setting)
	wallet.Settings = string(data)

	return true
}
