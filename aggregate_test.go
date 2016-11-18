package goes

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAggregate(t *testing.T) {
	agg, err := NewAggregate()
	assert.Nil(t, err)
	assert.NotNil(t, agg)
	assert.IsType(t, new(Aggregate), agg)
}

func TestGenerateID(t *testing.T) {
	pattern := "[a-f0-9]{8}-[a-f0-9]{4}-[4][a-f0-9]{3}-(8|9|a|b)[a-f0-9]{3}-[a-f0-9]{12}"

	uuid, err := GenerateID()
	assert.Nil(t, err)
	assert.NotEmpty(t, uuid)

	t.Logf("%s", uuid)
	m, _ := regexp.MatchString(pattern, uuid)
	assert.True(t, m)

}
