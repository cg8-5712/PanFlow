package service

import (
	"fmt"

	"panflow/internal/config"

	"gopkg.in/gomail.v2"
)

type MailService struct {
	cfg *config.PanflowConfig
}

func NewMailService(cfg *config.PanflowConfig) *MailService {
	return &MailService{cfg: cfg}
}

// Send sends an email with the given subject and body
func (s *MailService) Send(subject, body string) error {
	if !s.cfg.MailSwitch {
		return nil
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", s.cfg.MailFromAddress, s.cfg.MailFromName)
	m.SetAddressHeader("To", s.cfg.MailToAddress, s.cfg.MailToName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.cfg.MailHost, s.cfg.MailPort, s.cfg.MailUsername, s.cfg.MailPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send mail: %w", err)
	}
	return nil
}

// SendAccountAlert sends an alert when an account encounters an error
func (s *MailService) SendAccountAlert(accountName, reason string) error {
	subject := fmt.Sprintf("[PanFlow] 账号异常通知：%s", accountName)
	body := fmt.Sprintf(`
<h3>账号异常通知</h3>
<p>账号 <strong>%s</strong> 发生异常：</p>
<p>%s</p>
`, accountName, reason)
	return s.Send(subject, body)
}

// SendParseError sends an alert when a parse error occurs
func (s *MailService) SendParseError(surl, errMsg string) error {
	subject := "[PanFlow] 解析失败通知"
	body := fmt.Sprintf(`
<h3>解析失败通知</h3>
<p>分享链接 <strong>%s</strong> 解析失败：</p>
<p>%s</p>
`, surl, errMsg)
	return s.Send(subject, body)
}
