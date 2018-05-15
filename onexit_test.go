package onexit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)

	var order []int

	// the functions are being registered from lowest priority
	// to highest.
	for i := -10; i <= 10; i++ {
		i := i
		Register(func() {
			order = append(order, i)
		}, i)
	}

	go forceExit(0, false)

	<-Done()

	// the highest priority will run first.
	for i := 0; i < 21; i++ {
		assert.Equal(20-i-10, order[i])
	}
}
