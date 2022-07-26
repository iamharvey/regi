package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDB_Current(t *testing.T) {
	db, err := NewDB()
	assert.NoError(t, err)
	assert.NotEmpty(t, db)

	reg, err := db.Current()
	assert.NoError(t, err)
	assert.NotEmpty(t, reg)
	assert.Equal(t, "192.168.0.168:5000", reg.Name)
}
