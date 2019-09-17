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

### Application Environment Variables
Use the following variables to configure the application.

| Environment Variable | Description | Example |
| --- | --- | --- |
| `ADDR` | Network interface for proxy server to listen on | `"0.0.0.0"` |
| `APP_ID` | ID of the application to reflect in the logs | `"goproygzip"` |
| `CONTENT_TYPE` | Forces the proxied `Content-Type` header to whatever you want. Useful for when the MIME type cannot automatically be detected | `"application/some-custom-format"` |
| `FLUENTD_HOST` | Hostname of the FluentD service | `"somefluentd"` |
| `FLUENTD_INIT_RETRY_COUNT` | Number of times the logger should retry reconnecting to the FluentD service | `50` |
| `FLUENTD_INIT_RETRY_INTERVAL` | Duration between the logger's attempt to connect to the FluentD service | `"somefluentd"` |
| `FLUENTD_PORT` | Port which the FluentD service is listening on | `"24224"` |
| `FORWARD_TO` | URL to forward to. When this is left empty, the service will simply be an echo server that echoes what it will sent to the next hop server if `FORWARD_TO` had been specified. | `"https://my.api.somewhere.com"` |
| `LOG_FORMAT` | Sets the logs to the format you desire for development/production. | `"text"`, or `"json"` |
| `PORT` | Port for proxy server to listen on | `"1337"` |

## Deployment

### Docker Compose
See [the example file](./deploy/docker/docker-compose.yml).

### Kubernetes
See [the example manifests](./deploy/kubernetes).

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

- GitLab CI (for GitLab hosting)
- TraviS CI (for GitHub hosting)

### Environment setup

#### Generating SSH Deploy Keys
Run the following to generate a set of keys to use for your deploy keys:

```sh
make ssh.keys
```

You can find the keys in the `./bin` directory as `id_rsa` (the private key) and `id_rsa.pub` (the public key). There should also be a Base 64 encoded version of the private key named `id_rsa.b64` which you can copy the contents and paste into the relevant `x_SSH_DEPLOY_KEY` in [the Pipeline Environment Variables table below](#pipeline-environment-variables).

#### Inserting SSH Deploy Keys - GitHub
Go to your repository and click on **Settings > Deploy keys** and hit the **Add deploy key** button on the top right of the page. Paste the contents of `id_rsa.pub` there as the key. **Note: do not base64 encode this one**.

#### Inserting SSH Deploy Keys - GitLab
Go to your repository's side menu in **Settings > Repository > Deploy Keys** and add the contents of `id_rsa.pub` there as the key. **Note: do not base64 encode this one**.

#### Pipeline Environment Variables
Before running the CI pipeline, you need to input the following build pipeline variables:

| Environment Variable | Description | Example |
| --- | --- | --- |
| `BINARY_FILENAME` | Filename for the binary | `"proxy-gzip"` |
| `DOCKER_REGISTRY_HOSTNAME` | Hostname for the Docker registry | `"docker.io"` |
| `DOCKER_REGISTRY_USERNAME` | *Optional*: Username for the Docker registry (if not present, the job will be skipped) | `"username"` |
| `DOCKER_REGSITRY_PASSWORD` | *Optional*: Password for the Docker registry (if not present, the job will be skipped) | `"password123"` |
| `DOCKER_IMAGE_NAMESPACE` | docker.io/**THIS**/image:tag | `"zephinzer"` |
| `DOCKER_IMAGE_NAME` | docker.io/namespace/**THIS**:tag | `"proxy-gzip"` |
| `GITHUB_REPOSITORY_URL` | When present, bumps the patch version and  *Optional*: SSH URL of the GitHub repository to release to, if not present, releasing to GitHub will be skipped | `"git@github.com:zephinzer/go-proxy-gzip.git"` |
| `GITHUB_SSH_DEPLOY_KEY` | Base64 encoded deploy key for the GitHub repository *Optional*: Only in play when `GITHUB_REPOSITORY_URL` is specified. | `""` |
| `GITHUB_OAUTH_TOKEN` | *Optional*: When specified, deploys the built binaries to GitHub under releases. | `""` |
| `GITLAB_REPOSITORY_URL` | *Optional*: SSH URL of the GitLab repository to release to, if not present, releasing to GitLab will be skipped | `"git@gitlab.com:zephinzer/go-proxy-gzip.git"` |
| `GITLAB_SSH_DEPLOY_KEY` | Base64 encoded deploy key for the GitHub repository *Optional*: Only in play when `GITLAB_REPOSITORY_URL` is specified. | `""` |
| `VERSION_BUMP` | *Optional*: One of "patch", "minor", or "major". Only of use if either the GitHub or GitLab URL is specified. Indicates whether the semver version bump should be a patch, minor, or major one accordingly | `"patch"` |

# TODOS

- Healthchecks to verify next hop server is alive
- Healthchecks to verify self-health 
- Distributed tracing
- ~~Logs collation~~

> (help anyone?)

# License
This project is licensed under the [MIT license](./LICENSE).

# Cheers
