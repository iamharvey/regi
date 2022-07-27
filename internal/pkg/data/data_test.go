package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDB_Current(t *testing.T) {
	db, err := NewDB()
	assert.NoError(t, err)
	assert.NotEmpty(t, db)
}
