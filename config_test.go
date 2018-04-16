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
	"testing"
)

func TestShouldReturnDefaultConfigWhenEnvironmentVariablesNotSet(t *testing.T) {
	env = &mockEnvironment{} // does not contain any values, thus defaults if set will be used

	env, _ := getConfig(
		value{name: "HOST", dfault: "https://something-default.com", required: true},
		value{name: "SERVICE_ACCOUNT", dfault: "ServiceAccountDefault", required: false},
		value{name: "SERVICE_ACCOUNT_PASS", required: false}, // note: no default set
	)

	if env.values["HOST"] != "https://something-default.com" {
		t.Errorf("HOST is %s, should be https://something-default.com", env.values["HOST"])
	}
	if env.values["SERVICE_ACCOUNT"] != "ServiceAccountDefault" {
		t.Errorf("SERVICE_ACCOUNT is %s, should be ServiceAccountDefault", env.values["SERVICE_ACCOUNT"])
	}
	if env.values["SERVICE_ACCOUNT_PASS"] != "" {
		t.Error("SERVICE_ACCOUNT_PASS should be blank.")
	}
}

func TestShouldReturnEnvironmentConfigWhenSet(t *testing.T) {
	env = &mockEnvironment{
		values: map[string]string{
			"HOST":            "https://www.something-set-in-environment.com",
			"SERVICE_ACCOUNT": "ServiceAccountFromEnvironment",
		},
	}

	env, _ := getConfig(
		value{name: "HOST", dfault: "https://something-default.com", required: true},
		value{name: "SERVICE_ACCOUNT", dfault: "ServiceAccountDefault", required: false},
	)

	if env.values["HOST"] != "https://www.something-set-in-environment.com" {
		t.Errorf("HOST is %s, should be https://www.something-set-in-environment.com", env.values["HOST"])
	}
	if env.values["SERVICE_ACCOUNT"] != "ServiceAccountFromEnvironment" {
		t.Errorf("SERVICE_ACCOUNT is %s, should be ServiceAccountFromEnvironment", env.values["SERVICE_ACCOUNT"])
	}
}

func TestShouldReturnDefaultConfigIfEnvironmentVariableIsBlank(t *testing.T) {
	env = &mockEnvironment{
		values: map[string]string{
			"SERVICE_ACCOUNT_PASS": "", // note key exists, but not value
		},
	}

	env, _ := getConfig(
		value{name: "SERVICE_ACCOUNT_PASS", dfault: "work1234", required: false}, // default set
	)

	if env.values["SERVICE_ACCOUNT_PASS"] != "work1234" {
		t.Errorf("SERVICE_ACCOUNT_PASS is %s, should be work1234", env.values["SERVICE_ACCOUNT_PASS"])
	}
}

func TestShouldReturnErrorIfValueIsRequiredButNotSet(t *testing.T) {
	env = &mockEnvironment{
		values: map[string]string{
			"HOST": "", // not set in environment variables...
		},
	}

	_, err := getConfig(
		value{name: "HOST", required: true}, // not set with a default value, but required...
	)

	if err == nil || err.Error() != "HOST is required!!" {
		t.Error("HOST error should've been thrown!!!")
	}
}

type mockEnvironment struct {
	values map[string]string
}

func (m *mockEnvironment) getValueFor(name string) string {
	return m.values[name]
}
