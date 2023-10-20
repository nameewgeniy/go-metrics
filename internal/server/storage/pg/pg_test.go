package pg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Successfully pings the database
func TestPingDatabaseSuccessfully(t *testing.T) {

	// Create a new Pg instance with the mock config
	pg := Pg{}

	// Call the Ping method
	err := pg.Ping()

	// Assert that the Ping method returns nil error
	assert.Nil(t, err)
}
