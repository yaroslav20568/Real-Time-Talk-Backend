package email

import (
	"fmt"
	"net/smtp"

	"gin-real-time-talk/config"
)

type EmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		host:     config.Env.SMTP.Host,
		port:     config.Env.SMTP.Port,
		username: config.Env.SMTP.User,
		password: config.Env.SMTP.Password,
		from:     config.Env.SMTP.From,
	}
}

func (e *EmailService) SendVerificationCode(to, code string) error {
	subject := "Код подтверждения"
	body := fmt.Sprintf(`
Здравствуйте!

Ваш код подтверждения: %s

Этот код действителен в течение 10 минут.

Если вы не запрашивали этот код, проигнорируйте это письмо.
`, code)

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	if e.username == "" || e.password == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}

	auth := smtp.PlainAuth("", e.username, e.password, e.host)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", to, subject, body))

	addr := fmt.Sprintf("%s:%s", e.host, e.port)
	err := smtp.SendMail(addr, auth, e.from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (e *EmailService) IsConfigured() bool {
	return e.username != "" && e.password != "" && e.from != ""
}
