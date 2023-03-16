package main

import (
	"bytes"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

var (
	templateHTMLFile      = "./templates/mail.html.gohtml"
	templatePlainTextFile = "./templates/mail.plain.gohtml"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func NewMail() *Mail {
	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	return &Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        mailPort,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		FromName:    os.Getenv("MAIL_FROM"),
	}
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	htmlMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainTextMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	smtpServer := &mail.SMTPServer{
		Host:           m.Host,
		Port:           m.Port,
		Username:       m.Username,
		Password:       m.Password,
		Encryption:     m.getEncryption(m.Encryption),
		KeepAlive:      false,
		ConnectTimeout: 10 * time.Second,
		SendTimeout:    10 * time.Second,
	}

	smtpClient, err := smtpServer.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject).
		SetBody(mail.TextPlain, plainTextMessage).
		AddAlternative(mail.TextHTML, htmlMessage)

	if len(msg.Attachments) > 0 {
		for _, a := range msg.Attachments {
			email.AddAttachment(a)
		}
	}

	if err = mail.SendMessage(
		email.GetFrom(),
		email.GetRecipients(),
		email.GetMessage(),
		smtpClient); err != nil {
		return err
	}

	return nil
}

// buildHTMLMessage builds HTML-formatted message and get it.
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	t, err := template.New("email-html").ParseFiles(templateHTMLFile)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

// buildPlainTextMessage builds plain-text based message and get it.
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	t, err := template.New("email-plain").ParseFiles(templatePlainTextFile)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

// inlineCSS converts css styles from HTML <style> tag
// into inline CSS inside html tags.
func (m *Mail) inlineCSS(fmtMsg string) (string, error) {
	options := premailer.Options{
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(fmtMsg, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

// getEncryption gets mail encryption type from specifying string.
func (m *Mail) getEncryption(encryptionStr string) mail.Encryption {
	switch strings.ToLower(encryptionStr) {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	default:
		return mail.EncryptionNone
	}
}
