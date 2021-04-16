package uds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
		tb = gb * 1024
	)

	for _, tc := range []struct {
		description string
		given       int64
		want        string
		ok          bool
	}{
		{"invalid - negative size", -1, "", false},
		{"", 1000, "1000.0 bytes", true},
		{"", kb, "1.0 KB", true},
		{"", 800 * kb, "800.0 KB", true},
		{"", mb, "1.0 MB", true},
		{"", tb, "1.0 TB", true},
		{"", 1024 * tb, "1024.0 TB", true},
	} {
		t.Run(tc.description, func(t *testing.T) {
			size, err := format(tc.given)
			if tc.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.want, size)
		})
	}
}
