# forward

[![Codecov](https://img.shields.io/codecov/c/github/worldline-go/forward?logo=codecov&style=flat-square)](https://app.codecov.io/gh/worldline-go/forward)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/worldline-go/forward/Test?logo=github&style=flat-square&label=ci)](https://github.com/worldline-go/forward/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/worldline-go/forward.svg)](https://pkg.go.dev/github.com/worldline-go/forward)

Export socket connection to HTTP service with filter options.

## Usage

```
Flags:
  -h, --help                 help for forward
  -H, --host string          Host to listen on, default: 0.0.0.0:8080 (default "0.0.0.0:8080")
  -s, --socket stringArray   Socket to export: /var/run/docker.sock:/docker:*,-POST,-PUT,-DELETE
  -v, --version              version for forward
```

Show the lists of sockets to export with show http methods.

HTTP methods could be any string(case insensitive), and starting with `-` will not allowed.

```sh
# allow all except POST, PUT, DELETE
/var/run/docker.sock:/docker/:-POST,-PUT,-DELETE

# Only allow GET requests
/var/run/docker.sock:/docker/:GET

# Disable all methods of requests
/var/run/docker.sock:/docker/:-*
```

If program cannot access the socket, you will get `502 Bad Gateway` on every request.

```sh
sudo ./forward -s '/var/run/docker.sock:/docker/:-POST,-PUT,-DELETE'
```

### Docker

There are scratch and alpine version of container image.

```sh
docker run -p 8080:8080 -v /var/run/docker.sock:/docker.sock ghcr.io/worldline-go/forward -s /docker.sock:/docker/:-POST,-PUT,-DELETE,-PATCH
```

## Development

Local generate binary and docker image

```sh
goreleaser release --snapshot --rm-dist
```
