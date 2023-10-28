
dev:
	go mod tidy && go run main.go

test:
	go test ./...

build:
	gcloud builds submit --tag gcr.io/bluedot/bluedot

deploy:
	gcloud run deploy quotient \
		--image gcr.io/bluedot/bluedot \
		--platform managed

ship:
	make test && make build && make deploy