# Development

## Local development

For interactive local development, use `make run` to build and run go-httpbin
or `make watch` to automatically re-build and re-run go-httpbin on every
change:

    make run
    make watch

By default, the server will listen on `http://127.0.0.1:8080`, but the host,
port, or any other [configuration option][config] may be overridden by
specifying the relevant environment variables:

    make run PORT=9999
    make run PORT=9999 MAX_DURATION=60s
    make watch HOST=0.0.0.0 PORT=8888

## Testing

Run `make test` to run unit tests, using `TEST_ARGS` to pass arguments through
to `go test`:

    make test
    make test TEST_ARGS="-v -race -run ^TestDelay"

### Integration tests

go-httpbin includes its own minimal WebSocket echo server implementation, and
we use the incredibly helpful [Autobahn Testsuite][] to ensure that the
implementation conforms to the spec.

These tests can be slow to run (~40 seconds on my machine), so they are not run
by default when using `make test`.

They are run automatically as part of our extended "CI" test suite, which is
run on every pull request:

    make testci

### WebSocket development

When working on the WebSocket implementation, it can also be useful to run
those integration tests directly, like so:

    make testautobahn

Use the `AUTOBAHN_CASES` var to run a specific subset of the Autobahn tests,
which may or may not include wildcards:

    make testautobahn AUTOBAHN_CASES=6.*
    make testautobahn AUTOBAHN_CASES=6.5.*
    make testautobahn AUTOBAHN_CASES=6.5.4


### Test coverage

We use [Codecov][] to measure and track test coverage as part of our continuous
integration test suite. While we strive for as much coverage as possible and
the Codecov CI check is configured with fairly strict requirements, 100% test
coverage is not an explicit goal or requirement for all contributions.

To view test coverage locally, use

    make testcover

which will run the full suite of unit and integration tests and pop open a web
browser to view coverage results.


## Linting and code style

Run `make lint` to run our suite of linters and formatters, which include
gofmt, [revive][], and [staticcheck][]:

    make lint


## OpenAPI documentation

The checked-in [`docs/openapi.yaml`](docs/openapi.yaml) file is generated from
the same route registry that serves `/openapi.json` at runtime. Regenerate it
after adding, removing, or changing a route:

    make openapi

Run `make openapi-check` to verify that the checked-in document is current
without modifying it. This check is also included in `make lint`.

### Contract test

Run the bounded [Schemathesis][] contract test against a locally started
go-httpbin instance:

    make openapi-test

The initial suite covers representative parameter-free operations. It is kept
separate from the unit test suite while the OpenAPI document gains the parameter
constraints and response details needed for exhaustive contract testing.


## Docker images

To build a docker image locally:

    make image

To build a docker image an push it to a remote repository:

    make imagepush

By default, images will be tagged as `mccutchen/go-httpbin:${COMMIT}` with the
current HEAD commit hash.

Use `VERSION` to override the tag value

    make imagepush VERSION=v1.2.3

or `DOCKER_TAG` to override the remote repo and version at once:

    make imagepush DOCKER_TAG=my-org/my-fork:v1.2.3

### Automated docker image builds

When a new release is created, the [Release][] GitHub Actions workflow
automatically builds and pushes new Docker images for both linux/amd64 and
linux/arm64 architectures.


[config]: /README.md#configuration
[revive]: https://github.com/mgechev/revive
[staticcheck]: https://staticcheck.dev/
[Release]: /.github/workflows/release.yaml
[Codecov]: https://app.codecov.io/gh/mccutchen/go-httpbin
[Autobahn Testsuite]: https://github.com/crossbario/autobahn-testsuite
[Schemathesis]: https://schemathesis.readthedocs.io/
