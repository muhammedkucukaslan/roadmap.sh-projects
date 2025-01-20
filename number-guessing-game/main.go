package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/fatih/color"
)

type Choises int

const (
	Easy Choises = iota + 1
	Medium
	Hard
)

func main() {
	printWelcome()
	var gameAttempt int
	for {
		if gameAttempt > 0 {
			printThinkingAgain()
		}
		start := time.Now()
		target := (rand.IntN(100))
		fmt.Println("TARGET:", target)
		var right, choise, attempt, guess int
		var win bool
		for {
			printDifficultyChoises()
			fmt.Print("Enter your choise: ")
			fmt.Scanln(&choise)
			if right = determineRight(choise); right != -1 {
				printChosingCelebrating(choise)
				break
			}
			color.New(color.FgRed).Printf("\nInvalid Choise. Try again")
		}

		fmt.Printf("Let's start the game!\n\n")
		for right > 0 {
			fmt.Print("Enter your guess:")
			fmt.Scanln(&guess)
			if guess < 0 || guess > 100 {
				fmt.Println("Invalid number! Number should be between 0 and 100")
				continue
			}
			attempt++
			if compareNumbers(target, guess) {
				win = true
				break
			}
			right--
		}
		printingResult(win, attempt)

		var again string
		elapsed := time.Since(start)
		fmt.Printf("You finished the game in the %.3v sec\n", elapsed)
		fmt.Println("Do you want to play one more time? (N for NO)")
		fmt.Scanln(&again)
		if again == "N" || again == "n" {
			fmt.Println("Have a good day!")
			break
		}
		gameAttempt++
	}

}
