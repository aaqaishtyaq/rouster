# rouster

go tool to build [aaqaishtyaq/dockerfiles](https://github.com/aaqaishtyaq/dockerfiles)

## What

For building docker images with specified path and image tag. It reads `platform` file in the root of the `Dockerfile` context
to figure out the architectures for multi-arch docker image builds. It utilises `docker buildx` to do so. Eventually generating
the equivalent `docker buildx` build command.

For example, If I run the following command in the root of [my dockerfiles](https://github.com/aaqaishtyaq/dockerfiles) repo, It translates to

```diff
- rouster buildx -i base-ubuntu -t 0.0.1 dockerfiles
+ docker buildx build --load --platform=linux/amd64,linux/arm64 -t ghcr.io/aaqaishtyaq/base-debian:0.0.1 dockerfiles/base/debian
```

## Installation

`Rouster` requires Go 1.16+

```console
go install github.com/aaqaishtyaq/rouster@latest
```

## Usage

```console
% rouster buildx
Error: requires at least 1 arg(s), only received 0
Usage:
  rouster buildx [flags]

Flags:
  -d, --deadline int        Image build timeout. (default 20)
  -h, --help                help for buildx
  -i, --image stringArray   Image to be built
  -p, --push                Push the image to container registry.
  -t, --tag string          Tag to be used for tagging the image. (default "latest")

requires at least 1 arg(s), only received 0
```
