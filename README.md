# bluedot

![Coverage](https://img.shields.io/badge/Coverage-7.6%25-red)

## Development

Runs on port 8085.

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
make postman
```

### Run tests and coverage badge

```sh
make test
make cover
```
