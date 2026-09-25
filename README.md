# Retry functions for Go

[![Release][release-image]][releases]

Simple retry functions to add exponential backoff to golang apps.

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

[release-image]: https://img.shields.io/github/v/tag/snowplow-devops/go-retry?sort=semver&label=golang&color=6ad7e5&style=flat
[releases]: https://github.com/snowplow-devops/go-retry/tags
