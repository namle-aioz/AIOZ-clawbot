package mail_actions

import "context"

type SendMailInput struct {
	To      []string
	Subject string
	Html    string
	Cc      []string
	Bcc     []string
	ReplyTo string
}

type MailSender interface {
	SendOTPMail(ctx context.Context, input SendMailInput) error
}
