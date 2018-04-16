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

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type logger struct {
	calldepth int
	afterLog  func(msg string)
	log       *log.Logger
}

func newLogger(prefix string, afterLog func(msg string), writer io.Writer) logger {

	prefix = fmt.Sprintf("[%s] ", prefix)
	return logger{
		calldepth: 3,
		afterLog:  afterLog,
		log:       log.New(writer, prefix, log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (logger logger) Log(v ...interface{}) {

	msg := fmt.Sprint(v...)
	logger.log.Output(logger.calldepth, msg)
	logger.afterLog(msg)
}

func (logger logger) Logf(format string, v ...interface{}) {

	msg := fmt.Sprintf(format, v...)
	logger.log.Output(logger.calldepth, msg)
	logger.afterLog(msg)
}

var (
	Debug  = func(...interface{}) {}
	Debugf = func(string, ...interface{}) {}
	Info   = func(...interface{}) {}
	Infof  = func(string, ...interface{}) {}
	Error  = func(...interface{}) {}
	Errorf = func(string, ...interface{}) {}
	Fatal  = func(...interface{}) {}
	Fatalf = func(string, ...interface{}) {}
)

var logLevels = map[string]bool{
	"DEBUG": true,
	"INFO":  true,
	"ERROR": true,
	"FATAL": true,
}

// --- logger initialization ---

var (
	LogLevel  string
	LogWriter io.Writer
	AfterFatal func(string)
)

func init() {
	LogLevel = getLogLevel()
	LogWriter = getWriter()
	AfterFatal = func(string) { os.Exit(1) }
	InitLoggers()
}

// Use only if you set LogLevel, LogWriter and/or AfterFatal to different values and need to re-initialize logger.
func InitLoggers() {

	resetLoggers()
	log.SetOutput(LogWriter)

	if _, ok := logLevels[LogLevel]; !ok {
		if LogLevel != "" {
			log.Printf("invalid log level supplied %q", LogLevel)
		}
		log.Print("log level can be set by using env variable LOGLEVEL=DEBUG|INFO|ERROR")
		LogLevel = "INFO"
	}
	log.Printf("log level set to: %s", LogLevel)

	switch LogLevel {
	case "DEBUG":
		l := newLogger("DEBUG", func(string) {}, LogWriter)
		Debug = l.Log
		Debugf = l.Logf
		fallthrough
	case "INFO":
		l := newLogger("INFO", func(string) {}, LogWriter)
		Info = l.Log
		Infof = l.Logf
		fallthrough
	case "ERROR":
		l := newLogger("ERROR", func(string) {}, LogWriter)
		Error = l.Log
		Errorf = l.Logf
		fallthrough
	case "FATAL":
		l := newLogger("FATAL", AfterFatal, LogWriter)
		Fatal = l.Log
		Fatalf = l.Logf
	}
}

// sets loggers to dummy functions that do nothing
func resetLoggers() {

	var nilLog = func(...interface{}) {}
	var nilLogf = func(string, ...interface{}) {}

	Debug = nilLog
	Debugf = nilLogf
	Info = nilLog
	Infof = nilLogf
	Error = nilLog
	Errorf = nilLogf
	Fatal = nilLog
	Fatalf = nilLogf
}

func getWriter() io.Writer {

	logFile := getLogFile()
	if logFile == "" {
		log.Println("LOGFILE not set, logging to stdout only")
		return os.Stdout
	}

	log.Printf("log file set to: %s", logFile)
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Fatalf("cannot create log file: %v", err)
	}
	return io.MultiWriter(os.Stdout, file)
}

func getLogLevel() string {
	return os.Getenv("LOGLEVEL")
}

func getLogFile() string {
	return os.Getenv("LOGFILE")
}
