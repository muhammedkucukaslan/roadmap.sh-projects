package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	usageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println(errorStyle.Render("Error: Username not provided"))
		fmt.Println()
		fmt.Println(usageStyle.Render("Usage:"))
		fmt.Printf("    %s\n", commandStyle.Render("./main <username>"))
		return
	}
	username := args[1]
	client := NewClient(username)

	activities, err := client.makeRequest()
	if err != nil {
		fmt.Println(errorStyle.Render("Error:", err.Error()))
		return
	}
	if err := printActivites(username, activities); err != nil {
		fmt.Println(err)
		return
	}
}

func printActivites(username string, activities []activity) error {
	if len(activities) == 0 {
		return fmt.Errorf("no activity found")
	}
	fmt.Println(

		lipgloss.NewStyle().
			Bold(true).
			Padding(1, 0).
			Foreground(lipgloss.Color("12")).
			Render(fmt.Sprintf("%s's recent activity(s)", username)),
	)
	for _, event := range activities {
		var action string
		switch event.Type {
		case "PushEvent":
			commitCount := len(event.Payload.Commits)
			action = fmt.Sprintf("Pushed %d commit(s) to %s", commitCount, event.Repo.Name)
		case "PullRequestEvent":
			pullRequestAction := strings.ToTitle(event.Payload.Action[:1]) + (event.Payload.Action[1:])
			action = fmt.Sprintf("%s Pull Request at %s ", pullRequestAction, event.Repo.Name)
		case "IssuesEvent":
			action = fmt.Sprintf("%s an issue in %s", event.Payload.Action, event.Repo.Name)
		case "WatchEvent":
			action = fmt.Sprintf("Starred %s", event.Repo.Name)
		case "ForkEvent":
			action = fmt.Sprintf("Forked %s", event.Repo.Name)
		case "CreateEvent":
			action = fmt.Sprintf("Created %s in %s", event.Payload.Ref_Type, event.Repo.Name)
		default:
			action = fmt.Sprintf("%s in %s", event.Type, event.Repo.Name)
		}

		arrowStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#45475A"))

		textStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).
			Width(61).
			Foreground(lipgloss.Color("#CBA6F7"))

		actionStyle := fmt.Sprintf("%s %s",
			arrowStyle.Render("▶"),
			textStyle.Render(action))

		fmt.Println(actionStyle)
		actionStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("#3C3C3C")).
			Render(fmt.Sprintf("- %s", action))

	}
	return nil
}
