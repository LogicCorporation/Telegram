package utils

import (
	"fmt"
	"github-telegram-notify/types"
	"html"
	"strings"
)

func CreateContents(meta *types.Metadata) (text string, markupText string, markupUrl string, err error) {
	event, _ := meta.ParseEvent()
	switch meta.EventName {
	case "fork":
		event := event.(*types.ForkEvent)

		// No Activity Types

		text = createForkText(event)
		markupText = fmt.Sprintf("Total Forks: %d", event.Repo.ForksCount)
		markupUrl = event.Repo.HTMLURL + "/network/members"
	case "issue_comment":
		event := event.(*types.IssueCommentEvent)

		if !Contains([]string{"created", "deleted"}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createIssueCommentText(event)
		markupText = "Open Comment"
		markupUrl = event.Comment.HTMLURL
	case "issues":
		event := event.(*types.IssuesEvent)

		if !Contains([]string{
			"created", "closed", "opened", "reopened", "locked", "unlocked", // More to be added.
		}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createIssuesText(event)
		markupText = "Open Issue"
		markupUrl = event.Issue.HTMLURL
	case "pull_request":
		event := event.(*types.PullRequestEvent)

		if !Contains([]string{
			"created", "opened", "reopened", "locked", "unlocked", "closed", "synchronize", // More to be added.
		}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createPullRequestText(event)
		markupText = "Open Pull Request"
		markupUrl = event.PullRequest.HTMLURL
	case "pull_request_review_comment":
		event := event.(*types.PullRequestReviewCommentEvent)

		if !Contains([]string{"created", "deleted"}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createPullRequestReviewCommentText(event)
		markupText = "Open Comment"
		markupUrl = event.Comment.HTMLURL
	case "push":
		event := event.(*types.PushEvent)
		// No Activity Types
		text = createPushText(event)
		markupText = "Open Changes"
		markupUrl = event.Compare
	case "release":
		event := event.(*types.ReleaseEvent)
		if !Contains([]string{"published", "released"}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createReleaseText(event)
		markupText = "🌐"
		markupUrl = event.Release.HTMLURL
	case "watch":
		event := event.(*types.WatchEvent)

		if !Contains([]string{"started"}, event.Action) {
			err = fmt.Errorf("unsupported event type '%s' for %s", event.Action, meta.EventName)
			return
		}

		text = createStarText(event)
		markupText = fmt.Sprintf("✨ Total stars: %d", event.Repo.StargazersCount)
		markupUrl = event.Repo.HTMLURL + "/stargazers"
	}
	return text, markupText, markupUrl, nil
}

func createPushText(event *PushEvent) string {
    // 1. Extrair e formatar o nome do repositório
    repoFullName := event.GetRepo().GetFullName() // Exemplo: "LogicCorporation/ProxyChecker-v2"
    repoParts := strings.Split(repoFullName, "/")
    repoName := repoParts[len(repoParts)-1]               // "ProxyChecker-v2"
    repoName = strings.ReplaceAll(repoName, "-", " ")      // "ProxyChecker v2"

    // 2. Remover a versão se necessário (opcional)
    // Por exemplo, remover "-v2" para ficar apenas "ProxyChecker"
    // repoName = strings.Split(repoName, " v")[0] // Descomente se quiser remover a versão

    // 3. Obter a data de publicação no formato [DD/MM/AA]
    pubDate := event.GetCreatedAt().Time.Format("02/01/06") // "DD/MM/AA"

    // 4. Cabeçalho com emoji e formatação em negrito e sublinhado
    header := fmt.Sprintf("🚀 <b><u>%d New Update(s) to %s</u></b>\n\n",
        len(event.Commits),
        repoName,
    )

    // 5. Seção de Atualizações com marcadores
    updatesText := "<b><u>📌 Updates:</u></b>\n"
    for _, commit := range event.Commits {
        // Escapar caracteres HTML nas mensagens de commit
        commitMessage := html.EscapeString(commit.GetMessage())
        // Adicionar marcador e mensagem do commit
        updatesText += fmt.Sprintf("• %s\n", commitMessage)
    }

    // 6. Rodapé com agradecimento e data
    footer := fmt.Sprintf("\nSpecial thanks to accompany, stay tuned for more. [ %s ]", pubDate)

    // Combinar todas as partes
    return header + updatesText + footer
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

