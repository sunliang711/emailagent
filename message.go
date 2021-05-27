// Package emailagent
// @Project:       emailagent
// @File:          message.go
// @Author:        eagle
// @Create:        2021/05/27 15:49:27
// @Description:
package emailagent

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
)

type Message struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	attachments map[string][]byte
}

func NewMessage(subject string, body string) *Message {
	return &Message{
		Subject:     subject,
		Body:        body,
		attachments: make(map[string][]byte),
	}
}

func (m *Message) Attach(filename string, contents []byte) {
	m.attachments[filename] = contents
}

func (m *Message) Build(isHTML bool) (data []byte, err error) {
	if len(m.To) == 0 {
		err = errors.New("no recipient")
		return
	}

	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.To, ",")))

	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\r\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	}
	if isHTML {
		buf.WriteString("Content-Type: text/html; charset=utf-8\r\n")
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	}

	buf.WriteString("\r\n")
	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.attachments {
			buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.WriteString("\r\n")
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\r\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	data = buf.Bytes()

	return
}
