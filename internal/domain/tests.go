package domain

import (
	"os"
	"testing"
)

const all = "all"

func TestIsIntegration(t *testing.T) {
	if os.Getenv("TEST_TYPE") != "integration" && os.Getenv("TEST_TYPE") != all {
		t.Skip("skipping integration test")
	}
}

func TestIsUnit(t *testing.T) {
	if os.Getenv("TEST_TYPE") != "unit" && os.Getenv("TEST_TYPE") != all {
		t.Skip("skipping unit test")
	}
}

func TestIsE2E(t *testing.T) {
	if os.Getenv("TEST_TYPE") != "e2e" && os.Getenv("TEST_TYPE") != all {
		t.Skip("skipping e2e test")
	}
}
