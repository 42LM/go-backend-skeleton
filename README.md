# go-backend-skeleton
This repository acts as a template for a simple go backend with the focus on not using too much external dependencies.

It utilizes [Cobra](https://github.com/spf13/cobra) to create a cli command to run the server.

## Quickstart
Copy `.envrc.example`
```sh
cp .envrc.example .envrc
```
> [!TIP]
> Use [direnv](https://github.com/direnv/direnv) to load and unload environment variables depending on the current directory.

Install tools
```sh
go install tool
```

Setup DB
```sh
docker compose up -d
```

Run tests
```
go test ./...
```

Run server
```sh
go run main.go server
```

Make a request
```sh
curl localhost:8080/none
```

## Development
To install the tools needed for development just run
```sh
go install tool
```

### Mockery
Denerate mocks from interfaces with [mockery](https://vektra.github.io/mockery/latest/installation/).
```sh
mockery
```

### Build
```sh
go build -ldflags "-X go-backend-skeleton/app/cmd.version=v0.0.0" -o gbs
```
