package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) (string, error) {
	for _, m := range sqsEvent.Records {
		var userRequest User
		if err := json.Unmarshal([]byte(m.Body), &userRequest); err != nil {
			fmt.Printf("error %v", err)
		}

		saveDataToDB(&userRequest)
		saveFile(&userRequest)
		readFile()
	}
	return "Successful", nil
}

func main() {
	lambda.Start(handler)
}

func saveDataToDB(user *User) {
	db, _ := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/test")
	db.Query("INSERT INTO test VALUES (?, ?);", &user.Username, &user.Email)
}

func saveFile(user *User) {
	fileContent := []byte(user.Username + " " + user.Email)

	filePath := filepath.Join("/tmp/", "output.txt")
	err := os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("File created and saved successfully")
}

func readFile() {
	f, err := os.ReadFile("/tmp/output.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("File output: %v \n", string(f))
}

//docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root --restart unless-stopped mysql:8
