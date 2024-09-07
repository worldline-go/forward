# forward

[![License](https://img.shields.io/github/license/worldline-go/forward?color=red&style=flat-square)](https://raw.githubusercontent.com/worldline-go/forward/main/LICENSE)
[![Coverage](https://img.shields.io/sonar/coverage/worldline-go_forward?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=worldline-go_forward)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/worldline-go/forward/test.yml?branch=main&logo=github&style=flat-square&label=ci)](https://github.com/worldline-go/forward/actions)
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

```sh
sudo ./forward -s '/var/run/docker.sock:/:-POST,-PUT,-DELETE'
```

If you need multiple ports with multiple sockets than give name before prefix of `@`.

```sh
forward -H localhost@127.0.0.1:8082 -s 'localhost@/var/run/docker.sock'
```

Use `-H` and `-s` more than once to make connections.

```sh
forward -H localhost@127.0.0.1:8082 -s 'localhost@/var/run/docker.sock' -H share@0.0.0.0:8080 -s 'share@/var/run/docker.sock:/:-POST,-PUT,-DELETE'

# 2024-09-07 21:01:45 CEST INF forward version: v0.0.0 commit: - buildDate:-
# 2024-09-07 21:01:45 CEST INF localhost - route [/] to [/var/run/docker.sock]; allow: *; deny:
# 2024-09-07 21:01:45 CEST INF share - route [/] to [/var/run/docker.sock]; allow: *; deny: DELETE,POST,PUT
# 2024-09-07 21:01:45 CEST INF server share on [0.0.0.0:8080]
# 2024-09-07 21:01:45 CEST INF server localhost on [127.0.0.1:8082]
```

If program cannot access the socket, you will get `502 Bad Gateway` on every request.

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
