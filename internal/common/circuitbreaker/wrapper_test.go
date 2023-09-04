package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"gitlab.karlson.dev/individual/pet_gonote/note_service/internal/domain"

	"github.com/rs/zerolog"
)

func TestNewCircuitBreaker(t *testing.T) {
	domain.TestIsUnit(t)
	logger := zerolog.New(nil)
	cfg := &Config{
		RecordLength:     10,
		Timeout:          2 * time.Second,
		Percentile:       0.6,
		RecoveryRequests: 5,
	}

	//nolint:go-critic
	cb := newCircuitBreaker(cfg, &logger)

	if cb.State != CLOSED {
		t.Errorf("Expected State to be CLOSED, got %v", cb.State)
	}
}

func TestAttemptWithClosedState(t *testing.T) {
	domain.TestIsUnit(t)
	logger := zerolog.New(nil)
	cfg := &Config{Timeout: 1 * time.Second}
	cb := newCircuitBreaker(cfg, &logger)

	err := cb.Attempt()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAttemptWithOpenState(t *testing.T) {
	domain.TestIsUnit(t)
	logger := zerolog.New(nil)
	cfg := &Config{Timeout: 1 * time.Second}
	cb := newCircuitBreaker(cfg, &logger)
	cb.State = OPEN
	cb.LastAttemptedAt = time.Now()

	err := cb.Attempt()
	if err == nil || err.Error() != "circuit breaker is open" {
		t.Errorf("Expected error 'circuit breaker is open', got %v", err)
	}
}

func TestFail(t *testing.T) {
	domain.TestIsUnit(t)
	logger := zerolog.New(nil)
	cfg := &Config{
		RecordLength: 5,
		Percentile:   0.6,
	}
	cb := newCircuitBreaker(cfg, &logger)

	for i := 0; i < 3; i++ {
		cb.Fail(errors.New("some error"))
	}

	if cb.State != OPEN {
		t.Errorf("Expected State to be OPEN, got %v", cb.State)
	}
}

func TestSuccess(t *testing.T) {
	domain.TestIsUnit(t)
	logger := zerolog.New(nil)
	cfg := &Config{
		RecordLength:     5,
		Percentile:       0.6,
		RecoveryRequests: 2,
	}
	cb := newCircuitBreaker(cfg, &logger)
	cb.State = HALFOPEN

	cb.Success()
	cb.Success()

	if cb.State != CLOSED {
		t.Errorf("Expected State to be CLOSED, got %v", cb.State)
	}
}
