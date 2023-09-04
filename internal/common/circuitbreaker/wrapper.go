package circuitbreaker

import (
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type status int

const (
	CLOSED   status = 1
	OPEN     status = 2
	HALFOPEN status = 3
)

type CircuitBreaker interface {
	Attempt() error
	Fail(err error)
	Success()
}

type Impl struct {
	logger *zerolog.Logger
	mu     sync.Mutex
	// CLOSED - work!, OPEN - fail!, HALFOPEN - work until fail!
	State status
	// Length of last requests
	RecordLength int
	// Time for recovery of CB
	Timeout time.Duration

	LastAttemptedAt time.Time
	// Percentile requests after which CB opens
	Percentile float64
	// Buffer stores data about request results
	Buffer []bool
	// Pos increases for each next request, then resets to 0
	Pos int
	// Amount of successful requests in a row to go to CLOSED
	RecoveryRequests int
	// How many successful requests in HALFOPEN have already been made
	SuccessCount int
}

func NewCircuitBreaker(cfg *Config, logger *zerolog.Logger) CircuitBreaker {
	return &Impl{
		logger:           logger,
		State:            CLOSED,
		RecordLength:     cfg.RecordLength,
		Timeout:          cfg.Timeout,
		Percentile:       cfg.Percentile,
		Buffer:           make([]bool, cfg.RecordLength),
		Pos:              0,
		RecoveryRequests: cfg.RecoveryRequests,
		SuccessCount:     0,
	}
}

// for tests
func newCircuitBreaker(cfg *Config, logger *zerolog.Logger) *Impl {
	return &Impl{
		logger:           logger,
		State:            CLOSED,
		RecordLength:     cfg.RecordLength,
		Timeout:          cfg.Timeout,
		Percentile:       cfg.Percentile,
		Buffer:           make([]bool, cfg.RecordLength),
		Pos:              0,
		RecoveryRequests: cfg.RecoveryRequests,
		SuccessCount:     0,
	}
}

func (c *Impl) Attempt() error {
	c.mu.Lock()
	// only OPEN
	if c.State == OPEN {
		if elapsed := time.Since(c.LastAttemptedAt); elapsed > c.Timeout {
			c.logger.Info().Msg("switching to halfopen")
			c.State = HALFOPEN
			c.SuccessCount = 0
		} else {
			c.mu.Unlock()
			return errors.New("circuit breaker is open")
		}
		c.mu.Unlock()
	} else {
		c.mu.Unlock()
	}
	return nil
}

func (c *Impl) Fail(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Buffer[c.Pos] = err != nil
	c.Pos = (c.Pos + 1) % c.RecordLength
	// 0 f, 1 f, 2 t, 3 f, 4 f, 5t, 6t ... 36t

	// only HALFOPEN
	if c.State == HALFOPEN {
		if err != nil {
			c.logger.Info().Msg("Switching back to open state due to an error")

			c.State = OPEN
			c.LastAttemptedAt = time.Now()
			c.SuccessCount = 0 // reset counter
		} else {
			c.SuccessCount++
			if c.SuccessCount >= c.RecoveryRequests {
				c.logger.Info().Msg("Switching to closed state")
				c.Reset()
			}
		}
		return
	}

	// only CLOSED
	failureCount := 0
	for _, failed := range c.Buffer {
		if failed {
			failureCount++
		}
	}
	if float64(failureCount)/float64(c.RecordLength) >= c.Percentile {
		c.logger.Info().Msg("Switching to open state due to exceeding percentile")

		c.State = OPEN
		c.LastAttemptedAt = time.Now()
	}
}

func (c *Impl) Success() {
	c.Fail(nil)
}

func (c *Impl) Reset() {
	c.State = CLOSED
	c.Buffer = make([]bool, c.RecordLength) // reset buffer
	c.Pos = 0                               // reset position
	c.SuccessCount = 0                      // reset success counter
}
