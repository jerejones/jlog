package target

import (
	"strings"
	"testing"
)

func assertContains(t *testing.T, str, substr string) {
	t.Helper()
	if !strings.Contains(str, substr) {
		t.Error(`Error: "%s" should have contained "%s"`, str, substr)
	}
}

func assertNotContains(t *testing.T, str, substr string) {
	t.Helper()
	if strings.Contains(str, substr) {
		t.Error(`Error: "%s" should not have contained "%s"`, str, substr)
	}
}
