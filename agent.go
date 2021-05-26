package emailagent

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"

	log "github.com/sirupsen/logrus"
)

const (
	HTML_MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

// buildMessage
// 2019/09/23 21:53:09
func buildMessage(recipients []string, from, subject, body string, isHtml bool) string {
	f := mail.Address{Address: from}

	recipientsMsg := ""
	for _, recipient := range recipients {
		t := mail.Address{Address: recipient}
		recipientsMsg += fmt.Sprintf("To: %s\r\n", t.String())
	}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = f.String()
	// headers["To"] = t.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += recipientsMsg
	if isHtml {
		message += HTML_MIME
	}
	message += "\r\n" + body
	return message
}

// EmailAgent
// 2019/09/23 22:00:36
type EmailAgent struct {
	Host     string
	Port     int
	User     string
	Password string
	Client   *smtp.Client
}

// NewEmailAgent
func NewEmailAgent(host string, port int, user, password string) (*EmailAgent, error) {
	agent := &EmailAgent{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Client:   nil,
	}
	err := agent.init()
	if err != nil {
		return nil, err
	}
	return agent, nil
}

// init
func (a *EmailAgent) init() error {
	auth := smtp.PlainAuth("", a.User, a.Password, a.Host)
	tlsConfig := &tls.Config{
		ServerName: a.Host,
		// InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", a.Host, a.Port), tlsConfig)
	if err != nil {
		log.Errorf("Dial error: %v", err)
		return err
	}

	c, err := smtp.NewClient(conn, a.Host)
	if err != nil {
		log.Errorf("NewClient error: %v", err)
		return err
	}

	if err = c.Auth(auth); err != nil {
		log.Errorf("Auth error: %v", err)
		return err
	}

	a.Client = c

	return nil
}


// SendEmail
func (a *EmailAgent) SendEmail(recipients []string, subject, body string, isHtml bool) error {
	if a.Client == nil {
		msg := fmt.Sprintf("EmailAgent must init!")
		log.Errorf(msg)
		err := fmt.Errorf(msg)
		return err
	}
	if err := a.Client.Mail(a.User); err != nil {
		log.Errorf("Mail error: %v", err)
		return err
	}

	for _, recipient := range recipients {
		if err := a.Client.Rcpt(recipient); err != nil {
			log.Errorf("Rcpt error: %v", err)
			return err
		}
	}

	w, err := a.Client.Data()
	if err != nil {
		log.Errorf("Data error: %v", err)
		return err
	}
	defer w.Close()

	message := buildMessage(recipients, a.User, subject, body, isHtml)
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Errorf("Write error: %v", err)
		return err
	}
	return nil
}

// Close
// 2019/09/23 22:17:53
func (a *EmailAgent) Close() error {
	if a.Client == nil {
		msg := fmt.Sprintf("EmailAgent need init!")
		err := fmt.Errorf("%s", msg)
		log.Errorf(msg)
		return err
	}
	return a.Client.Close()
}
