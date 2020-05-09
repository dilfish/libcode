// default test files

package libcode

import (
	"testing"
)

func TestLibcodeEncodeDefault(t *testing.T) {
	r := EncodeDefault(1334)
	if r != 1334 {
		t.Error("expect 1334, got", r)
	}
}

func TestLibcodeDecodeDefault(t *testing.T) {
	d := DecodeDefault(1335)
	if d != 1335 {
		t.Error("expect 1335, got", d)
	}
}
