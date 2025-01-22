package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	today    = time.Now().Format("2006-01-02")
	jsonPath = "data.json"
)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type List []*Task

func (l *List) readJsonFile() error {

	file, err := os.OpenFile(jsonPath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer file.Close()
	if info, err := file.Stat(); err != nil {
		return fmt.Errorf("%v", err)
	} else if info.Size() == 0 {
		*l = make([]*Task, 0)
		return nil
	}

	var todoBytes []byte

	todoBytes, err = io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err = json.Unmarshal(todoBytes, l); err != nil {
		return fmt.Errorf("Unmarshalling error: %s", err)
	}
	return nil
}

func (l *List) writeJsonFile() {
	jsonByte, err := json.MarshalIndent(l, "", " ")
	if err != nil {
		fmt.Println("Error parsing data to json:", err)
		return
	}

	if err := os.WriteFile(jsonPath, jsonByte, 0644); err != nil {
		fmt.Println("Error writing file:")
		return
	}
}

func (l *List) add(description string) {
	var id int
	if len((*l)) > 0 {
		id = (*l)[len((*l))-1].Id + 1
	} else {
		id = 1
	}
	task := Task{
		Id:          id,
		Description: description,
		Status:      "todo",
		CreatedAt:   today,
		UpdatedAt:   today,
	}
	*l = append(*l, &task)
}

func (l List) list() []Task {
	var tasks []Task
	for _, task := range l {
		tasks = append(tasks, *task)
	}
	return tasks
}

func (l List) listDone() []Task {
	var tmp []Task
	for _, task := range l {
		if task.Status == "done" {
			tmp = append(tmp, *task)
		}
	}
	return tmp
}

func (l List) listTodo() []Task {
	var tmp []Task
	for _, task := range l {
		if task.Status == "todo" {
			tmp = append(tmp, *task)
		}
	}
	return tmp
}
func (l List) listInProgress() []Task {
	var tmp []Task
	for _, task := range l {
		if task.Status == "in-progress" {
			tmp = append(tmp, *task)
		}
	}
	return tmp
}

func (l List) update(id int, description string) error {
	for i, task := range l {
		if task.Id == id {
			(*l[i]).Description = description
			(*l[i]).UpdatedAt = today
			return nil
		}
	}
	return fmt.Errorf("Invalid index")

}

func (l *List) delete(id int) error {
	if len(*l) == 1 {
		*l = make(List, 0)
		return nil
	}
	for i, task := range *l {
		if task.Id == id {
			*l = append((*l)[:i], (*l)[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Invalid index")
}

func (l *List) deleteDone() error {
	newSlice := make([]*Task, 0, len(*l))
	for i, task := range *l {
		if task.Status != "done" {
			newSlice = append(newSlice, (*l)[i])
		}
	}
	*l = newSlice
	return nil
}

// we can add deleteAll() ?

func (l *List) markDone(id int) error {
	for i, task := range *l {
		if task.Id == id {
			(*l)[i].Status = "done"
			(*l)[i].UpdatedAt = today
			return nil
		}
	}
	return fmt.Errorf("Invalid index")
}
func (l *List) markInProgress(id int) error {
	for i, task := range *l {
		if task.Id == id {
			(*l)[i].Status = "in-progress"
			(*l)[i].UpdatedAt = today
			return nil
		}
	}
	return fmt.Errorf("Invalid index")
}
