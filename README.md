# ProxyGzip
A proxy service that assists in compressing your requests via gzip.

> Note: As of 22nd February 2019, this component is actively maintained but inactively upgraded

# Why this exists

Some policies enforce analysis of ingoing and outgoing through decoding of body data. A method needed to be devised to re-encode this data upon exiting the internal network before hitting the final destination.

This is it.

# Usage

## How it works
This server takes in any request and transforms the request so that:
1. The body data is compressed with `gzip`
1. The `Content-Length` header will be set to the length of the `gzip`ped content
1. The `Content-Encoding` header will be set to `"gzip"`
1. The `Content-Type` header will be derived from your data, or if the `CONTENT_TYPE` environment variable is set, that
1. All other headers sent are re-constructed from the original client request and forwarded to the next hop server **AS-IS**

The transformed request is then forwarded to the next hop server and a response is received. This response is processed such that:
1. The next hop response headers will be reflected in the response to the client **AS-IS**
1. The client request headers will be reflected in the response headers prefixed with `pgz-request-*`
1. The proxy server configuration will be reflected in the response headers prefixed with `pgz-config-*`

## Configuration

| Environment Variable | Description | Example |
| --- | --- | --- |
| `ADDR` | Network interface to listen on | `"0.0.0.0"` |
| `PORT` | Port to listen on | `"1337"` |
| `FORWARD_TO` | URL to forward to. When this is left empty, the service will simply be an echo server that echoes what it will sent to the next hop server if `FORWARD_TO` had been specified. | `"https://my.api.somewhere.com"` |
| `CONTENT_TYPE` | Forces the proxied `Content-Type` header to whatever you want. Useful for when the MIME type cannot automatically be detected | `"application/some-custom-format"` |
| `LOG_FORMAT` | Sets the logs to the format you desire for development/production. | `"text"`, or `"json"` |

## Deployment

### Docker Compose
See [the example file](./example/deployment/docker-compose.yml).

## Response Interpretation

### Status Code
The status code will reflect the status code of the next hop response.

### Headers
Headers prefixed with `pgz-config-` display the configuration of the ProxyGzip server.

Headers prefixed with `pgz-request-` display headers from the client's request.

Headers not prefixed are headers from the next hop server.

# Contributing
- Development tooling is performed via the Makefile.
- Dependency management is done via Go Modules (from Go 1.11)
- Non-collabs: Fork, make changes on your fork's `master` branch, and then issue a pull request back for changes
- Collabs: Make changes on a branch other than `master`, issue a merge request
- For quickening the process/if this repo seems dead, ping [the maintainers](./MAINTAINERS)

## Running locally

```sh
make
```

## (Manually) Getting dependencies in

> The `make start`/`make` recipe already does this for you before running the application

```sh
make deps
```

## Running tests

```sh
make test
```

## Binary generation (compilation)

```sh
# compiles for all operating systems and architectures
make compile

# for windows
make compile.windows

# for macos
make compile.macos

# for linux
make compile.linux
```

## Bundling Docker image for production

```sh
make package
```

## Publish Docker image

```sh
make release
```

## Continuous Integration Configuration

### Available pipelines
Currently, there exists integrations for:

- GitLab CI

### Environment setup

| Environment Variable | Description | Example |
| --- | --- | --- |
| `BINARY_FILENAME` | Filename for the binary | `"proxy-gzip"` |
| `DOCKER_REGISTRY_HOSTNAME` | Hostname for the Docker registry | `"docker.io"` |
| `DOCKER_IMAGE_NAMESPACE` | docker.io/**THIS**/image:tag | `"zephinzer"` |
| `DOCKER_IMAGE_NAME` | docker.io/namespace/**THIS**:tag | `"proxy-gzip"` |
| `GITHUB_REPOSITORY_URL` | *Optional*: SSH URL of the GitHub repository to release to | `"git@github.com:zephinzer/go-proxy-gzip.git"` |

# TODOS

- Healthchecks to verify next hop server is alive
- Healthchecks to verify self-health 
- Distributed tracing
- Logs collation

> (help anyone?)

# License
This project is licensed under the [MIT license](./LICENSE).

# Cheers
