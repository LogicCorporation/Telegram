package utils

import (
	"fmt"
	"github-telegram-notify/types"
	"html"
	"strings"
	"time"
)

func CreateContents(meta *types.Metadata) (text string, err error) {
	event, _ := meta.ParseEvent()
	switch meta.EventName {
	case "fork":
		event := event.(*types.ForkEvent)

		// No Activity Types

		text = createForkText(event)
	case "issue_comment":
		event := event.(*types.IssueCommentEvent)

		if !Contains([]string{"created", "deleted"}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createIssueCommentText(event)
	case "issues":
		event := event.(*types.IssuesEvent)

		if !Contains([]string{
			"created", "closed", "opened", "reopened", "locked", "unlocked", // More to be added.
		}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createIssuesText(event)
	case "pull_request":
		event := event.(*types.PullRequestEvent)

		if !Contains([]string{
			"created", "opened", "reopened", "locked", "unlocked", "closed", "synchronize", // More to be added.
		}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createPullRequestText(event)
	case "pull_request_review_comment":
		event := event.(*types.PullRequestReviewCommentEvent)

		if !Contains([]string{"created", "deleted"}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createPullRequestReviewCommentText(event)
	case "push":
		event := event.(*types.PushEvent)
		// No Activity Types
		text = createPushText(event)
	case "release":
		event := event.(*types.ReleaseEvent)
		if !Contains([]string{"published", "released"}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createReleaseText(event)
	case "watch":
		event := event.(*types.WatchEvent)

		if !Contains([]string{"started"}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createStarText(event)
	}
	return text, nil
}

func createPushText(event *types.PushEvent) string {
    // 1. Extract and format the repository name
    repoFullName := event.Repo.FullName // Example: "LogicCorporation/ProxyChecker-v2"
    repoParts := strings.Split(repoFullName, "/")
    repoName := repoParts[len(repoParts)-1]               // "ProxyChecker-v2"
    repoName = strings.ReplaceAll(repoName, "-", " ")      // "ProxyChecker v2"

    // 2. Get the publication date in [DD/MM/AA] format
    pubDate := time.Now().Format("02/01/06")

    // 3. Create the header with emoji, bold, and underline
    header := fmt.Sprintf("🚀 <b>%d New Update(s) to <u>%s</u></b>\n\n",
        len(event.Commits),
        repoName,
    )

    // 4. Create the updates section with bold and underline for "📌 Updates:"
    updatesText := "<blockquote><b><u>📌 Updates:</u></b>\n"
    for _, commit := range event.Commits {
        commitMessage := html.EscapeString(commit.Message)
        updatesText += fmt.Sprintf("• %s\n", commitMessage)
    }
    updatesText += "</blockquote>"

    // 5. Create the footer with a thank you message and publication date
    footer := fmt.Sprintf("\nSpecial thanks to accompany, stay tuned for more. [ %s ]", pubDate)

    // 6. Combine all parts into the final message
    text := header + updatesText + footer

    return text
}



func createForkText(event *types.ForkEvent) string {
	return fmt.Sprintf("🍴 <a href='%s'>%s</a> forked <a href='%s'>%s</a> → <a href='%s'>%s</a>",
		event.Sender.HTMLURL,
		event.Sender.Login,
		event.Repo.HTMLURL,
		event.Repo.FullName,
		event.Forkee.HTMLURL,
		event.Forkee.FullName,
	)
}

func createIssueCommentText(event *types.IssueCommentEvent) string {
	return fmt.Sprintf("🗣 <a href='%s'>%s</a> commented on issue <a href='%s'>%s</a> in <a href='%s'>%s</a>",
		event.Sender.HTMLURL,
		event.Sender.Login,
		event.Issue.HTMLURL,
		html.EscapeString(event.Issue.Title),
		event.Repo.HTMLURL,
		event.Repo.FullName,
	)
}

func createIssuesText(event *types.IssuesEvent) string {
	return fmt.Sprintf("🐛 <a href='%s'>%s</a> %s issue <a href='%s'>%s</a> in <a href='%s'>%s</a>",
		event.Sender.HTMLURL,
		event.Sender.Login,
		event.Action,
		event.Issue.HTMLURL,
		html.EscapeString(event.Issue.Title),
		event.Repo.HTMLURL,
		event.Repo.FullName,
	)
}

func createPullRequestText(event *types.PullRequestEvent) (text string) {
	text = fmt.Sprintf("🔌 <a href='%s'>%s</a> ", event.Sender.HTMLURL, event.Sender.Login)
	text += event.Action
	if event.Action == "opened" {
		text += " a new"
	}
	text += " pull request "
	text += fmt.Sprintf("<a href='%s'>%s</a>", event.PullRequest.HTMLURL, html.EscapeString(event.PullRequest.Title))
	text += fmt.Sprintf(" in <a href='%s'>%s</a>", event.Repo.HTMLURL, event.Repo.FullName)
	return text
}

func createPullRequestReviewCommentText(event *types.PullRequestReviewCommentEvent) string {
	return fmt.Sprintf("🧐 <a href='%s'>%s</a> commented on PR review <a href='%s'>%s</a> in <a href='%s'>%s</a>",
		event.Sender.HTMLURL,
		event.Sender.Login,
		event.PullRequest.HTMLURL,
		html.EscapeString(event.PullRequest.Title),
		event.Repo.HTMLURL,
		event.Repo.FullName,
	)
}

func createReleaseText(event *types.ReleaseEvent) (text string) {
	text = "🎊 A new "
	if event.Release.Prerelease {
		text += "pre"
	}
	text += fmt.Sprintf("release was %s in <a href='%s'>%s</a> by <a href='%s'>%s</a>\n",
		event.Action,
		event.Repo.HTMLURL,
		event.Repo.FullName,
		event.Sender.HTMLURL,
		event.Sender.Login,
	)
	text += fmt.Sprintf("\n📍 <a href='%s'>%s</a> (<code>%s</code>)\n\n", event.Release.HTMLURL, event.Release.Name, event.Release.TagName)
	if event.Release.Assets != nil {
		text += "📦 <b>Assets:</b>\n"
		for _, asset := range event.Release.Assets {
			text += fmt.Sprintf("• <a href='%s'>%s</a>\n", asset.BrowserDownloadURL, html.EscapeString(asset.Name))
		}
	}

	return
}

func createStarText(event *types.WatchEvent) string {
	return fmt.Sprintf("🌟 <a href='%s'>%s</a> starred <a href='%s'>%s</a>",
		event.Sender.HTMLURL,
		event.Sender.Login,
		event.Repo.HTMLURL,
		event.Repo.FullName,
	)
}

// Helper function to extract the first name from a full name
func getFirstName(fullName string) string {
    // Split the full name by spaces
    parts := strings.Fields(fullName)
    if len(parts) > 0 {
        return parts[0]
    }
    return fullName // Return as is if splitting fails
}

