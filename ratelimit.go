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
	// Get the bucket, and if it does not exist, create it.
	bucket, exists := s.rateLimitBuckets[endPoint.bucket]
	if !exists {
		bucket = &rateBucket{remaining: 1}
		s.rateLimitBuckets[endPoint.bucket] = bucket
	}

	// Lock this bucket
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()

	// Wait for the bucket to expire if we're out of attempts
	now := time.Now()
	if bucket.remaining == 0 && bucket.reset.After(now) {
		logger.Warnf("We are out of slots for %s, waiting...", endPoint.bucket)
		time.Sleep(bucket.reset.Sub(now))
	}

	// Once we're past the bucket lock, lock globally
	s.globalRateLimit.Lock()
	defer s.globalRateLimit.Unlock()

	// Wait for the globalRateLimit lock if we're being globally ratelimited
	now = time.Now()
	if s.globalReset.After(now) {
		logger.Warnf("We are waiting for the globalRateLimit ratelimit...")
		time.Sleep(s.globalReset.Sub(now))
	}

	// Okay, we've exhausted all possible ratelimit timers, let's send
	response, err := call()
	now = time.Now()

	// Read the headers
	var (
		headerRemaining  = response.Header.Get("X-RateLimit-Remaining")
		headerLimit      = response.Header.Get("X-RateLimit-Limit")
		headerReset      = response.Header.Get("X-RateLimit-Reset")
		headerRetryAfter = response.Header.Get("Retry-After")
		headerGlobal     = response.Header.Get("X-RateLimit-Global")
	)

	// Are we being ratelimited because of that last request?
	if response.StatusCode == 429 {
		if headerRetryAfter == "" {
			return errors.New("We are being ratelimited, but Discord didn't send a Retry-After header")
		}

		retryAfter, parseError := strconv.Atoi(headerRetryAfter)
		if parseError != nil {
			return errors.New("We are being ratelimited, but Discord didn't send a valid Retry-After header")
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

		// Automatically queue a retry, but this one will wait for the timers to expire
		return s.rateLimit(endPoint, call)
	}

	// Nope, not ratelimited, but let's update our bucket first
	var parseError error
	if headerRemaining != "" {
		bucket.remaining, parseError = strconv.Atoi(headerRemaining)
	}
	if headerLimit != "" {
		bucket.limit, parseError = strconv.Atoi(headerLimit)
	}
	if endPoint.resetTime == -1 {
		if headerReset != "" {
			var unix int64
			unix, parseError = strconv.ParseInt(headerReset, 10, 64)
			if parseError == nil {
				bucket.reset = time.Unix(unix, 0)
			}
		}
	} else {
		bucket.reset = now.Add(time.Duration(endPoint.resetTime) * time.Millisecond)
	}

	// Check for errors
	if err != nil {
		return err // If the call previously errored, return that now (we still wanted to try and read the headers
	}
	if parseError != nil {
		return parseError // Did we have any issues reading the headers?
	}

	// No errors? Awesome.
	return nil
}
