package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Person struct {
	Name    string `csv:"Name"`
	Surname string `csv:"Surname"`
	Age     int    `csv:"Age"`
}

var (
	sourceFile = "report.csv"
)

func main() {

	args := os.Args

	//! Todo : Add better argumant validation
	if len(args) == 1 {
		color.New(color.FgHiYellow).Fprintln(os.Stdout, "Usage:")
		color.New(color.FgBlue).Fprintln(os.Stdout, "\t- expense-tracker <operation>")
		return
	}
	// Something like this would be grate
	//if err := checkArgs(args); err!=nil{
	//		fmt.Println(err)
	// 		return
	//}

	report := NewReport()
	if err := report.readFromCSV(sourceFile); err != nil {
		fmt.Println(err)
		return
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	description := addCmd.String("description", "", "expense itself")
	amount := addCmd.Float64("amount", -1, "amount you spent")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	id := deleteCmd.Int("id", -1, "id of expense")
	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	month := summaryCmd.Int("month", -1, "the order of the month you want to know how much you spend ")

	switch args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		id, err := report.add(*description, *amount)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Expense added successfully (ID: %d)\n", id)
	case "list":
		report.print()
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if err := report.delete(*id); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Expense deleted successfully")
	case "summary":
		summaryCmd.Parse(os.Args[2:])
		if *month != -1 {
			monthStr, sum, err := report.summaryOfMonth(*month)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Total expenses for %s: $%.3f\n", monthStr, sum)
			break
		}
		fmt.Printf("Total expenses: $%.3f\n", report.summary())
	}

	if err := report.writeToCSV(sourceFile); err != nil {
		fmt.Println(err)
		return
	}
}
