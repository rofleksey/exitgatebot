package twitch_irc

import (
	"context"
	"exitgatebot/app/client/twitch_api"
	"exitgatebot/app/config"
	"regexp"
	"strings"
	"sync"
	"time"

	"log/slog"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/samber/do"
)

const firstMessageReplyText = "Hello fucking burger? Sosal?"

var englishRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\p{P}]+$`)

type MessageHandler func(channel, username, messageID, text string, tags map[string]string)

var _ do.Shutdownable = (*Client)(nil)

type Client struct {
	appCtx    context.Context
	cfg       *config.Config
	apiClient *twitch_api.Client
	ircClient *twitch.Client

	mutex sync.RWMutex
}

func NewClient(di *do.Injector) (*Client, error) {
	cfg := do.MustInvoke[*config.Config](di)
	apiClient := do.MustInvoke[*twitch_api.Client](di)

	client := &Client{
		appCtx:    do.MustInvoke[context.Context](di),
		cfg:       cfg,
		apiClient: apiClient,
	}

	client.ircClient = twitch.NewClient(cfg.Twitch.ClientID, "oauth:"+apiClient.AccessToken())
	client.setupIRCListeners()

	return client, nil
}

func (c *Client) setupIRCListeners() {
	c.ircClient.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		if !msg.FirstMessage {
			return
		}

		username := strings.ToLower(msg.User.Name)
		channel := strings.ToLower(strings.TrimPrefix(msg.Channel, "#"))
		text := strings.TrimSpace(msg.Message)

		c.handleFirstMessage(channel, username, msg.ID, text)
	})

	c.ircClient.OnConnect(func() {
		for _, profile := range c.cfg.Profiles {
			c.ircClient.Join(strings.ToLower(profile.Username))
		}
		slog.Info("Connected to Twitch IRC")
	})

	c.ircClient.OnReconnectMessage(func(message twitch.ReconnectMessage) {
		for _, profile := range c.cfg.Profiles {
			c.ircClient.Join(strings.ToLower(profile.Username))
		}
		slog.Info("Reconnecting to Twitch IRC")
	})
}

func (c *Client) handleFirstMessage(channel, username, msgID, text string) {
	if !englishRegex.MatchString(text) {
		return
	}

	err := c.apiClient.ReplyMessage(channel, msgID, firstMessageReplyText)
	if err != nil {
		slog.Error("Failed to reply to first message",
			slog.String("username", username),
			slog.Any("error", err),
		)
		return
	}

	slog.Info("Replied to first message",
		slog.String("username", username),
		slog.Bool("telegram", true),
	)
}

func (c *Client) RunRefreshLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.refreshToken()
		}
	}
}

func (c *Client) RunConnectLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := c.ircClient.Connect(); err != nil {
				slog.Error("irc connect error", "error", err)
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func (c *Client) Shutdown() error {
	c.ircClient.Disconnect()
	return nil
}

func (c *Client) refreshToken() {
	c.ircClient.SetIRCToken("oauth:" + c.apiClient.AccessToken())
}
