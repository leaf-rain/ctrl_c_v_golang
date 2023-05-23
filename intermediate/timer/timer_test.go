package timer

import (
	"testing"
	"time"
)

func TestToYYYYWW(t *testing.T) {
	t.Log(ToYYYYWW(time.Now()))
	t.Log(time.Now().ISOWeek())
}
