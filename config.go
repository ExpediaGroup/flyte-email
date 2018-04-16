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
	"fmt"
	"os"
)

var env environment = &osEnvironment{}

type environment interface {
	getValueFor(name string) string
}

type osEnvironment struct{}

func (o *osEnvironment) getValueFor(name string) string {
	return os.Getenv(name)
}

type config struct {
	values map[string]string
}

type value struct {
	name     string
	dfault   string
	required bool
}

func getConfig(values ...value) (config, error) {
	config := config{
		values: make(map[string]string),
	}
	for _, v := range values {
		configValue := getConfigValue(v)
		if v.required && len(configValue) == 0 {
			return config, fmt.Errorf("%s is required!!", v.name)
		}
		config.values[v.name] = configValue
	}
	return config, nil
}

func getConfigValue(v value) string {
	configValue := env.getValueFor(v.name)
	if len(configValue) == 0 {
		configValue = v.dfault
	}
	return configValue
}
