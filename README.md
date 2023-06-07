# go-email-lambda
This process is configured with smtp details to send emails with Go and is hosted as an AWS lambda.

## configuration
Please update relevant code in `main.go` to proper a smtp configuration.

## building
To build an executable, run:
```shell
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build main.go
```

## deploying
You'll want to zip and upload it a lambda
```shell
zip main.zip main
```
