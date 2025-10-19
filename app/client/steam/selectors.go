package steam

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	commentThreadSelector = ".commentthread_area"
	commentSelector       = "div.commentthread_comment"
	authorSelector        = "a.commentthread_author_link"
	avatarSelector        = "div.commentthread_comment_avatar img"
	contentSelector       = "div.commentthread_comment_text"
	timestampSelector     = "span.commentthread_comment_timestamp"
)

func getCommentID(selection *goquery.Selection) string {
	if id, exists := selection.Attr("id"); exists {
		return id
	}
	return ""
}

func getAuthorURL(selection *goquery.Selection) string {
	if href, exists := selection.Attr("href"); exists {
		return href
	}
	return ""
}

func getAvatarURL(selection *goquery.Selection) string {
	if src, exists := selection.Attr("src"); exists {
		return src
	}

	if srcset, exists := selection.Attr("srcset"); exists {
		urls := strings.Split(srcset, " ")
		if len(urls) > 0 {
			return urls[0]
		}
	}

	return ""
}
