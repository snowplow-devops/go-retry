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
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// Attempt to execute a function 5 times with backoff
	res, err := retry.ExponentialWithInterface(5, time.Second, "Error Message Prefix", func() (interface{}, error) {
		// Run your code here!
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// Cast the result back to the expected type
	resCast := res.(ResultType)
}
```

[release-image]: http://img.shields.io/badge/golang-0.1.0-6ad7e5.svg?style=flat
[releases]: https://github.com/snowplow-devops/go-retry/releases/
