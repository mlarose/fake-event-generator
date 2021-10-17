homework-event-generator
------------------------

# Requirements

The provided instructions require the installation of these software installations on your local environment

- go 1.7+
- git
- github cli
- docker for desktop


# Building distribution binaries

A [Dockerfile](./Dockerfile) is provided to simplify the compilation of cross-platform binaries for distribution. The
following commands will create the binaries in the local `./dist` folder when executed at the root of this repository.

## Linux or OS X

```shell
mkdir -p ./dist
docker build -t builder .
docker run -it --rm -v $(pwd)\dist:/export builder
```

## Windows

```shell
mkdir ./dist
docker build -t builder .
docker run -it --rm -v %CD%\dist:/export builder
```
