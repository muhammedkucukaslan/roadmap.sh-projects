package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

var (
	layout = "2006-01-02"
	header = fmt.Sprintf("%-*s  %-*s  %-*s %-*s",
		idWidth, "ID", dateWidth, "Date",
		descWidth, "Description",
		amountWidth, "Amount",
	)
	idWidth     = 5
	descWidth   = 15
	dateWidth   = 7
	amountWidth = 5
)

type Expense struct {
	Id          int     `csv:"id"`
	Date        string  `csv:"date"`
	Description string  `csv:"description"`
	Amount      float64 `csv:"amount"`
}

type Report []*Expense

func NewReport() *Report {
	return &Report{}
}

func (r *Report) readFromCSV(sourcePath string) error {
	file, err := os.OpenFile(sourcePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("opening file error: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("getting file info error: %w", err)
	}

	if info.Size() == 0 {
		*r = Report{}
		return nil
	}
	if err := gocsv.UnmarshalFile(file, r); err != nil {
		return fmt.Errorf("unmarshaling file error: %w", err)
	}

	return nil
}

func (r *Report) writeToCSV(sourcePath string) error {

	file, err := os.OpenFile(sourcePath, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Opening File Error:%s", err)
	}
	defer file.Close()
	if err := gocsv.Marshal(r, file); err != nil {
		return fmt.Errorf("Writing File Error:%s", err)
	}
	return nil
}

func (r *Report) add(description string, amount float64) (int, error) {
	var id int
	if len(*r) > 0 {
		id = (*r)[len(*r)-1].Id + 1
	} else {
		id = 1
	}
	*r = append(*r, &Expense{
		Id:          id,
		Description: description,
		Amount:      amount,
		Date:        time.Now().Format(layout),
	})
	return id, nil
}

func (r *Report) print() {
	fmt.Println(header)
	for _, task := range *r {
		fmt.Println(fmt.Sprintf("%-*d  %-*s  %-*s %-*.2f",
			idWidth, task.Id, dateWidth, task.Date,
			descWidth, task.Description,
			amountWidth, task.Amount,
		))
	}
}

func (r *Report) delete(id int) error {
	var found bool
	for i, expense := range *r {
		if expense.Id != id {
			*r = append((*r)[:i], (*r)[i+1:]...)
			found = true
		}
	}
	if !found {
		return fmt.Errorf("expense with id %d not found", id)
	}
	return nil
}

func (r *Report) summary() float64 {
	var sum float64
	for _, expense := range *r {
		sum += expense.Amount
	}
	return sum
}

func (r *Report) summaryOfMonth(order int) (string, float64, error) {
	if order <= 0 || order > 12 {
		return "", 0, fmt.Errorf("Month can not be less than 12 or a negative number")
	}

	var sum float64
	month := time.Month(order)
	for _, expense := range *r {
		parsedDate, err := time.Parse(layout, expense.Date)
		if err != nil {
			return "", 0, fmt.Errorf("Invalid date: %v", err)
		}

		if month == parsedDate.Month() {
			sum += expense.Amount
		}
	}
	return month.String(), sum, nil
}
