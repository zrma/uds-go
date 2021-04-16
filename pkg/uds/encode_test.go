package uds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	const (
		str    = `abc123!?$*&()'-=@~`
		base64 = `YWJjMTIzIT8kKiYoKSctPUB+`
	)

	t.Run("should encode", func(t *testing.T) {
		assert.Equal(t, base64, encode([]byte(str)))
	})

	t.Run("should decode", func(t *testing.T) {
		actual, err := decode(base64)
		assert.NoError(t, err)
		assert.Equal(t, []byte(str), actual)
	})
}
