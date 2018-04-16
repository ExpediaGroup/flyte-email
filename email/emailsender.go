/*
Copyright (C) 2018 Expedia Group.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package email

import (
	"bytes"
	"net/smtp"
	"strings"
)

const crlf = "\r\n"

type EmailSender interface {
	Send(from string, to []string, subject, body string, isHtmlEmail bool) error
}

type emailSender struct {
	send func(from string, to []string, msg []byte) error
}

func NewEmailSender(smtpServerAddr string) EmailSender {
	return &emailSender{send: sendFunc(smtpServerAddr)}
}

func (e *emailSender) Send(from string, to []string, subject, body string, isHtmlEmail bool) error {
	return e.send(from, to, createEmail(from, to, subject, body, isHtmlEmail))
}

func createEmail(from string, to []string, subject, body string, isHtmlEmail bool) []byte {
	msg := new(bytes.Buffer)
	addHeaders(from, to, subject, isHtmlEmail, msg)
	msg.WriteString(body)
	return msg.Bytes()
}

func addHeaders(from string, to []string, subject string, isHtmlEmail bool, msg *bytes.Buffer) {
	addHeader("MIME-version", "1.0;", msg)

	if isHtmlEmail {
		addHeader("Content-Type", "text/html; charset=\"UTF-8\";", msg)
	} else {
		addHeader("Content-Type", "text/plain; charset=\"UTF-8\";", msg)
	}

	addHeader("From", from, msg)
	addHeader("To", strings.Join(to, ", "), msg)
	addHeader("Subject", subject, msg)

	msg.WriteString(crlf)
}

func addHeader(key, value string, msg *bytes.Buffer) {
	msg.WriteString(key + ": " + value + crlf)
}

func sendFunc(smtpServerAddr string) func(string, []string, []byte) error {
	return func(from string, to []string, msg []byte) error {
		c, err := smtp.Dial(smtpServerAddr)
		if err != nil {
			return err
		}
		defer c.Close()
		if err = c.Mail(from); err != nil {
			return err
		}
		for _, addr := range to {
			if err = c.Rcpt(addr); err != nil {
				return err
			}
		}
		w, err := c.Data()
		if err != nil {
			return err
		}
		_, err = w.Write(msg)
		if err != nil {
			return err
		}
		err = w.Close()
		if err != nil {
			return err
		}
		return c.Quit()
	}
}
