package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewErrorAndErrorString(t *testing.T) {
	e := NewError("NOT_FOUND", "missing", http.StatusNotFound)
	require.Equal(t, "NOT_FOUND", e.Code)
	require.Equal(t, "missing", e.Message)
	require.Equal(t, http.StatusNotFound, e.HTTPStatusCode)
	require.Equal(t, "NOT_FOUND: missing", e.Error())
}
