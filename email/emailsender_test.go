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
	"io/ioutil"
	"net/mail"
	"strings"
	"testing"
)

func TestEmailSender_SendHtmlEmail(t *testing.T) {
	f, r := mockSendWithRecorder(nil)

	sender := &emailSender{send: f}

	from := "flyte@email.com"
	to := []string{"dude@email.com"}
	subject := "Welcome"
	body := "Yo!!!"

	err := sender.Send(from, to, subject, body, true)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if r.from != from {
		t.Errorf("wrong from address\nwant: %s\ngot : %s\n", from, r.from)
	}
	if strings.Join(r.to, " ") != strings.Join(to, " ") {
		t.Errorf("wrong to address\nwant: %s\ngot : %s\n", to, strings.Join(r.to, " "))
	}

	m, err := mail.ReadMessage(strings.NewReader(string(r.msg)))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if m.Header.Get("Content-Type") != "text/html; charset=\"UTF-8\";" {
		t.Errorf("wrong header Content-Type\nwant: %s\ngot : %s\n", "text/html; charset=\"UTF-8\";", m.Header.Get("Content-Type"))
	}
	if m.Header.Get("From") != from {
		t.Errorf("wrong header From\nwant: %s\ngot : %s\n", from, m.Header.Get("From"))
	}
	if m.Header.Get("To") != strings.Join(to, " ") {
		t.Errorf("wrong header To\nwant: %s\ngot : %s\n", to, m.Header.Get("To"))
	}
	if m.Header.Get("Subject") != subject {
		t.Errorf("wrong header Subject\nwant: %s\ngot : %s\n", subject, m.Header.Get("Subject"))
	}

	msgBody, err := ioutil.ReadAll(m.Body)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(msgBody) != body {
		t.Errorf("wrong message body\nwant: %s\ngot : %s\n", body, string(msgBody))
	}
}

func TestEmailSender_SendPlainTextEmail(t *testing.T) {
	f, r := mockSendWithRecorder(nil)

	sender := &emailSender{send: f}

	body := "Yo!!!"

	err := sender.Send("flyte@email.com", []string{"dude@email.com"}, "Welcome", body, false)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	m, err := mail.ReadMessage(strings.NewReader(string(r.msg)))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if m.Header.Get("Content-Type") != "text/plain; charset=\"UTF-8\";" {
		t.Errorf("wrong header Content-Type\nwant: %s\ngot : %s\n", "text/plain; charset=\"UTF-8\";", m.Header.Get("Content-Type"))
	}

	msgBody, err := ioutil.ReadAll(m.Body)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(msgBody) != body {
		t.Errorf("wrong message body\nwant: %s\ngot : %s\n", body, string(msgBody))
	}
}

func mockSendWithRecorder(errToReturn error) (func(string, []string, []byte) error, *emailRecorder) {
	r := new(emailRecorder)
	return func(from string, to []string, msg []byte) error {
		*r = emailRecorder{from, to, msg}
		return errToReturn
	}, r
}

type emailRecorder struct {
	from string
	to   []string
	msg  []byte
}
