// Copyright 2024 Grafana Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"errors"
	"time"
)

// RateLimiter controls the rate of API requests.
type RateLimiter struct {
	delay time.Duration
}

// NewRateLimiter creates a new RateLimiter with the specified delay.
func NewRateLimiter(delayInSeconds int) *RateLimiter {
	return &RateLimiter{
		delay: time.Duration(delayInSeconds) * time.Second,
	}
}

// Wait waits for the necessary amount of time before making the next request.
func (rl *RateLimiter) Wait() {
	time.Sleep(rl.delay)
}

type Config struct {
	BearerToken string
	RateLimiter *RateLimiter
	NodeIds     []int
}

var (
	errNoBearerToken = errors.New("bearer token must be specified")
)

func (c *Config) Validate() error {
	if c.BearerToken == "" {
		return errNoBearerToken
	}
	return nil
}

func NewConfig(bearerToken string, rateLimiter *RateLimiter) *Config {
	return &Config{
		BearerToken: bearerToken,
		RateLimiter: rateLimiter,
	}
}
