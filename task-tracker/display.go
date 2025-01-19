package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	usage      = lipgloss.NewStyle().Foreground(lipgloss.Color("#fcba03")).Render("Usage:")
	operations = lipgloss.NewStyle().Foreground(lipgloss.Color("#fcba03")).Render("Operations:")
	filters    = lipgloss.NewStyle().Foreground(lipgloss.Color("#fcba03")).Render("Filters:")
	seperator  = lipgloss.NewStyle().Foreground(lipgloss.Color("#25b4e8"))
	horizontal = strings.Repeat("-", 70)
	header     = fmt.Sprintf("%-*s %s %-*s %s %-*s",
		idWidth, "ID",
		seperator.Render("|"),
		descWidth, "Description",
		seperator.Render("|"),
		dateWidth, "Date")
	idWidth   = 5
	descWidth = 40
	dateWidth = 15
)

func displayList(tasks []Task) {
	fmt.Println(header)
	fmt.Println(seperator.Render(horizontal))
	for _, task := range tasks {
		taskLine := fmt.Sprintf("%-*d %s %-*s %s %-*s",
			idWidth, task.Id,
			seperator.Render("|"),
			descWidth, task.Description,
			seperator.Render("|"),
			dateWidth, task.CreatedAt)
		fmt.Println(taskLine)
		fmt.Println(seperator.Render(horizontal))
	}
}

func checkArgs(args []string) int {
	ln := len(args)
	switch {
	case ln == 1:
		displayOperations()
		return 1
	case ln == 2, ln > 2:
		switch args[1] {
		case "add":
			if ln != 3 {
				fmt.Println("Invalid usage. Use quotation marks (\"\")")
				fmt.Println(usage)
				fmt.Println("\ttask-cli add <task>")
				return 1
			}
		case "update":
			if ln != 4 {
				fmt.Println("Invalid usage")
				fmt.Println(usage)
				fmt.Println("\ttask-cli update <id> <task>")
				return 1
			}
		case "delete-done":
			if ln > 2 {
				fmt.Println("Invalid usage")
				fmt.Println(usage)
				fmt.Printf("\ttask-cli delete-done \n")
				return 1
			}
		case "delete", "mark-done", "mark-in-progress":
			if ln != 3 {
				fmt.Println("Invalid usage")
				fmt.Println(usage)
				fmt.Printf("\ttask-cli %s <id>\n", args[1])
				return 1
			}
		case "list":
			if ln > 3 || (ln == 3 && (args[2] != "done" && args[2] != "todo" && args[2] != "in-progress")) {
				fmt.Println("Invalid usage")
				fmt.Println(usage)
				fmt.Printf("\ttask-cli list <filter>\n")
				fmt.Println(filters)
				fmt.Println("\t done")
				fmt.Println("\t todo")
				fmt.Println("\t in-progress")

				return 1
			}
		default:
			fmt.Println("Invalid operation:")
			displayOperations()
			return 1
		}
	}
	return 0
}

func displayOperations() {
	fmt.Println(usage)
	fmt.Println("\ttask-cli <operation>")
	fmt.Println(operations)
	fmt.Println("\t add <task>")
	fmt.Println("\t list <filter>")
	fmt.Println("\t delete <id>")
	fmt.Println("\t delete-done")
	fmt.Println("\t update <id> <task>")
	fmt.Println("\t mark-in-progress <id>")
	fmt.Println("\t mark-done <id>")
}
