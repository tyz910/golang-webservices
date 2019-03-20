package main

import (
	"bytes"
	"testing"
)

var testOkInput = `1
2
3
3
4
5`

var testOkResult = `1
2
3
4
5
`

var testFailInput = `1
2
1`

func TestOK(t *testing.T) {
	in := bytes.NewBufferString(testOkInput)
	out := bytes.NewBuffer(nil)
	err := uniq(in, out)
	if err != nil {
		t.Errorf("Test OK failed: %s", err)
	}
	result := out.String()
	if result != testOkResult {
		t.Errorf("Test OK failed, result not match")
	}
}

func TestFail(t *testing.T) {
	in := bytes.NewBufferString(testFailInput)
	out := bytes.NewBuffer(nil)
	err := uniq(in, out)
	if err == nil {
		t.Errorf("Test FAIL failed: expected error")
	}
}
