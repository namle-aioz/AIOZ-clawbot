package mail_actions

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"
)

type ResendMailSender struct {
	client *resend.Client
	from   string
}

func NewResendMailSenderAction(client *resend.Client, from string) MailSender {
	return &ResendMailSender{
		client: client,
		from:   from,
	}
}

func (s *ResendMailSender) SendOTPMail(ctx context.Context, input SendMailInput) error {
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      input.To,
		Subject: input.Subject,
		Html:    input.Html,
		Cc:      input.Cc,
		Bcc:     input.Bcc,
		ReplyTo: input.ReplyTo,
	}

	res, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("send mail failed: %w", err)
	}

	fmt.Println("email sent:", res.Id)
	return nil
}
