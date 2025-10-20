package checker

import (
	"context"
	"exitgatebot/app/client/openai"
	"exitgatebot/app/client/steam"
	"exitgatebot/app/client/twitch"
	"exitgatebot/app/config"
	"fmt"
	"log/slog"
	"time"

	"github.com/elliotchance/pie/v2"
	"github.com/rofleksey/meg"
	"github.com/samber/do"
)

const notificationFormat = "У тебя появился новый комментарий в профиле PogChamp %s"

type Service struct {
	cfg          *config.Config
	twitchClient *twitch.Client
	steamClient  *steam.Client
	openaiClient *openai.Client
	db           *CommentDatabase
}

func New(di *do.Injector) (*Service, error) {
	return &Service{
		cfg:          do.MustInvoke[*config.Config](di),
		twitchClient: do.MustInvoke[*twitch.Client](di),
		steamClient:  do.MustInvoke[*steam.Client](di),
		openaiClient: do.MustInvoke[*openai.Client](di),
		db:           newCommentDatabase(),
	}, nil
}

func (s *Service) runCheck(ctx context.Context) {
	slog.Info("Starting check")
	defer slog.Info("Check finished")

	for _, profile := range s.cfg.Profiles {
		if err := s.runCheckForProfile(ctx, profile); err != nil {
			slog.ErrorContext(ctx, "Check for profile failed",
				slog.String("username", profile.Username),
				slog.Any("error", err),
			)
		}

		time.Sleep(5 * time.Second)
	}
}

func (s *Service) runCheckForProfile(ctx context.Context, profile config.Profile) error {
	slogger := slog.With(
		slog.String("username", profile.Username),
	)

	slogger.Debug("Starting check for profile")
	defer slogger.Debug("Check for profile finished")

	existingCommentCount := s.db.Count(profile.Username)
	needReaction := existingCommentCount > 0

	steamRes, err := s.steamClient.ParseCommentsFromURL(ctx, profile.SteamLink)
	if err != nil {
		return fmt.Errorf("ParseCommentsFromURL: %w", err)
	}
	slogger.Debug("Got comments from steam",
		slog.String("count", profile.Username),
	)

	remoteComments := steamRes.Comments
	newComments := pie.Filter(remoteComments, func(c steam.Comment) bool {
		return !s.db.Has(profile.Username, c.CommentID)
	})
	if len(newComments) == 0 {
		return nil
	}

	slogger.Debug("Got new comments",
		slog.Int("count", len(newComments)),
	)

	for _, comment := range pie.Reverse(newComments) {
		if needReaction {
			if err := s.notifyNewComment(ctx, slogger, profile, comment); err != nil {
				return fmt.Errorf("notifyNewComment: %w", err)
			}
		}

		s.db.Store(profile.Username, comment.CommentID)

		slogger.Debug("Added comment to database",
			slog.String("comment_id", comment.CommentID),
		)
	}

	return nil
}

func (s *Service) notifyNewComment(ctx context.Context, slogger *slog.Logger, profile config.Profile, comment steam.Comment) error {
	slogger = slog.With(
		slog.String("comment_id", comment.CommentID),
	)

	slogger.Debug("Reacting to comment")

	summary, err := s.openaiClient.SummarizeComment(ctx, comment.Content)
	if err != nil {
		return fmt.Errorf("SummarizeComment: %w", err)
	}

	notificationText := fmt.Sprintf(notificationFormat, summary)

	if !s.cfg.Twitch.DisableNotifications {
		if err = s.twitchClient.SendMessage(profile.Username, notificationText); err != nil {
			return fmt.Errorf("SendMessage: %w", err)
		}

		slogger.Info("Successfully reacted to comment",
			slog.String("message", notificationText),
			slog.Bool("telegram", true),
		)
	} else {
		slogger.Info("Would react to comment, but notifications are disabled",
			slog.String("message", notificationText),
			slog.Bool("telegram", true),
		)
	}

	return nil
}

func (s *Service) RunCheckLoop(ctx context.Context) {
	interval := time.Duration(s.cfg.Steam.CheckInterval) * time.Minute

	s.runCheck(ctx)
	meg.RunTicker(ctx, interval, func() {
		s.runCheck(ctx)
	})
}
