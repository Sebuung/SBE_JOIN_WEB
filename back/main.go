package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type StudentForm struct {
	Grade   string `json:"grade"`
	Class   string `json:"class"`
	Number  int    `json:"number"`
	Name    string `json:"name"`
}

func main() {
	http.HandleFunc("/submit-form", submitFormHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func submitFormHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    numberValue, err := strconv.Atoi(r.FormValue("number"))
    if err != nil {
        http.Error(w, "Invalid value for 'number'", http.StatusBadRequest)
        return
    }

    grade := r.FormValue("grade")
    class := r.FormValue("class")
    studentNumber := generateStudentNumber(grade, class, numberValue)

    form := StudentForm{
        Grade:   grade,
        Class:   class,
        Number:  numberValue,
        Name:    r.FormValue("name"),
    }

    err = saveFormToDynamoDB(form, studentNumber)
    if err != nil {
        log.Println("Error saving form data:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Set the thank you message
    thankYouMessage := "제출해주셔서 감사합니다!"

    // Create a template to display the alert message
    template := fmt.Sprintf(`
        <html>
        <head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">

		<title>Thank You</title>
			<script type="text/javascript" src="script.js" charset="utf-8">

			alert("%s");
                window.location.href = "http://sebuung.team";
            </script>
        </head>
        <body>
            <h1>%s</h1>
        </body>
        </html>
    `, thankYouMessage, thankYouMessage)

    // Set the response headers
    w.Header().Set("Content-Type", "text/html")

    // Write the response
    w.Write([]byte(template))
}


func generateStudentNumber(grade, class string, number int) string {
	gradeClass := grade + class
	studentNumber := fmt.Sprintf("%s%02d", gradeClass, number)
	return studentNumber
}

func saveFormToDynamoDB(form StudentForm, studentNumber string) error {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
        // Provide your AWS credentials or use environment variables, IAM roles, or other methods of authentication.
    })
    if err != nil {
        return err
    }

    svc := dynamodb.New(sess)

    item, err := dynamodbattribute.MarshalMap(form)
    if err != nil {
        return err
    }
    item["StudentNumber"] = &dynamodb.AttributeValue{S: aws.String(studentNumber)} // Add the StudentNumber attribute

    tableName := "user_tbl" // Replace with your DynamoDB table name

    input := &dynamodb.PutItemInput{
        Item:      item,
        TableName: aws.String(tableName),
    }

    _, err = svc.PutItem(input)
    if err != nil {
        return err
    }

    return nil
}
