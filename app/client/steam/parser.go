package steam

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	timestampPatterns = []string{
		"2 Jan @ 3:04pm",
		"2 Jan @ 3:04PM",
		"02 Jan @ 3:04pm",
		"02 Jan @ 3:04PM",
	}

	currentYear = time.Now().Year()
)

func ParseComments(doc *goquery.Document) ([]Comment, error) {
	commentThread := doc.Find(commentThreadSelector)
	if commentThread.Length() == 0 {
		return nil, fmt.Errorf("no comment thread found on the page")
	}

	comments := extractComments(commentThread)

	return comments, nil
}

func extractComments(commentThread *goquery.Selection) []Comment {
	var comments []Comment

	commentThread.Find(commentSelector).Each(func(i int, commentElem *goquery.Selection) {
		comment := parseSingleComment(commentElem)
		if comment.Author != "" && comment.Content != "" {
			comments = append(comments, comment)
		}
	})

	return comments
}

func parseSingleComment(commentElem *goquery.Selection) Comment {
	comment := Comment{
		CommentID: getCommentID(commentElem),
	}

	authorElem := commentElem.Find(authorSelector).First()
	if authorElem.Length() > 0 {
		comment.Author = strings.TrimSpace(authorElem.Text())
		comment.AuthorURL = getAuthorURL(authorElem)
	}

	avatarElem := commentElem.Find(avatarSelector).First()
	if avatarElem.Length() > 0 {
		comment.AvatarURL = getAvatarURL(avatarElem)
	}

	contentElem := commentElem.Find(contentSelector).First()
	if contentElem.Length() > 0 {
		comment.Content = strings.TrimSpace(contentElem.Text())
	}

	timestampElem := commentElem.Find(timestampSelector).First()
	if timestampElem.Length() > 0 {
		timestampText := strings.TrimSpace(timestampElem.Text())
		comment.TimestampRaw = timestampText
		comment.Timestamp = parseTimestamp(timestampText)
	}

	return comment
}

func parseTimestamp(timestampStr string) time.Time {
	timestampStr = strings.TrimSpace(timestampStr)

	for _, layout := range timestampPatterns {
		fullTimestamp := fmt.Sprintf("%d %s", currentYear, timestampStr)
		fullLayout := fmt.Sprintf("2006 %s", layout)

		if t, err := time.Parse(fullLayout, fullTimestamp); err == nil {
			return t
		}

		if t, err := time.Parse(layout, timestampStr); err == nil {
			t = time.Date(currentYear, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
			return t
		}
	}

	return time.Time{}
}
