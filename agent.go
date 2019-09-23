package emailagent

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"

	log "github.com/sirupsen/logrus"
)

// buildMessage TODO
// 2019/09/23 21:53:09
func buildMessage(from, to, subject, body string) string {
	f := mail.Address{"", from}
	t := mail.Address{"", to}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = f.String()
	headers["To"] = t.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	return message
}

// EmailAgent TODO
// 2019/09/23 22:00:36
type EmailAgent struct {
	Host     string
	Port     int
	User     string
	Password string
	Client   *smtp.Client
}

//NewEmailAgent TODO
func NewEmailAgent(host string, port int, user, password string) *EmailAgent {
	return &EmailAgent{host, port, user, password, nil}
}

//Init TODO
func (a *EmailAgent) Init() error {
	auth := smtp.PlainAuth("", a.User, a.Password, a.Host)
	tlsConfig := &tls.Config{
		ServerName: a.Host,
		// InsecureSkipVerify: true,
	}

	log.Infof("Dial %s:%d...", a.Host, a.Port)
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", a.Host, a.Port), tlsConfig)
	if err != nil {
		log.Errorf("Dial error: %v", err)
		return err
	}
	log.Infof("Dial OK.")

	c, err := smtp.NewClient(conn, a.Host)
	if err != nil {
		log.Errorf("NewClient error: %v", err)
		return err
	}

	log.Infof("Auth ...")
	if err = c.Auth(auth); err != nil {
		log.Errorf("Auth error: %v", err)
		return err
	}

	log.Infof("Auth OK.")
	a.Client = c

	return nil
}

// SendEmail TODO
// 2019/09/23 22:09:58
func (a *EmailAgent) SendEmail(to, subject, body string) error {
	if a.Client == nil {
		msg := fmt.Sprintf("EmailAgent must init!")
		err := fmt.Errorf(msg)
		log.Errorf(msg)
		return err
	}
	log.Infof("Mail ...")
	if err := a.Client.Mail(a.User); err != nil {
		log.Errorf("Mail error: %v", err)
		return err
	}
	log.Infof("Mail OK.")

	log.Infof("Rctp ...")
	if err := a.Client.Rcpt(to); err != nil {
		log.Errorf("Rcpt error: %v", err)
		return err
	}
	log.Infof("Rctp OK.")

	log.Infof("Data ...")
	w, err := a.Client.Data()
	if err != nil {
		log.Errorf("Data error: %v", err)
		return err
	}
	defer w.Close()
	log.Infof("Data OK.")

	message := buildMessage(a.User, to, subject, body)
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Errorf("Write error: %v", err)
		return err
	}
	return nil
}

// Close TODO
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
