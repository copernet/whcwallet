package logic

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/util"
	"github.com/gin-gonic/gin"
)

func TestGetSession(t *testing.T) {
	session := util.CryptoSha256(config.GetConf().Private.SessionSecret + "f86a48e2-89a4-4327-9b1c-a22f193d042e")
	sessions, exist := model.GetSessionById(session)
	if exist {
		t.Error("uuid not exist")
	}

	fmt.Println(sessions)

}

func TestRsa(t *testing.T) {
	pemstr := "-----BEGIN PUBLIC KEY-----\nMIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgHSEy24daT29CGLEPuKT1adgji1E\nOuMFoOCz5hv9yx/MJx+Npn0Bz1cnALHPAbnaKrhQvuqtKmkxnQq93cw9tl5yo+Im\n+3Srrm/6hQTsYgXhovCbgo729jNTGlSmVkkEXabKso7VvNqgDIckBpaVkmvKP5Yu\ncPWhrMl/ujjtWnBTAgMBAAE=\n-----END PUBLIC KEY-----"
	block, _ := pem.Decode([]byte(pemstr))
	pubkey, _ := x509.ParsePKIXPublicKey(block.Bytes)

	hash := sha1.New()
	hash.Write([]byte("9i9B80bn8Rkhns91UTm77TZP8kdJItFu"))
	sign, _ := hex.DecodeString("00fea99ad43a1d21338486d35f9e02ef128bca40353ad847a3834d58604a3ca5c844ca94e5c8b084b051197195a5667f70e94a822a28e7c9af37198d9056af00c17ad4d912408d0abf0831c865f2594f19cfa420062f12f5a067d6015b60016a4a9e20465882dba6ce96fdc178e8046697b00fbfccc587d8546bdaa4af59464e")

	err := rsa.VerifyPKCS1v15(pubkey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ok")
}

func TestFailedRecaptcha(t *testing.T) {
	content := "03AL4dnxoP1oiSHG-enRJiK5FqMdGvXcuOanLIZUnj9SnIgBPlLCAQTpjZrZXnMNputYUETwi8-cHGKXuldE9Mioj8m0vhxkHmE85H7T3JwBfRCSPCVmsQF_oYMW3sjSO6GIYRV0ywRArUZ9j2kZ6J_31N_v8BANoQB2o0nlPz5_a-hYpZlMgNWo09lwM7kgRUDVE1PIHhRgTb3vPEPNmzZUPcLFHzxVoAkAR916wDMb8bIXNAy60axUkbifnxVd8umyjN5Jj1QWPlhG-hsKWv0Hm-_6DklS4KVQ"
	res := failedRecaptcha(&gin.Context{}, content)
	if !res {
		t.FailNow()
	}
}

func TestMfa(t *testing.T) {
	sec, _ := util.AesEncrypt([]byte("test"), []byte(config.GetConf().Private.AesSecret))
	fmt.Println(sec)

}

func TestVerifyChallenge(t *testing.T) {
	bool := failedChallenge("430e58b86534cbe9fb1dafb317340859", "2678")
	fmt.Println(bool)
}

func TestBase64Decode(t *testing.T) {
	byts, _ := base64.StdEncoding.DecodeString("LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0NCk1JR2ZNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0R05BRENCaVFLQmdRQ2Y1ZFpvNzRtZ0ZDbTFxNXhiSEg1YkRHdVQNCkpnQUVnTU5XQ2pOaTdLMTdYZXZGN1VMaDhOVHJ0YlRDRW9MNFJudjhEZ1g3dHVvc1NHY2tXeXJpMzZ6MTZ0c0MNCmRtZDhmQUtMRFBHYUtsOFlxWVo2ZVZBN3VVVjRFajIrTmZrQ0FMN2FqSjRVcEd0Umw4MUdqSHY0aFViL0I5YisNCnpMcEp0KzgwNFU1UkY4RmVqd0lEQVFBQg0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tDQo=")
	fmt.Println(string(byts))
}

func TestVerifyMfa(t *testing.T) {

	//secret, _ := util.AesDecrypt(setting.MfaSecret.Value, []byte(config.GetConf().Private.AesSecret))
	//res.Mfa = util.VerifyMfa(param.Mfatoken, secret)

	val, _ := util.AesDecrypt("9p163XG46XT93Uyivq7WmOeaGf3v7EsEvHbAGSWraFs=", []byte(config.GetConf().Private.AesSecret))
	var qa model.QuestionAnswer
	err := json.Unmarshal([]byte(val), &qa)
	if err != nil {
		t.Errorf("Unmarshal asq error:%s", err.Error())
	}
	fmt.Println(qa)

	mfaSecret := "PZ2Z4KZPDZP5V25E"
	sec, err := util.AesEncrypt([]byte(mfaSecret), []byte(config.GetConf().Private.AesSecret))
	fmt.Println(sec)
	fmt.Println(err)

	b := VerifyMfa("582121", "GZV6UBTHWXIYXHBR")
	fmt.Println(b)

}
