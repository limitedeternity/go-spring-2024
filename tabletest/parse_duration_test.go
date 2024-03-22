package tabletest

import (
	"math/rand"
	"testing"
	"time"
)

var parseDurationTests = []struct {
	input    string
	duration time.Duration
	ok       bool
}{
	// units
	{"0", 0, true},
	{"5s", 5 * time.Second, true},
	{"5m", 5 * time.Minute, true},
	{"5h", 5 * time.Hour, true},
	{"1488ms", 1488 * time.Millisecond, true},
	{"1337ns", 1337 * time.Nanosecond, true},
	{"228us", 228 * time.Microsecond, true},
	{"228µs", 228 * time.Microsecond, true}, // U+00B5
	{"228μs", 228 * time.Microsecond, true}, // U+03BC

	// signs
	{"-5s", -5 * time.Second, true},
	{"+5s", 5 * time.Second, true},
	{"-0", 0, true},
	{"+0", 0, true},

	// floating
	{"5.0s", 5 * time.Second, true},
	{"5.3s", 5*time.Second + 300*time.Millisecond, true},
	{"5.s", 5 * time.Second, true},
	{".5s", 500 * time.Millisecond, true},
	{"5.00s", 5 * time.Second, true},
	{"5.000s", 5 * time.Second, true},
	{"5.004s", 5*time.Second + 4*time.Millisecond, true},
	{"5.0040s", 5*time.Second + 4*time.Millisecond, true},

	// composite
	{"5h30m", 5*time.Hour + 30*time.Minute, true},
	{"30m3h", 30*time.Minute + 3*time.Hour, true},
	{"5m10.5s", 5*time.Minute + 10*time.Second + 500*time.Millisecond, true},
	{"10.5s5m", 5*time.Minute + 10*time.Second + 500*time.Millisecond, true},
	{"-5m10.5s", -(5*time.Minute + 10*time.Second + 500*time.Millisecond), true},
	{"-10.5s5m", -(5*time.Minute + 10*time.Second + 500*time.Millisecond), true},
	{"1h2m3s4ms5ns6µs", 1*time.Hour + 2*time.Minute + 3*time.Second + 4*time.Millisecond + 6*time.Microsecond + 5*time.Nanosecond, true},
	{"1h2m14.625s", 1*time.Hour + 2*time.Minute + 14*time.Second + 625*time.Millisecond, true},

	// large
	{"92233720368ns", 92233720368 * time.Nanosecond, true},
	{"9223372036854775807ns", (1<<63 - 1) * time.Nanosecond, true},

	// https://golang.org/issue/6617
	{"0.3333333333333333333h", 20 * time.Minute, true},

	// https://golang.org/issue/15011
	{"0.100000000000000000000h", 6 * time.Minute, true},

	// overflow check
	{"0.9223372036854775807h", 55*time.Minute + 20*time.Second + 413933267*time.Nanosecond, true},

	// errors
	{"", 0, false},
	{"5", 0, false},
	{"-", 0, false},
	{"s", 0, false},
	{".", 0, false},
	{"-.", 0, false},
	{"+.", 0, false},
	{".s", 0, false},
	{"+.s", 0, false},

	// overflows
	{"3000000h", 0, false},
	{"9223372036854775808ns", 0, false},
	{"-9223372036854775808ns", 0, false},
	{"9223372036854775.808us", 0, false},
	{"9223372036854ms775μs808ns", 0, false},
}

func TestParseDuration(t *testing.T) {
	for _, tc := range parseDurationTests {
		res, err := ParseDuration(tc.input)
		if tc.ok && (err != nil || res != tc.duration) {
			t.Errorf("ParseDuration(%q) = (%v, %v), want (%v, nil)", tc.input, res, err, tc.duration)
		} else if !tc.ok && err == nil {
			t.Errorf("ParseDuration(%q) = (_, nil), want (_, non-nil)", tc.input)
		}
	}
}

func TestParseDurationRoundTrip(t *testing.T) {
	for i := 0; i < 100; i++ {
		// Resolutions finer than milliseconds will result in imprecise round-trips.
		res0 := time.Duration(rand.Int31()) * time.Millisecond
		s := res0.String()
		res1, err := ParseDuration(s)

		if err != nil || res0 != res1 {
			t.Errorf("round-trip failed: %d => %q => %d, %v", res0, s, res1, err)
		}
	}
}
