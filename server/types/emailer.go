package types

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	mailgun "github.com/mailgun/mailgun-go"
)

type Emailer interface {
	// sends an email with subject, body, & recipients
	SendEmail(string, string, ...string) error
}

type MailGunEmailer struct {
	*mailgun.MailgunImpl
	Domain string
	ApiKey string
	Sender string
}

func NewMailGunEmailer(domain, key, sender string) *MailGunEmailer {
	return &MailGunEmailer{
		MailgunImpl: mailgun.NewMailgun(domain, key),
		Domain:      domain,
		ApiKey:      key,
		Sender:      sender,
	}
}

func (e *MailGunEmailer) SendEmail(subject, body string, recipients ...string) error {
	message := e.NewMessage(e.Sender, subject, body, recipients...)
	_, _, err := e.Send(message)
	if err != nil {
		return err
	}

	return nil
}

type DevEmailer struct {
	Host       string
	Port       int
	Sender     string
	SentEmails []string
}

// for the dev environment, we are using maildev (https://danfarrelly.nyc/MailDev/)
// which is a handy smtp server + web gui for viewing and testing emails during
// development! It is setup as a linked docker container (hostname = maildev)
// all we need to do is simply send email messages to the maildev host at port 25
// and then we can open up localhost:1080 to view our emails!
func NewDevEmailer(sender string) *DevEmailer {
	emailer := &DevEmailer{
		Host:       "maildev",
		Port:       25,
		Sender:     sender,
		SentEmails: []string{},
	}

	return emailer
}

func (d *DevEmailer) SendEmail(subject, body string, recipients ...string) error {
	// convert recipient slice into comma seperated string
	recipientsString := ""
	for idx, recipient := range recipients {
		recipientsString += recipient
		if len(recipients) != idx+1 {
			recipientsString += ", "
		}
	}

	// construct email headers
	headers := make(map[string]string)
	headers["From"] = d.Sender
	headers["To"] = recipientsString
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// customize the tls config to skip bad/non-existing certs
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         d.Host,
	}

	// open connection
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", d.Host, d.Port))
	if err != nil {
		return err
	}

	err = c.StartTLS(tlsconfig)
	if err != nil {
		return err
	}

	// To && From
	if err = c.Mail(d.Sender); err != nil {
		return err
	}

	// repeat the RCPT command for each recipient
	for _, recipient := range recipients {
		if err = c.Rcpt(recipient); err != nil {
			return err
		}
	}

	// send the DATA command, which tells the smtp server that we
	// are ready to write our message!
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}

// TODO (cw|10.19.2018) eventually create a TestEmailer which will take a channel
// so the tests can immediately recieve the verification email over it so they can
// automatically verify without having to actually check email...
