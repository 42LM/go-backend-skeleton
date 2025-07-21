# go-backend-skeleton
This repository acts as a template for a simple go backend with the focus on not using too much external dependencies.

It utilizes:
* [spf13/cobra](https://github.com/spf13/cobra) to create a cli command to run the server
* [vektra/mockery](https://github.com/vektra/mockery) to generate mocks from interfaces
* [stretchr/testify](https://github.com/stretchr/testify) for simpler testing in general
* [42lm/muxify](https://github.com/42LM/muxify) for a better handling of the default `*http.ServeMux`

## Development
Copy `.envrc.example`
```sh
cp .envrc.example .envrc
```
> [!TIP]
> Use a tool like [direnv](https://github.com/direnv/direnv) to load and unload environment variables depending on the current directory.

Install tools (currently only mockery)
```sh
go install tool
```

Setup DB (currently only dynamodb)
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
curl localhost:8080/v1/none
```

```sh
curl localhost:8080/v1/msg/777
```

```sh
curl -X POST localhost:8080/rpc/msg/777 -d '{"msg":"test-msg"}'
```

### Generate mocks
Generate mocks from interfaces with [mockery](https://vektra.github.io/mockery/latest/installation/).
```sh
mockery
```

### Generate protobuf
Generate mocks from interfaces with [mockery](https://vektra.github.io/mockery/latest/installation/).
```sh
go generate ./...
```

### Build
```sh
go build -ldflags "-X go-backend-skeleton/app/cmd.version=v0.0.0" -o gbs
```
