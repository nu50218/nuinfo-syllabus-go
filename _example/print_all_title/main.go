package main

import (
	"fmt"
	"github.com/nu50218/nuinfo-syllabus-go/syllabus"
	"log"
	"time"
)

const Endpoint string = "https://syllabus.i.nagoya-u.ac.jp/i/"

func main() {
	client := syllabus.NewClient(Endpoint, 1*time.Second)
	subjects, err := client.GetAllConciseSubjects()
	if err != nil {
		log.Fatalln(err)
	}
	for _, subject := range subjects {
		fmt.Println(subject.CourseTitle)
	}
}
