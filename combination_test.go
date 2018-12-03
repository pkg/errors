package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Combine(t *testing.T) {
	var (
		err error
		e0  = New("e0")
		e1  = New("e1")
		e2  = New("e2")
	)

	err = Combine()
	assert.NoError(t, err)

	err = Combine(nil)
	assert.NoError(t, err)

	err = Combine(nil, nil)
	assert.NoError(t, err)

	err = Combine(e0)
	assert.EqualError(t, err, "e0")
	assert.Equal(t, e0, Cause(err))

	err = Combine(nil, e0)
	assert.EqualError(t, err, "e0")
	assert.Equal(t, e0, Cause(err))

	err = Combine(e0, nil)
	assert.EqualError(t, err, "e0")
	assert.Equal(t, e0, Cause(err))

	err = Combine(e0, e1)
	assert.EqualError(t, err, "e0\ne1")
	assert.Equal(t, e0, Cause(err))
	assert.NotEqual(t, e1, Cause(err))

	err = Combine(e0, e1, nil)
	assert.EqualError(t, err, "e0\ne1")
	assert.Equal(t, e0, Cause(err))
	assert.NotEqual(t, e1, Cause(err))

	err = Combine(nil, e0, e1, nil)
	assert.EqualError(t, err, "e0\ne1")
	assert.Equal(t, e0, Cause(err))
	assert.NotEqual(t, e1, Cause(err))

	err = Combine(nil, e0, nil, e1, nil)
	assert.EqualError(t, err, "e0\ne1")
	assert.Equal(t, e0, Cause(err))
	assert.NotEqual(t, e1, Cause(err))

	err = Combine(e0, e1, e2)
	assert.EqualError(t, err, "e0\ne1\ne2")

	err = Combine(e0, nil, e2)
	assert.EqualError(t, err, "e0\ne2")
}
