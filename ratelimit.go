package disgo

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/slf4go/logger"
)

type rateBucket struct {
	remaining int
	limit     int
	reset     time.Time
	mutex     sync.Mutex
}

func (s *Session) rateLimit(endPoint EndPoint, call func() (*http.Response, error)) error {
	return s.rateLimitRecursive(endPoint, call, false)
}

func (s *Session) rateLimitRecursive(endPoint EndPoint, call func() (*http.Response, error), recursive bool) error {
	// Get the bucket, and if it does not exist, create it.
	bucket, exists := s.rateLimitBuckets[endPoint.bucket]
	if !exists {
		bucket = &rateBucket{remaining: 1}
		s.rateLimitBuckets[endPoint.bucket] = bucket
	}

	// Lock this bucket
	if !recursive {
		bucket.mutex.Lock()
		defer bucket.mutex.Unlock()
	}

	// Wait for the bucket to expire if we're out of attempts
	now := time.Now()
	if bucket.remaining == 0 && bucket.reset.After(now) {
		logger.Warnf("We are out of slots for %s, waiting...", endPoint.bucket)
		time.Sleep(bucket.reset.Sub(now))
	}

	// Once we're past the bucket lock, lock globally
	if !recursive {
		s.globalRateLimit.Lock()
		defer s.globalRateLimit.Unlock()
	}

	// Wait for the globalRateLimit lock if we're being globally rate limited
	now = time.Now()
	if s.globalReset.After(now) {
		logger.Warnf("We are waiting for the globalRateLimit...")
		time.Sleep(s.globalReset.Sub(now))
	}

	// Okay, we've exhausted all possible ratelimit timers, let's send
	response, err := call()
	now = time.Now()

	// Not all returned error status codes include rate limit headers, if that is the case we return now
	if err != nil {
		switch response.StatusCode {
		case http.StatusBadRequest:
			fallthrough
		case http.StatusUnauthorized:
			fallthrough
		case http.StatusForbidden:
			fallthrough
		case http.StatusMethodNotAllowed:
			return err
		}

		if response.StatusCode > 500 {
			return err
		}
	}

	// Read the headers
	var (
		headerDiscordTime = response.Header.Get("Date")
		headerRemaining   = response.Header.Get("X-RateLimit-Remaining")
		headerLimit       = response.Header.Get("X-RateLimit-Limit")
		headerReset       = response.Header.Get("X-RateLimit-Reset")
		headerRetryAfter  = response.Header.Get("Retry-After")
		headerGlobal      = response.Header.Get("X-RateLimit-Global")
	)

	logger.Tracef("Ratelimit headers: remaining: %s, limit: %s, reset: %s, retryAfter: %s", headerRemaining, headerLimit, headerReset, headerRetryAfter)

	// Are we being rate limited because of that last request?
	if response.StatusCode == http.StatusTooManyRequests {
		if headerRetryAfter == "" {
			return errors.New("We are being ratelimited, but Discord didn't send a Retry-After header")
		}

		retryAfter, err := strconv.Atoi(headerRetryAfter)
		if err != nil {
			return err
		}

		resetTime := now.Add(time.Duration(retryAfter) * time.Millisecond)

		if headerGlobal == "true" {
			logger.Error("We are being globally ratelimited!")
			s.globalReset = resetTime
		} else {
			logger.Errorf("We are being ratelimited on %s!", endPoint.bucket)
			bucket.reset = resetTime
			bucket.remaining = 0
		}

		return s.rateLimitRecursive(endPoint, call, true)
	}

	// Nope, not rate limited, but let's update our bucket first
	var parseError error
	if headerRemaining != "" {
		bucket.remaining, parseError = strconv.Atoi(headerRemaining)
		if parseError != nil {
			return parseError
		}
	}
	if headerLimit != "" {
		bucket.limit, parseError = strconv.Atoi(headerLimit)
		if parseError != nil {
			return parseError
		}
	}
	if endPoint.resetTime == -1 {
		if headerReset != "" {
			unix, parseError := strconv.ParseInt(headerReset, 10, 64)
			if parseError != nil {
				return parseError
			}
			resetTime := time.Unix(unix, 0)

			if headerDiscordTime == "" {
				bucket.reset = resetTime
			} else {
				discordTime, parseError := time.Parse(time.RFC1123, headerDiscordTime)
				if parseError != nil {
					return parseError
				}
				bucket.reset = now.Add(resetTime.Sub(discordTime))
			}
		}
	} else {
		bucket.reset = now.Add(time.Duration(endPoint.resetTime) * time.Millisecond)
	}

	// If we did encounter a http response that did include ratelimit headers, we return here after reading the headers
	return err
}
