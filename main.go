package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	users := generateUsers(100)

	var mg sync.WaitGroup
	mg.Add(1)
	saveUserInfo(users, &mg)
	mg.Wait()

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func saveUserInfo(users []User, mg *sync.WaitGroup) {
	task := make(chan User)

	for i := 0; i < 10; i++ {
		go saveUser(task)
	}

	go func() {
		for _, user := range users {
			task <- user
		}
		mg.Done()
	}()
}

func saveUser(task <-chan User) {
	for user := range task {
		fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

		filename := fmt.Sprintf("users/uid%d.txt", user.id)
		fmt.Println(filename)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.WriteString(user.getActivityInfo())
		if err != nil {
			log.Fatal(err)
		}
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second)
	}
}

func generateUsers(count int) []User {
	users := make([]User, count)
	task := make(chan int, 10)
	res := make(chan User, 10)

	for i := 0; i < 10; i++ {
		go makeUser(task, res)
	}

	go func() {
		for i := 0; i < count; i++ {
			task <- i
		}
		defer close(task)
	}()

	for i := 0; i < count; i++ {
		users[i] = <-res
	}
	defer close(res)

	return users
}

func makeUser(task <-chan int, res chan<- User) {
	for i := range task {
		user := User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@company.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("generated user %d\n", i+1)
		time.Sleep(time.Millisecond * 100)

		res <- user
	}
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
