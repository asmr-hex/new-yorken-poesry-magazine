package types

import (
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
	Sender     string
	SentEmails []string
}

func NewDevEmailer(sender string) *DevEmailer {
	emailer := &DevEmailer{
		Sender:     sender,
		SentEmails: []string{},
	}

	return emailer
}

func (d *DevEmailer) SendEmail(subject, body string, recipients ...string) error {
	// TODO (cw|10.19.2018) configure this to send emails from
	// the localhost smtp server using sendmail or mail commandline tool

	return nil
}

// TODO (cw|10.19.2018) eventually create a TestEmailer which will take a channel
// so the tests can immediately recieve the verification email over it so they can
// automatically verify without having to actually check email...
