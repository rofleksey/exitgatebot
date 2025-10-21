package openai

import (
	"context"
	"exitgatebot/app/config"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaznevich(t *testing.T) {
	cfg, err := config.Load("../../../config.yaml")
	require.NoError(t, err)

	di := do.New()
	do.ProvideValue(di, cfg)

	client, err := NewClient(di)
	require.NoError(t, err)

	summary, err := client.SummarizeComment(context.Background(), "-rep you should lose some weight, idiot. Maybe go to gym instead of being a toxic trash in dbd")
	require.NoError(t, err)

	assert.NotEmpty(t, summary)
}
