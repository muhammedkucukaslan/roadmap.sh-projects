package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

func displayMovies(movies Movies) {

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFA500")).
		MarginBottom(1)

	ratingStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#2E8B57")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)

	dateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4682B4")).
		Italic(true)

	overviewHeaderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9370DB")).
		Bold(true)

	overviewStyle := lipgloss.NewStyle().
		MaxWidth(60).
		Foreground(lipgloss.Color("#B8B8B8")).
		MarginLeft(2).
		MarginTop(1).
		MarginBottom(1)

	adultBadgeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#FF0000")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)

	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444"))

	var output strings.Builder

	for i, movie := range movies.Results {
		date, _ := time.Parse("2006-01-02", movie.ReleaseDate)

		formattedDate := date.Format("January 2, 2006")

		ratingText := fmt.Sprintf("★ %.1f/10", movie.Rating)

		title := fmt.Sprintf("%d. %s", i+1, movie.Title)

		movieBlock := titleStyle.Render(title) + "\n"

		movieBlock += ratingStyle.Render(ratingText) + " " +
			dateStyle.Render(formattedDate)

		if movie.Adult {
			movieBlock += " " + adultBadgeStyle.Render("18+")
		}
		movieBlock += "\n" + overviewHeaderStyle.Render("Overview:") + "\n" +
			overviewStyle.Render(wrapText(movie.Overview, 58))

		movieBlock += fmt.Sprintf("\nOriginal Language: %s", strings.ToUpper(movie.OriginalLanguage))

		if i < len(movies.Results)-1 {
			movieBlock += "\n" + dividerStyle.Render(strings.Repeat("─", 70)) + "\n"
		}

		output.WriteString(movieBlock + "\n")
	}

	fmt.Println(output.String())
}

func wrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}

func PrintUsage() {
	yellow := color.New(color.FgHiYellow)
	magenta := color.New(color.FgHiMagenta)

	yellow.Printf("Usage:\n")
	magenta.Printf("\t\t tmdb-cli --type <type>\n\n")

	yellow.Printf("Types:\n")
	for _, t := range validTypes {
		magenta.Printf("\t\t %s\n", t)
	}
}

func ValidateArgs(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("incorrect number of arguments")
	}
	if args[1] != "--type" {
		return fmt.Errorf("second argument must be --type")
	}
	if !slices.Contains(validTypes, args[2]) {
		return fmt.Errorf("invalid type: %s", args[2])
	}
	return nil
}

func HandleError(err error) {
	red := color.New(color.FgHiRed)
	red.Printf("Error: %v\n\n", err)
	PrintUsage()
	os.Exit(1)
}

func getProperParam(str string) string {
	switch str {
	case "playing":
		return "now_playing"
	case "popular":
		return "popular"
	case "top":
		return "top_rated"
	case "upcoming":
		return "upcoming"
	}
	return ""
}
