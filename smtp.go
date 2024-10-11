package fsq

import (
	"fmt"
	"net/smtp"
)

type SmtpConfig struct {
	Host        string
	Port        int
	DefaultFrom string
}

type Smtp struct {
	Cfg *SmtpConfig
}

func NewSmtp(Cfg *SmtpConfig) *Smtp {
	return &Smtp{
		Cfg: Cfg,
	}
}

func (s *Smtp) SendMail(to string, subject string, body string) error {
	headers := make(map[string]string)
	headers["From"] = s.Cfg.DefaultFrom
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%d", s.Cfg.Host, s.Cfg.Port)

	err := smtp.SendMail(
		addr,
		nil,
		s.Cfg.DefaultFrom,
		[]string{to},
		[]byte(message),
	)
	if err != nil {
		return err
	}

	return nil
}
