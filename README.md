# Retry functions for Go

[![Release][release-image]][releases]

Simple retry functions to add exponential backoff to golang apps.

## Installation

Requires Go 1.26 or later.

```bash
go get github.com/snowplow-devops/go-retry
```

## How to use?

```golang
import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/snowplow-devops/go-retry"
)

func main() {
	// Attempt to execute a function 5 times with backoff
	err := retry.Exponential(5, time.Second, "Error Message Prefix", func() error {
		// Run your code here!
		return nil
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// Attempt to execute a function 5 times with backoff
	res, err := retry.ExponentialWithInterface(5, time.Second, "Error Message Prefix", func() (interface{}, error) {
		// Run your code here!
		return ResultType{}, nil
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// Cast the result back to the expected type
	resCast := res.(ResultType)
	logrus.Infof("Result: %v", resCast)
}
```

## How it works

Both functions call `f` and, if it returns an error, wait and call it again until it succeeds or the attempts run out.

- **Attempts** is the total number of calls, including the first. **Always set it to at least `1`.** A value of `0` or less is not meaningful: `f` is still called once, exactly as with `1`, and is never retried. Use `2` or more to get any retries.
- **Backoff** starts at the `sleep` duration and at least doubles after each failed attempt. Up to 50% random jitter is added to every wait, so many clients failing together don't all retry at the same moment.
- **No backoff** happens when `sleep` is zero or negative: failed attempts are retried immediately.
- **Logging**: every failed attempt is logged as a warning through [logrus][logrus], including the remaining attempts, the prefix and the error.
- **Errors**: when every attempt fails, the last error is returned wrapped with the prefix (`<prefix>: <error>`). The original error can still be matched with `errors.Is` or `errors.As`.
- **Results**: `ExponentialWithInterface` returns the result of the first successful call, or `nil` if every attempt fails.

## Development

The Makefile covers the common tasks:

| Command       | Description                                                           |
|---------------|-----------------------------------------------------------------------|
| `make all`    | Build the package                                                     |
| `make test`   | Run the unit tests with the race detector and print a coverage report |
| `make lint`   | Run [golangci-lint][golangci-lint] using `.golangci.yml`              |
| `make format` | Format the code with golangci-lint's formatters                       |
| `make tidy`   | Tidy `go.mod` and `go.sum`                                            |

`make lint` and `make format` run a pinned golangci-lint version through `go run`, so there is nothing to install beyond Go itself.

## Copyright and license

Copyright (c) 2021-2026 Snowplow Analytics Ltd. All rights reserved.

Licensed under the [Apache License, Version 2.0][license] (the "License");
you may not use this software except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[release-image]: https://img.shields.io/github/v/tag/snowplow-devops/go-retry?sort=semver&label=golang&color=6ad7e5&style=flat
[releases]: https://github.com/snowplow-devops/go-retry/tags
[logrus]: https://github.com/sirupsen/logrus
[golangci-lint]: https://golangci-lint.run/
[license]: LICENSE-2.0.txt
