package steam

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaznevich(t *testing.T) {
	client, err := NewClient(do.New())
	require.NoError(t, err)

	res, err := client.ParseCommentsFromURL(t.Context(), "https://steamcommunity.com/id/Maznevich")
	require.NoError(t, err)

	assert.NotEmpty(t, res.Comments)
}
