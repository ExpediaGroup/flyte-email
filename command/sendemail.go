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
	"fmt"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"github.com/HotelsDotCom/go-logger"
	"github/ExpediaGroup/flyte-email/email"
)

var (
	EmailSentDef      = flyte.EventDef{Name: "EmailSent"}
	SendEmailErrorDef = flyte.EventDef{Name: "SendEmailFailed"}
)

type sendEmailInput struct {
	From        string   `json:"from"`
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	IsHtmlEmail bool     `json:"isHtmlEmail"`
}

func SendEmailCommand(emailSender email.EmailSender) flyte.Command {
	return flyte.Command{
		Name:    "SendEmail",
		Handler: sendEmailHandlerFunc(emailSender),
		OutputEvents: []flyte.EventDef{
			EmailSentDef,
			SendEmailErrorDef,
		},
	}
}

func sendEmailHandlerFunc(emailSender email.EmailSender) flyte.CommandHandler {
	return func(rawInput json.RawMessage) flyte.Event {
		var input sendEmailInput
		if err := json.Unmarshal(rawInput, &input); err != nil {
			err = fmt.Errorf("Could not unmarshal input for 'SendEmail' command: %v", err)
			logger.Error(err)
			return flyte.NewFatalEvent(err.Error())
		}

		err := emailSender.Send(input.From, input.To, input.Subject, input.Body, input.IsHtmlEmail)
		if err != nil {
			logger.Error(err)
			return newSendEmailFailedEvent(input, err.Error())
		}

		return newEmailSentEvent(input)
	}
}

type emailSentPayload struct {
	sendEmailInput
}

func newEmailSentEvent(commandInput sendEmailInput) flyte.Event {
	return flyte.Event{
		EventDef: EmailSentDef,
		Payload: emailSentPayload{
			commandInput,
		},
	}
}

type sendEmailFailedPayload struct {
	sendEmailInput
	Err string
}

func newSendEmailFailedEvent(commandInput sendEmailInput, err string) flyte.Event {
	return flyte.Event{
		EventDef: SendEmailErrorDef,
		Payload: sendEmailFailedPayload{
			commandInput,
			err,
		},
	}
}
