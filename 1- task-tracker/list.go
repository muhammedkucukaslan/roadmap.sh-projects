package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var (
	today = time.Now().Format("2006-01-02")
)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type List []*Task

func (l *List) readJsonFile(jsonSource string) error {

	file, err := os.OpenFile(jsonSource, os.O_CREATE|os.O_RDONLY, 0644)
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

	if err := json.NewDecoder(file).Decode(l); err != nil {
		return fmt.Errorf("Reading file error:%w", err)
	}
	return nil
}

func (l *List) writeJsonFile(jsonSource string) error {
	file, err := os.OpenFile(jsonSource, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("Opening file error:%w", err)
	}
	defer file.Close()
	if err = json.NewEncoder(file).Encode(&l); err != nil {
		return fmt.Errorf("Writing file error:%w", err)
	}
	return nil
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
