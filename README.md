fake-event-generator
------------------------

This homebrew project was a proof of concept that merely devolved into a prop for homework assignments.

# Usage

## event-gen

The `event-gen` script supports three modes of output:

- `event-gen stdout` will print out to the standard output.
- `event-gen tcp` will connect to a tcp server and send events as json documents separated by newlines.
- `event-gen http` will connect to a http server and send events as individual documents in HTTP POST requests.

Both commands support the following arguments

	-H, --host string   remote host address (default "localhost")
	-p, --port uint16   remote port (default 3333)

Global flags include:

	-seed int	initialize random number generation with this seed

## tcp-server and http-server

Both utilities are provided only as code example for receiving output from event-gen and printing it to the stdout.


# Developer guide

## Requirements

The provided instructions require the installation of these software installations on your local environment

- go 1.17+
- docker for desktop

## Testing and running from sources with golang

```shell
go test ./...
go run main.go help
```

## Building distribution binaries with docker

A [Dockerfile](./Dockerfile) is provided to simplify the compilation of cross-platform binaries for distribution. The
following commands will create the binaries in the local `./dist` folder when executed at the root of this repository.

### Linux or OS X

```shell
mkdir -p ./dist
docker build -t builder .
docker run -it --rm -v $(pwd)\dist:/export builder
```

### Windows

```shell
mkdir ./dist
docker build -t builder .
docker run -it --rm -v %CD%\dist:/export builder
```
