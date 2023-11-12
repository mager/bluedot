
dev:
	go mod tidy && go run main.go

test:
	go test ./...

postman:
	openapi2postmanv2 -s openapi.yaml -o collection.json

build:
	docker build -t bluedot .

publish:
	gcloud builds submit --tag gcr.io/geotory/bluedot

deploy:
	gcloud run deploy bluedot \
		--image gcr.io/geotory/bluedot \
		--platform managed \
		--port 8085 \
		--set-env-vars BLUEDOT_PGPASSWORD=$(BLUEDOT_PGPASSWORD)

ship:
	make build && make publish && make deploy