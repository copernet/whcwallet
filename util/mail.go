package util

import (
	"bytes"
	"context"
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model/view"
	"github.com/domodwyer/mailyak"
	"html/template"
	"net/smtp"
	"path/filepath"
	"strconv"
	"strings"
)

var mailConf = config.GetConf().Mail

func Mail(param *view.WalletCreateParam, c context.Context) {
	err := MailTo(param)
	if err != nil {
		log.WithCtx(c).Errorf("Send Email Error:%s", err.Error())
	}
}

func MailTo(param *view.WalletCreateParam) error {
	mail := mailyak.New(mailConf.Host+":"+strconv.Itoa(mailConf.Port), smtp.PlainAuth("", mailConf.SmtpUser, mailConf.SmtpPass, mailConf.Host))

	mail.To(param.Email)
	mail.From(mailConf.UserName)
	mail.FromName(mailConf.FromName)

	mail.Subject("Your Wallet ID")

	tpl, err := loadTemplate()
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)
	data := MailContent{param.Uuid}
	tpl.Execute(b, data)
	mail.HTML().Set(b.String())
	// And you're done!
	return mail.Send()
}

type MailContent struct {
	WalletID string
}

func loadTemplate() (*template.Template, error) {
	path, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	lastIndex := strings.Index(path, config.ProjectLastDir) + len(config.ProjectLastDir)
	correctPath := path[:lastIndex] + "/static/mail_content.html"

	return template.ParseFiles(correctPath)
}
