build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/alert alert/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/delete delete/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/register register/main.go
.PHONY: clean
clean:
	rm -rf ./bin
.PHONY: deploy-dev
deploy-dev: clean build
	serverless deploy -v --stage development
.PHONY: deploy-prod
deploy-prod: clean build
	serverless deploy -v --stage production
