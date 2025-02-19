package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
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
		________________
		1. View tasks
		2. Add task
		3. Modify task progress
		4. Remove task
		5. Remove all tasks
		6. Exit
		_________________
		Select an option: `

	fmt.Println(menu)
}
func clearConsole() {

	switch runtime.GOOS {
	case "windows":
		fmt.Print("\033[H\033[J")
	case "linux", "darwin":
		fmt.Print("\033[H\033[J")
	default:
		fmt.Print("\n")
	}
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
		fmt.Println("No tasks to display! ")
		return
	}

	clearConsole()
	fmt.Println("Tasks:")
	fmt.Println("_______")
	fmt.Print("")
	for i, task := range tasks {
		fmt.Printf(" << %d Description: %s, Status: %s >>\n", i+1, task.Description, task.Status)
	}
	fmt.Println("_______")

}

func addTask() {

	clearConsole()
	fmt.Print("Add a description for the new task: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		description := scanner.Text() // Get the scanned text and assign it to description

		if taskExists(description) {
			fmt.Println("Task with such description already exists:", description)
			return
		}

		tasks = append(tasks, Task{Description: description, Status: NotStarted})
		saveTasksToFile()
		fmt.Println("Task added successfully:", description)
	} else {
		fmt.Println("Error reading input:", scanner.Err())
	}
}
func modifyTaskStatus() {
	clearConsole()
	displayTasks()

	if len(tasks) == 0 {
		return
	}

	fmt.Print("Select the task to modify status (1-", len(tasks), "): ")
	var taskIndex int
	_, err := fmt.Scanln(&taskIndex)
	if err != nil || taskIndex < 1 || taskIndex > len(tasks) {
		fmt.Println("Invalid task index. Please choose a number between 1 and", len(tasks))
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
	if err != nil || newStatus < 1 || newStatus > 3 {
		fmt.Println("Invalid status. Please choose a number between 1 and 3")
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
	clearConsole()
	fmt.Println("Task status successfully modified")
	saveTasksToFile()
}

func removeTask() {
	clearConsole()
	displayTasks()

	if len(tasks) == 0 {
		return
	}

	fmt.Printf("Select the task to remove (1-%d): ", len(tasks))

	var taskIndex int
	_, err := fmt.Scanln(&taskIndex)
	if err != nil || taskIndex < 1 || taskIndex > len(tasks) {
		fmt.Println("Invalid task index. Please choose a number between 1 and", len(tasks))
		return
	}

	taskIndex -= 1

	tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)
	saveTasksToFile()
	fmt.Println("Task removed successfully")
}

func removeAllTasks() {
	clearConsole()
	tasks = []Task{}
	saveTasksToFile()
	fmt.Println("All tasks removed successfully")
}

func exit() {
	clearConsole()
	fmt.Println("Bye!")
	os.Exit(0)
}

func main() {
	loadTasksFromFile()

	for {
		displayMenu()

		var option int
		_, err := fmt.Scan(&option)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number between 1 and 6")
			// Clear the input buffer
			bufio.NewReader(os.Stdin).ReadString('\n')
			continue
		}

		// Validate menu option range
		if option < 1 || option > 6 {
			fmt.Println("Invalid option. Please enter a number between 1 and 6")
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
			removeAllTasks()
		case 6:
			exit()
		}

		// Save after any operation
		if option != 6 { // Don't save if exiting
			saveTasksToFile()
		}
	}
}
