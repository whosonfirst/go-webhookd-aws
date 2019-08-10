tools:
	go build -mod vendor -o bin/webhookd-lambda cmd/webhookd-lambda/main.go
	go build -mod vendor -o bin/webhookd-lambda-task cmd/webhookd-lambda-task/main.go
	go build -mod vendor -o bin/webhookd-config cmd/webhookd-config/main.go
	go build -mod vendor -o bin/webhookd-flatten-config cmd/webhookd-flatten-config/main.go

lambda: lambda-webhookd lambda-task

lambda-webhookd:
	if test -f main; then rm -f main; fi
	if test -f webhookd.zip; then rm -f webhookd.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/webhookd-lambda/main.go
	zip webhookd.zip main
	rm -f main

lambda-task:
	if test -f main; then rm -f main; fi
	if test -f webhookd-task.zip; then rm -f webhookd-task.zip; fi
	GOOS=linux go build -mod vendor -o main cmd/webhookd-lambda-task/main.go
	zip webhookd-task.zip main
	rm -f main
