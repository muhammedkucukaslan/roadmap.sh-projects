package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func printWelcome() {
	fmt.Println("Welcome to the Number Guessing Game!")
	fmt.Println("I'm thinking of a number between 1 and 100.")
}

func printThinkingAgain() {
	fmt.Println("Happy to see you again")
	fmt.Println("I'm thinking of a number between 1 and 100.")
}

func printDifficultyChoises() {
	fmt.Println("Please select the difficulty level:")
	color.New(color.FgGreen).Fprintln(os.Stdout, "1. Easy (10 chances)")
	color.New(color.FgYellow).Fprintln(os.Stdout, "2. Medium (5 chances)")
	color.New(color.FgRed).Fprintln(os.Stdout, "3. Hard (3 chances)")
}

func determineRight(choise int) int {
	switch choise {
	case int(Hard):
		return 3
	case int(Medium):
		return 5
	case int(Easy):
		return 10
	}
	return -1
}

func printChosingCelebrating(choise int) {
	switch choise {
	case int(Hard):
		fmt.Println("Great! You have selected the Hard difficulty level.")
	case int(Medium):
		fmt.Println("Great! You have selected the Medium difficulty level.")
	case int(Easy):
		fmt.Println("Great! You have selected the Easy difficulty level.")
	}
}

func compareNumbers(target int, guess int) bool {
	switch {
	case target > guess:
		fmt.Printf("Incorrect! The number is bigger than %d.\n", guess)
		return false
	case target < guess:
		fmt.Printf("Incorrect! The number is less than %d.\n", guess)
		return false
	}
	return true
}

func printingResult(win bool, attempt int) {
	if win {
		color.New(color.FgGreen).Printf("Congratulations! You guessed the correct number in %d attempts.\n", attempt)
		return
	}
	color.New(color.FgCyan).Println("Unfortunetly, you lost the game! Hide up your head champ!")
}
