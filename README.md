# Redirector

[![Build Status](https://travis-ci.org/c0va23/redirector.svg?branch=master)](https://travis-ci.org/c0va23/redirector)

## About `redirector`

`redirector` is dynamic configurable HTTP-redirection service.

Supported store engine:
- `memory` (not store configuration between restarts)
- `redis` (now support only single-server configuration)

Supported match / resolve algorithm:
- `simple` (full match of host and path, not change target)
- `pattern` (full match host, match path by regexp, replace target placeholders by pattern groups)


### Use cases

- Short links with full path name control
- Short links with activation codes (promo codes, confirmation of email addresses, password recovery, etc.)


### How it work?

1. Run service on your server(s) ([See usage](#Usage))
2. Configure service via [API](/api.yml) (Web UI expected in the future)
3. Direct the DNS records of the domains used by the redirector to the server where the service is started
4. Test by open any configured link in browser or with `curl`


### Algorithm matching rules and resolving targets

1. Find host by full match.
2. If host not found, then return empty response with HTTP code 404.
3. Find rule with selected resolve algorithm.
4. If rule not found, then redirect to default target url and HTTP code.
5. If rule found, then redirect to rule URL and HTTP code.


### Pattern resolver example

If we have: source `^/(.)$` and target path `https://example.org/promocode/{0}`.

Then request to path `/PROMOCODE` resolved to `https://example.org/promocode/PROMOCODE`.


### Ready for production use?

**Not yet used in production.**


## Test and build

### Requirements

- `go` (tested on version 1.10)
- `make`

### Install dependencies

```bash
# Install development dependencies
make dev-deps

# Install runtime dependencies
make deps

# Generate swagger files
make gen-swagger
```


### Run linters and tests

```bash
make lint
make run-test
```


### Build server binary

```bash
make bin/redirector-server
```


## Build docker image and run container

```bash
# Build image
docker build -t redirector .

# Run container
docker run -p 8080:8080 redirector
```


## Usage

Show usage with command:
```bash
./bin/redirector-server --help
```

```
Usage:
  redirector-server [OPTIONS]

Redirector configure API

Application Options:
      --scheme=                   the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=          grace period for which to wait before shutting down the server (default: 10s)
      --max-header-size=          controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of
                                  the request body. (default: 1MiB)
      --socket-path=              the unix socket to listen on (default: /var/run/redirector.sock)
      --host=                     the IP to listen on (default: localhost) [$HOST]
      --port=                     the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=             limit the number of outstanding requests
      --keep-alive=               sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default: 3m)
      --read-timeout=             maximum duration before timing out read of the request (default: 30s)
      --write-timeout=            maximum duration before timing out write of the response (default: 60s)
      --tls-host=                 the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=                 the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=          the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=                  the private key to use for secure conections [$TLS_PRIVATE_KEY]
      --tls-ca=                   the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=         limit the number of outstanding requests
      --tls-keep-alive=           sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
      --tls-read-timeout=         maximum duration before timing out read of the request
      --tls-write-timeout=        maximum duration before timing out write of the response

store:
  -s, --store-type=[memory|redis] Type of store engine (default: memory)
      --redis-uri=                Connection URI for Redis. Required if store-type equal redis
      --redis-pool-size=          Redis pool size (default: 10)
  -u, --basic-username=           Username for Basic auth [$BASIC_USERNAME]
  -p, --basic-password=           Password for Basic auth. [$BASIC_PASSWORD]

Help Options:
  -h, --help                      Show this help message
```

### Logger level

Logger level can be configured with ennvar LOG_LEVEL.

Allowed values:
- debug
- info
- warn
- error
- fatal
- panic
