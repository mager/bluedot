
dev:
	go mod tidy && go run main.go

test:
	go test ./...

cover:
	go install github.com/AlexBeauchemin/gobadge@latest
	go test ./... -covermode=count -coverprofile=coverage.out fmt
	go tool cover -func=coverage.out -o=coverage.out
	$(HOME)/go/bin/gobadge -filename=coverage.out

postman:
	openapi2postmanv2 -s openapi.yaml -o collection.json

publish:
	gcloud builds submit --tag gcr.io/geotory/bluedot

deploy:
	gcloud run deploy bluedot \
		--image gcr.io/geotory/bluedot \
		--platform managed \
		--port 8085 \
		--set-env-vars BLUEDOT_PGPASSWORD=$(BLUEDOT_PGPASSWORD)

ship:
	make publish && make deploy

openapi:
	swag init  --parseDependency --parseInternal