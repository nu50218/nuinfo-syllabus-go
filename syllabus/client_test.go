package syllabus

import (
	"testing"
	"time"
)

var testURL = []string{
	"https://syllabus.i.nagoya-u.ac.jp/i/",
	"https://syllabus.i.nagoya-u.ac.jp/gsi/",
}

func makeTestClients() []*Client {
	res := []*Client{}
	for _, url := range testURL {
		res = append(res, NewClient(url, 1*time.Second))
	}
	return res
}

func TestClient_GetTable(t *testing.T) {
	for _, client := range makeTestClients() {
		table, err := client.GetFirstTable()
		if err != nil {
			t.Error(err)
			return
		}
		if table.NextStartIndex == "" {
			return
		}

		table, err = client.GetTable(table.NextStartIndex)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = client.GetTable(table.PrevStartIndex)
		if err != nil {
			t.Error(err)
		}
		return
	}
}

func TestClient_GetAllConciseSubjects(t *testing.T) {
	for _, client := range makeTestClients() {
		if _, err := client.GetAllConciseSubjects(); err != nil {
			t.Error(err)
		}
	}
}

func TestClient_GetSubject(t *testing.T) {
	for _, client := range makeTestClients() {
		table, err := client.GetFirstTable()
		if err != nil {
			t.Error(err)
			return
		}
		for i, subject := range table.Subjects {
			if i > 5 {
				break
			}
			if _, err := client.GetSubject(subject.TimetableCode); err != nil {
				t.Error(err)
				return
			}
		}
	}
}
