package util

import (
	"bytes"
	"fmt"
	logCtx "github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model/view"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSync(t *testing.T) {
	go testGo()
	fmt.Println("hello")
}
func testGo() {
	fmt.Println("hi")
}

func TestMailTo(t *testing.T) {
	to := "hongjia.hu@bitmain.com"
	mailConf := config.GetConf().Mail

	// Create a new message.
	m := gomail.NewMessage()

	tpl, _ := loadTemplate()
	b := bytes.NewBuffer(nil)
	data := MailContent{"Test UUID"}
	tpl.Execute(b, data)
	// Set the main email part to use HTML.
	m.SetBody("text/html", b.String())

	// Construct the message headers, including a Configuration Set and a Tag.
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(mailConf.UserName, mailConf.FromName)},
		"To":      {to},
		"Subject": {"Your Wallet ID"},
	})

	d := gomail.NewDialer(mailConf.Host, mailConf.Port, mailConf.SmtpUser, mailConf.SmtpPass)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent!")
	}
}

func TestMailTpl(t *testing.T) {
	path, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	lastIndex := strings.Index(path, config.ProjectLastDir) + len(config.ProjectLastDir)
	correctPath := path[:lastIndex] + "/tpl/mail_content.html"

	tpl, err := template.ParseFiles(correctPath)
	if err != nil {
		log.Fatal(err)
	}

	data := MailContent{"2abb1c16-083a-4659-9f5b-3e83419988e0"}
	err = tpl.Execute(os.Stdout, data)
	if err != nil {
		log.Fatal(err)
	}
}

func TestMail(t *testing.T) {
	walletId := "64db5a00-4470-42ee-a674-33284073f103"
	param := view.WalletCreateParam{"hongjia.hu@bitmain.com", "10837", "test", walletId, "d09a641fdf12d09b1a8b9655fea26440dd3f1cf34bc1a643c0f62e0743da09e3eeff2783126cf24d96fb4aebac75a19d94df33bc9f7f46e407d35d1435272418d61e6cb40b1abe2c93c04fd6a19f8c0e88f8c6ae244749c3f18192f2bbc2f493", ""}
	Mail(&param, logCtx.NewContext())
}
