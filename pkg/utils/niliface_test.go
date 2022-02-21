package utils_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"code.tokarch.uk/mainnika/nikt-link-proxy/pkg/utils"
)

func TestIsNilInterface(t *testing.T) {
	t.Run("expect detect nil", func(t *testing.T) {

		var nilIP *net.IP = nil
		var falseTester = func(i interface{}) bool {
			return i == nil
		}
		var trueTester = func(i interface{}) bool {
			return utils.IsNilInterface(i)
		}

		assert.False(t, falseTester(nilIP))
		assert.True(t, trueTester(nilIP))
	})
	t.Run("expect works fine", func(t *testing.T) {

		var nilIP *net.IP = &net.IP{}
		var falseTester = func(i interface{}) bool {
			return i == nil
		}
		var trueTester = func(i interface{}) bool {
			return utils.IsNilInterface(i)
		}

		assert.False(t, falseTester(nilIP))
		assert.False(t, trueTester(nilIP))
	})
}
