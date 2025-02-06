package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Task struct {
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
}

var tasks []Task

type TaskStatus int

const (
	NotStarted TaskStatus = iota
	InProgress
	Completed
)

func (status TaskStatus) String() string {
	return [...]string{"NotStarted", "InProgress", "Completed"}[status]
}

// helper functions

func taskExists(description string) bool {
	for _, t := range tasks {
		if description == t.Description {
			return true
		}
	}
	return false
}
func displayMenu() {
	menu := `
	Main Menu:
	1. View tasks
	2. Add task
	3. Modify task progress
	4. Remove task
	5. Show completed tasks
	6. Exit
	Select an option: `
	fmt.Println(menu)
}

// main functions

func loadTasksFromFile() {
	file, err := os.Open("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Tasks cannot be loaded - no file found:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Println("Error decoding tasks:", err)
	}
}
func addTask() {
	fmt.Print("Add a description for the new task")
	var description string
	_, err := fmt.Scanln(&description)
	if err != nil {
		fmt.Println("Error reading description:", err)
		return
	}
	if taskExists(description) {
		fmt.Println("Task with such description already exists:", description)
		return
	}

	tasks = append(tasks, Task{description, NotStarted})
	saveTasksToFile()
	fmt.Println("Task added successfully:", description)
}
func modifyTaskProgress() {}
func removeTask()         {}
func showCompletedTasks() {}
func exit()               {}
func saveTasksToFile() {
	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("Error encoding tasks:", err)
		return
	}
}

func main() {
	loadTasksFromFile()

	for {
		displayMenu()

		var option int
		_, err := fmt.Scan(&option)
		if err != nil {
			fmt.Println("Invalid input, please try again :", err)
		}

		switch option {
		case 1:
			loadTasksFromFile()
		case 2:
			addTask()
		case 3:
			modifyTaskProgress()
		case 4:
			showCompletedTasks()
		case 5:
			removeTask()
		case 6:
			exit()
		default:
			saveTasksToFile()
		}
	}
}
