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

package main

import (
	"github.com/HotelsDotCom/flyte-client/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"github.com/HotelsDotCom/go-logger"
	"github/ExpediaGroup/flyte-email/command"
	"github/ExpediaGroup/flyte-email/email"
	"net/url"
	"time"
)

const flyteApiUrl = "FLYTE_API_URL"
const smtpServer = "SMTPSERVER"

func main() {
	config, err := getConfig(
		value{name: flyteApiUrl, required: true},
		value{name: smtpServer, required: true},
	)
	if err != nil {
		logger.Fatalf("Email config not set: %v", err)
	}

	packDef := flyte.PackDef{
		Name: "Email",
		Commands: []flyte.Command{
			command.SendEmailCommand(email.NewEmailSender(config.values[smtpServer])),
		},
		HelpURL: createURL("https://github.com/ExpediaGroup/flyte-email/blob/master/README.md"),
	}
	pack := flyte.NewPack(packDef, client.NewClient(createURL(config.values[flyteApiUrl]), 10*time.Second))

	pack.Start()
	sleepIndefinitely()
}

func createURL(u string) *url.URL {
	url, err := url.Parse(u)
	if err != nil {
		logger.Fatalf("Cannot parse url: %s error: %s", u, err)
	}
	return url
}

func sleepIndefinitely() {
	select {}
}
