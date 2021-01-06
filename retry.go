//
// Copyright (c) 2021 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//

package retry

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

// Exponential provides the ability to exponentially retry the execution of a function
func Exponential(attempts int, sleep time.Duration, prefix string, f func() error) error {
	err := f()
	if err != nil {
		logrus.Warnf("Retrying func (attempts: %d): %s: %s", attempts, prefix, err)

		if attempts--; attempts > 0 {
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2
			time.Sleep(sleep)
			return Exponential(attempts, 2*sleep, prefix, f)
		}
		return errors.Wrap(err, prefix)
	}

	return nil
}

// ExponentialWithInterface provides the ability to exponentially retry the execution of a function
// and return a result from the function
func ExponentialWithInterface(attempts int, sleep time.Duration, prefix string, f func() (interface{}, error)) (interface{}, error) {
	res, err := f()
	if err != nil {
		logrus.Warnf("Retrying func (attempts: %d): %s: %s", attempts, prefix, err)

		if attempts--; attempts > 0 {
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2
			time.Sleep(sleep)
			return ExponentialWithInterface(attempts, 2*sleep, prefix, f)
		}
		return nil, errors.Wrap(err, prefix)
	}

	return res, nil
}
