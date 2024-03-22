package iprange_test

import (
	"testing"

	"gitlab.com/slon/shad-go/iprange"
)

func FuzzParser(f *testing.F) {
	testcases := []string{"10.0.0.1, 10.0.0.5-10, 192.168.1.*, 192.168.10.0/24"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Parse(%q) panicked: %v", orig, r)
			}
		}()

		_, err := iprange.Parse(orig)

		if err != nil {
			return
		}
	})
}
