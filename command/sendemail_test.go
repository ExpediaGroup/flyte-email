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

package command

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestSendEmailCommand_shouldSendPlainTextEmail(t *testing.T) {
	f, r := mockSendWithRecorder(nil)

	input := sendEmailInput{
		From:        "flyte@email.com",
		To:          []string{"dude@email.com"},
		Subject:     "Welcome",
		Body:        "Yo!!!",
		IsHtmlEmail: true,
	}

	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	got := SendEmailCommand(&mockEmailSender{send: f}).Handler(b)

	if got.EventDef.Name != EmailSentDef.Name {
		t.Fatalf("wrong event\nwant: %s\ngot : %s\n", EmailSentDef.Name, got.EventDef.Name)
	}

	p := got.Payload.(emailSentPayload)
	if p.From != input.From {
		t.Errorf("wrong From in the event payload\nwant: %s\ngot : %s\n", input.From, p.From)
	}
	if strings.Join(p.To, " ") != strings.Join(input.To, " ") {
		t.Errorf("wrong To in the event payload\nwant: %s\ngot : %s\n", input.To, p.To)
	}
	if p.Subject != input.Subject {
		t.Errorf("wrong subject in the event payload\nwant: %s\ngot : %s\n", input.Subject, p.Subject)
	}
	if p.Body != input.Body {
		t.Errorf("wrong body in the event payload\nwant: %s\ngot : %s\n", input.Body, p.Body)
	}

	if r.from != input.From {
		t.Errorf("wrong from address\nwant: %s\ngot : %s\n", input.From, r.from)
	}
	if strings.Join(r.to, " ") != strings.Join(input.To, " ") {
		t.Errorf("wrong to address\nwant: %s\ngot : %s\n", input.To, r.to)
	}
	if r.subject != input.Subject {
		t.Errorf("wrong email subject\nwant: %s\ngot : %s\n", input.Subject, r.subject)
	}
	if r.body != input.Body {
		t.Errorf("wrong email body\nwant: %s\ngot : %s\n", input.Body, r.body)
	}
	if r.isHtmlEmail != input.IsHtmlEmail {
		t.Errorf("wrong email format\nwant: %t\ngot : %t\n", input.IsHtmlEmail, r.isHtmlEmail)
	}
}

func TestSendEmailCommand_shouldReturnSendEmailErrorEvent_whenErrorHappensWhileSendingEmail(t *testing.T) {
	errToReturn := errors.New("Uhhhh something went wrong!!!")

	input := sendEmailInput{
		From:    "flyte@email.com",
		To:      []string{"dude@email.com"},
		Subject: "Welcome",
		Body:    "Yo!!!",
	}

	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	got := SendEmailCommand(&mockEmailSender{send: mockSend(errToReturn)}).Handler(b)

	if got.EventDef.Name != SendEmailErrorDef.Name {
		t.Fatalf("wrong event\nwant: %s\ngot : %s\n", SendEmailErrorDef.Name, got.EventDef.Name)
	}

	p := got.Payload.(sendEmailFailedPayload)
	if p.From != input.From {
		t.Errorf("wrong From in the event payload\nwant: %s\ngot : %s\n", input.From, p.From)
	}
	if strings.Join(p.To, " ") != strings.Join(input.To, " ") {
		t.Errorf("wrong To in the event payload\nwant: %s\ngot : %s\n", input.To, p.To)
	}
	if p.Subject != input.Subject {
		t.Errorf("wrong subject in the event payload\nwant: %s\ngot : %s\n", input.Subject, p.Subject)
	}
	if p.Body != input.Body {
		t.Errorf("wrong body in the event payload\nwant: %s\ngot : %s\n", input.Body, p.Body)
	}
	if p.Err != errToReturn.Error() {
		t.Errorf("wrong body in the event payload\nwant: %s\ngot : %s\n", errToReturn.Error(), p.Err)
	}
}

func TestSendEmailCommand_shouldReturnFatalEvent_whenUnmarshalErrorHappens(t *testing.T) {
	got := SendEmailCommand(&mockEmailSender{send: mockSend(nil)}).Handler(json.RawMessage(`{"dodgy-json}`))

	if got.EventDef.Name != "FATAL" {
		t.Fatalf("wrong event\nwant: %s\ngot : %s\n", "FATAL", got.EventDef.Name)
	}
}

type mockEmailSender struct {
	send func(from string, to []string, subject, body string, isHtmlEmail bool) error
}

func (e mockEmailSender) Send(from string, to []string, subject, body string, isHtmlEmail bool) error {
	return e.send(from, to, subject, body, isHtmlEmail)
}

func mockSend(errToReturn error) func(string, []string, string, string, bool) error {
	return func(from string, to []string, subject, body string, isHtmlEmail bool) error {
		return errToReturn
	}
}

func mockSendWithRecorder(errToReturn error) (func(string, []string, string, string, bool) error, *emailRecorder) {
	r := new(emailRecorder)
	return func(from string, to []string, subject, body string, isHtmlEmail bool) error {
		*r = emailRecorder{from, to, subject, body, isHtmlEmail}
		return errToReturn
	}, r
}

type emailRecorder struct {
	from        string
	to          []string
	subject     string
	body        string
	isHtmlEmail bool
}
