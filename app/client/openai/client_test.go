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

	summary, err := client.SummarizeComment(context.Background(), "You are gonna die alone in your room still being a virgin, cause you are ugly and toxic af.")
	require.NoError(t, err)

	assert.NotEmpty(t, summary)
}
