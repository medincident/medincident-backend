package ignore_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/medincident/medincident-command-service/internal/util/ignore"
)

func TestAny_DiscardsSecondReturnRegardlessOfType(t *testing.T) {
	assert.Equal(t, 7, ignore.Any(7, errors.New("err")))
	assert.Equal(t, "ok", ignore.Any("ok", true))
	assert.Equal(t, 3.14, ignore.Any(3.14, struct{}{}))
}

func TestError_ReturnsValueOnSuccess(t *testing.T) {
	id := uuid.MustParse("00000000-0000-7000-8000-000000000001")
	got := ignore.Error(uuid.Parse("00000000-0000-7000-8000-000000000001"))
	assert.Equal(t, id, got)
}

func TestError_ReturnsZeroOnFailure(t *testing.T) {
	got := ignore.Error(uuid.Parse("garbage"))
	assert.Equal(t, uuid.Nil, got)
}

func TestError_PreservesValueEvenWithError(t *testing.T) {
	fn := func() (int, error) { return 42, errors.New("ignored") }
	assert.Equal(t, 42, ignore.Error(fn()))
}

func TestExists_ReturnsValueIgnoringFlag(t *testing.T) {
	lookup := func(present bool) (int, bool) {
		if present {
			return 42, true
		}
		return 0, false
	}
	assert.Equal(t, 42, ignore.Exists(lookup(true)))
	assert.Equal(t, 0, ignore.Exists(lookup(false)))
}
