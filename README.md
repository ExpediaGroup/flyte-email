# flyte-email

![Build Status](https://travis-ci.org/HotelsDotCom/flyte-email.svg?branch=master)
[![Docker Stars](https://img.shields.io/docker/stars/hotelsdotcom/flyte-email.svg)](https://hub.docker.com/r/hotelsdotcom/flyte-email)
[![Docker Pulls](https://img.shields.io/docker/pulls/hotelsdotcom/flyte-email.svg)](https://hub.docker.com/r/hotelsdotcom/flyte-email)

## Overview
The Email pack provides the ability to send emails.

## Build and Run
### Command Line
To build and run from the command line:
* Clone this repo
* Run `dep ensure` (must have [dep](https://github.com/golang/dep) installed )
* Run `go build`
* Run `FLYTE_API_URL=<URL> SMTPSERVER=<SERVER> ./flyte-email`
* Fill in this command with the relevant API url and smtp server value.




### Docker
To build and run from docker
* Run `docker build -t flyte-email .`
* Run `docker run -e FLYTE_API_URL=<URL> -e SMTPSERVER=<SERVER> flyte-email`
* All of these environment variables need to be set





## Commands
This pack provides the 'SendEmail' command, and as its name suggests, simply sends an email.
#### Input
This commands input requires the email from address, the email to address, the email subject, the email body and isHtmlEmail flag:
```
"input": {
    "from": "dave@discoveryone.com",
    "to": "hal9000@discoveryone.com",
    "subject": "Open the pod bay doors, HAL",
    "body": "...",
    "isHtmlEmail": true
    }
```
#### Output
This command can either return an `EmailSent` event meaning the email has successfully sent or an 
`SendEmailFailed` event, meaning there was a problem.
##### EmailSent
This is the success event, it contains the command name, sender, receiver, subject and email body. It returns them 
in the form:
```
"payload": {
        "from": "dave@discoveryone.com",
        "to": "hal9000@discoveryone.com",
        "subject": "Open the pod bay doors, HAL",
        "body": "...",
        "isHtmlEmail": true
}
```
##### SendEmailFailed
This contains the normal output fields plus the error if the command fails:
```
"payload": {
        "from": "dave@discoveryone.com",
        "to": "hal9000@discoveryone.com",
        "subject": "Open the pod bay doors, HAL",
        "body": "...",
        "isHtmlEmail": true,
        "error": "I'm sorry, Dave. I'm afraid I can't do that."
}
```
