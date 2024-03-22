//go:build !solution

package testequal

import (
	"bytes"
	"fmt"
	"reflect"
)

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.

func equal(expected, actual any) bool {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return false
	}

	if expected == nil || actual == nil {
		return expected == nil && actual == nil
	}

	switch e := expected.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return expected == actual

	case struct{}:
		return false

	case map[any]any:
		mapA, ok := actual.(map[any]any)

		if !ok {
			return false
		}

		if len(e) != len(mapA) {
			return false
		}

		for key, eVal := range e {
			var aVal any
			aVal, ok = mapA[key]

			if !ok || !equal(eVal, aVal) {
				return false
			}
		}

		return true
	}

	bytesE, okE := expected.([]byte)
	bytesA, okA := actual.([]byte)

	if !okE || !okA {
		return reflect.DeepEqual(expected, actual)
	}

	if bytesE == nil || bytesA == nil {
		return bytesE == nil && bytesA == nil
	}

	return bytes.Equal(bytesE, bytesA)
}

func errorf(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	format :=
		`
		expected: %v
        actual  : %v
        message : %v`

	msg := ""

	switch len(msgAndArgs) {
	case 0:
		break
	case 1:
		msg = msgAndArgs[0].(string)

	default:
		msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}

	t.Errorf(format, expected, actual, msg)
}

func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if equal(expected, actual) {
		return true
	}

	errorf(t, expected, actual, msgAndArgs...)
	return false
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !equal(expected, actual) {
		return true
	}

	errorf(t, expected, actual, msgAndArgs...)
	return false
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if equal(expected, actual) {
		return
	}

	errorf(t, expected, actual, msgAndArgs...)
	t.FailNow()
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !equal(expected, actual) {
		return
	}

	errorf(t, expected, actual, msgAndArgs...)
	t.FailNow()
}
