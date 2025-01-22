package main

import (
	"fmt"
	"os"
	"strconv"
)

var (
	jsonPath = "data.json"
)

func main() {
	args := os.Args
	if err := checkArgs(args); err == 1 {
		return
	}
	list := List{}
	if err := list.readJsonFile(jsonPath); err != nil {
		fmt.Println(err)
		return
	}

	switch args[1] {
	case "add":
		task := args[2]
		list.add(task)
	case "list":
		if len(args) > 2 {
			switch args[2] {
			case "done":
				displayList(list.listDone())
			case "in-progress":
				displayList(list.listInProgress())
			case "todo":
				displayList(list.listTodo())
			}
			break
		}
		displayList(list.list())
	case "update":
		id, _ := strconv.Atoi(args[2])
		newTask := args[3]
		if err := list.update(id, newTask); err != nil {
			fmt.Println(err)
			return
		}
	case "delete":
		id, _ := strconv.Atoi(args[2])
		if err := list.delete(id); err != nil {
			fmt.Println(err)
			return
		}
	case "delete-done":
		if err := list.deleteDone(); err != nil {
			fmt.Println(err)
			return
		}
	case "mark-in-progress":
		id, _ := strconv.Atoi(args[2])
		if err := list.markInProgress(id); err != nil {
			fmt.Println(err)
			return
		}
	case "mark-done":
		id, _ := strconv.Atoi(args[2])
		if err := list.markDone(id); err != nil {
			fmt.Println(err)
			return
		}
	}

	if err := list.writeJsonFile(jsonPath); err != nil {
		fmt.Println(err)
		return
	}
}
