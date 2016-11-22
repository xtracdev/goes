package uuid_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xtracdev/goes/uuid"
)

func TestGenerateUuidV4(t *testing.T) {
	uuid, err := uuid.GenerateUuidV4()
	assert.Nil(t, err)
	assert.NotEmpty(t, uuid)
	t.Logf("%s", uuid)
	//check that the uuid complies rfc4122 v4
	pattern := "[a-f0-9]{8}-[a-f0-9]{4}-[4][a-f0-9]{3}-(8|9|a|b)[a-f0-9]{3}-[a-f0-9]{12}"
	m, _ := regexp.MatchString(pattern, uuid)
	assert.True(t, m)
}
