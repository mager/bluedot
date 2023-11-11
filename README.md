# bluedot

## Development

### Run locally

```sh
make dev
```

### Build and ship

This command builds, publishes, and deploys:

```sh
make ship
```

#### Build Dockerfile

```sh
make build
```

#### Upload build to Cloud Build

```sh
make publish
```

#### Deploy latest build to production

```sh
make deploy
```

### Generate Postman collection

```sh
openapi2postmanv2 -s openapi.yaml -o collection.json
```
