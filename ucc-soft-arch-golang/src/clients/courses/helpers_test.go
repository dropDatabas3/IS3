package courses

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_parseUUID(t *testing.T) {
	id := uuid.New()
	// valid string
	got := parseUUID(id.String())
	require.Equal(t, id, got)
	// nil returns zero UUID
	require.Equal(t, uuid.Nil, parseUUID(nil))
}

func Test_toBool(t *testing.T) {
	require.True(t, toBool(true))
	require.True(t, toBool(int(1)))
	require.True(t, toBool(int64(1)))
	require.True(t, toBool(float64(1)))
	require.True(t, toBool([]byte("1")))
	require.True(t, toBool("1"))
	require.True(t, toBool([]byte("true")))
	require.True(t, toBool("TrUe"))

	require.False(t, toBool(false))
	require.False(t, toBool(int(0)))
	require.False(t, toBool(float64(0)))
	require.False(t, toBool([]byte("0")))
	require.False(t, toBool("0"))
	require.False(t, toBool("false"))
	require.False(t, toBool(struct{}{}))
}

func Test_toInt(t *testing.T) {
	require.Equal(t, 5, toInt(5))
	require.Equal(t, 6, toInt(int32(6)))
	require.Equal(t, 7, toInt(int64(7)))
	require.Equal(t, 3, toInt(float32(3.9)))
	require.Equal(t, 4, toInt(float64(4.1)))
	require.Equal(t, 12, toInt([]byte("12")))
	require.Equal(t, 0, toInt(struct{}{}))
}

func Test_toFloat64(t *testing.T) {
	require.InDelta(t, 1.5, toFloat64(float64(1.5)), 0.0001)
	require.InDelta(t, 2.5, toFloat64(float32(2.5)), 0.0001)
	require.InDelta(t, 3.0, toFloat64(int(3)), 0.0001)
	require.InDelta(t, 4.0, toFloat64(int64(4)), 0.0001)
	require.InDelta(t, 5.75, toFloat64([]byte("5.75")), 0.0001)
	require.InDelta(t, 0.0, toFloat64(struct{}{}), 0.0001)
}

func Test_GetAll_Empty_NoError(t *testing.T) {
	db := setupCoursesDB(t)
	c := NewCourseClient(db)
	got, err := c.GetAll("")
	require.NoError(t, err)
	require.Len(t, got, 0)
}
