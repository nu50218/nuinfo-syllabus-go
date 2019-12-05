package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nu50218/nuinfo-syllabus-go/syllabus"
)

const Endpoint string = "https://syllabus.i.nagoya-u.ac.jp/i/"

const keyTopic string = "授業概要"

func main() {
	client := syllabus.NewClient(Endpoint, 1*time.Second)

	table, err := client.GetFirstTable()
	if err != nil {
		log.Fatalln(err)
	}

	subject, err := client.GetSubject(table.Subjects[0].TimetableCode)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(subject.CourseTitle)
	fmt.Println(subject.JapaneseValues[keyTopic])
}
