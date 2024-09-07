# forward

[![License](https://img.shields.io/github/license/worldline-go/forward?color=red&style=flat-square)](https://raw.githubusercontent.com/worldline-go/forward/main/LICENSE)
[![Coverage](https://img.shields.io/sonar/coverage/worldline-go_forward?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=worldline-go_forward)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/worldline-go/forward/Test?logo=github&style=flat-square&label=ci)](https://github.com/worldline-go/forward/actions)
[![Go PKG](https://raw.githubusercontent.com/worldline-go/guide/main/badge/custom/reference.svg)](https://pkg.go.dev/github.com/worldline-go/forward)

Export socket connection to HTTP service with filter options.

## Usage

Use `LOG_LEVEL` and `LOG_PRETTY` env values to control the log level and format.

```
Flags:
  -h, --help                 help for forward
  -H, --host stringArray     Host to listen on (default [0.0.0.0:8080])
  -s, --socket stringArray   Socket to export: /var/run/docker.sock:/:*,-POST,-PUT,-DELETE
  -v, --version              version for forward
```

Show the lists of sockets to export with show http methods.

HTTP methods could be any string(case insensitive), and starting with `-` will not allowed.

```sh
# allow all methods to / path
/var/run/docker.sock

# allow all except POST, PUT, DELETE with basepath /docker/
/var/run/docker.sock:/docker/:-POST,-PUT,-DELETE

# Only allow GET requests
/var/run/docker.sock:/docker/:GET

# Disable all methods of requests
/var/run/docker.sock:/docker/:-*
```

If program cannot access the socket, you will get `502 Bad Gateway` on every request.

```sh
sudo ./forward -s '/var/run/docker.sock:/:-POST,-PUT,-DELETE'
```

### Docker

There are scratch and alpine version of container image.

```sh
docker run --rm -it -p 8080:8080 -v /var/run/docker.sock:/docker.sock ghcr.io/worldline-go/forward -s /docker.sock:/:-POST,-PUT,-DELETE,-PATCH
```

## Development

Local generate binary

```sh
make build
```
