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
func taskExists(description string) bool {
	for _, t := range tasks {
		if description == t.Description {
			return true
		}
	}
	return false
}

func displayMenu() {
	menu :=
		`	Main Menu:
		1. View tasks
		2. Add task
		3. Modify task progress
		4. Remove task
		5. Exit
		Select an option: `

	fmt.Println(menu)
}

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
			return
		}
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Println("Error decoding tasks:", err)
	}
}
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
			return
		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("Error encoding tasks:", err)
		return
	}
}

func displayTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks to display")
	}

	for i, task := range tasks {
		fmt.Printf("<< %d Description: %s, Status: %s >>\n", i+1, task.Description, task.Status)
	}
}

func addTask() {
	fmt.Print("Add a description for the new task: ")
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

func modifyTaskStatus() {
	displayTasks()

	fmt.Print("Select the task to modify status: ")
	var taskIndex int
	_, err := fmt.Scanln(&taskIndex)
	if err != nil || taskIndex < 0 || taskIndex > len(tasks) {
		fmt.Println("Invalid task selected", err)
		return
	}

	fmt.Print(`
		Choose new status:
		1. Not Started
		2. In Progress
		3. Complete
	`)
	var newStatus int
	_, err = fmt.Scanln(&newStatus)
	if err != nil {
		fmt.Println("Incorrect status selected. Try again:", err)
		return
	}

	switch newStatus {
	case 1:
		tasks[taskIndex-1].Status = NotStarted
	case 2:
		tasks[taskIndex-1].Status = InProgress
	case 3:
		tasks[taskIndex-1].Status = Completed
	}
	fmt.Println("Task status successfully modified")
	saveTasksToFile()

}

func removeTask() {
	displayTasks()
	fmt.Println("Select the task to remove: ")

	var taskIndex int
	_, err := fmt.Scanln(&taskIndex)
	if err != nil || taskIndex < 0 || taskIndex > len(tasks) {
		fmt.Println("Invalid task selected", err)
	}

	tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)
	saveTasksToFile()
	fmt.Println("Task removed successfully:")
}

func showCompletedTasks() {}

func exit() {
	os.Exit(0)
}

func main() {
	loadTasksFromFile()

	for {
		displayMenu()

		var option int
		_, err := fmt.Scan(&option)
		if err != nil {
			fmt.Println("Invalid input, please try again :", err)
			continue
		}

		switch option {
		case 1:
			displayTasks()
		case 2:
			addTask()
		case 3:
			modifyTaskStatus()
		case 4:
			removeTask()
		case 5:
			exit()
		default:
			saveTasksToFile()
		}
	}
}
