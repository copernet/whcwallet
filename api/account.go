package api

import (
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/api/errs"
	"github.com/copernet/whcwallet/logic"
	"github.com/copernet/whcwallet/model/view"
	"github.com/gin-gonic/gin"
)

var accountService logic.AccountService

func Challenge(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		c.JSON(200, apiError(ErrEmptyQueryParam))
		return
	}

	challenge, err := accountService.ChallengeLogic(uuid)
	if err != nil {
		log.WithCtx(c).Warnf("challenge error:%s", err.Error())

		c.JSON(200, apiError(ErrChallenge))
		return
	}

	c.JSON(200, apiSuccess(challenge))
}

func Create(c *gin.Context) {
	var param view.WalletCreateParam
	if err := c.ShouldBind(&param); err != nil {
		log.WithCtx(c).Errorf("Form valid error:%s", err.Error())
		c.JSON(200, apiError(ErrFormItems))
		return
	}

	err := accountService.CreateWallet(c, &param)
	if err != nil {
		c.JSON(200, apiErrorWithMsg(ErrCreateWallet, err.Error()))
		return
	}

	c.JSON(200, apiSuccess(nil))
}

func Update(c *gin.Context) {
	var param view.WalletUpdateParam
	if err := c.ShouldBind(&param); err != nil {
		log.WithCtx(c).Errorf("Form valid error:%s", err.Error())
		c.JSON(200, apiError(ErrFormItems))
		return
	}

	err := accountService.UpdateWallet(c, param)
	if err != nil {
		c.JSON(200, apiErrorWithMsg(ErrUpdateWallet, err.Error()))
		return
	}

	c.JSON(200, apiSuccess(nil))
}

func Login(c *gin.Context) {
	var param view.WalletLoginParam
	if err := c.ShouldBind(&param); err != nil {
		log.WithCtx(c).Errorf("Form valid error:%s", err.Error())
		c.JSON(200, apiError(ErrFormItems))
		return
	}

	wallet, err := accountService.LoginWallet(c, param)
	if err != nil {
		code := ErrLogin
		if err.Error() == errs.MfaTokenNil {
			code = ErrMfaToken
		}

		c.JSON(200, apiErrorWithMsg(code, err.Error()))
		return
	}

	c.JSON(200, apiSuccess(wallet))
}

/**
The api is used for web ajax verify email unique or phone unique
*/
func Verify(c *gin.Context) {
	var param view.WalletVerifyParam
	if err := c.ShouldBind(&param); err != nil {
		log.WithCtx(c).Errorf("Form valid error:%s", err.Error())
		c.JSON(200, apiError(ErrFormItems))
		return
	}

	err := accountService.VerifyWallet(param)
	if err != nil {
		c.JSON(200, apiErrorWithMsg(ErrVerify, err.Error()))
		return
	}

	c.JSON(200, apiSuccess(nil))
}

func NewMfa(c *gin.Context) {
	uuid, _ := c.GetPostForm("uuid")
	if uuid == "" {
		c.JSON(200, apiError(ErrEmptyQueryParam))
		return
	}

	vo := accountService.NewMfa(uuid)
	c.JSON(200, apiSuccess(vo))

}
