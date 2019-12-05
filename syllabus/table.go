package syllabus

import (
	"errors"
	"io"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// Table 科目一覧の表のページに相当
type Table struct {
	Subjects []*ConciseSubject
	// PrevStartIndex [前のn件]のボタンのstartIndex
	PrevStartIndex string
	// NextStartIndex [次のn件]のボタンのstartIndex
	NextStartIndex string
}

// ConciseSubject 科目一覧の表の科目（情報がSubjectより少ない）
type ConciseSubject struct {
	// TimetableCode 時間割コード
	TimetableCode string
	// CourseTitle 科目名
	CourseTitle string
	// Semester 開講期
	Semester string
	// DayAndPeriod 開講時間帯
	DayAndPeriod string
	// Grade 対象学年
	Grade string
	// Credits 単位数
	Credits string
	// Update 更新日
	Update string
}

func makeTable(body io.Reader) (*Table, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	table := &Table{}

	selection := doc.Find("table#ichiran > tbody > .ichiran_odd,.ichiran_even")

	if selection == nil {
		return nil, errors.New("$(table#ichiran > tbody > .ichiran_odd,.ichiran_even) is nil")
	}

	for i := range selection.Nodes {
		tr := selection.Eq(i)
		row := []string{}
		tr.Find("td").Each(func(_ int, td *goquery.Selection) {
			row = append(row, td.Text())
		})

		if len(row) != 7 {
			return nil, errors.New("len(#ichiran > tbody > .ichiran_odd,.ichiran_even > td) is not 7")
		}

		table.Subjects = append(table.Subjects, &ConciseSubject{
			TimetableCode: row[0],
			CourseTitle:   row[1],
			Semester:      row[2],
			DayAndPeriod:  row[3],
			Grade:         row[4],
			Credits:       row[5],
			Update:        row[6],
		})
	}

	forms := doc.Find("table.ichiran_idou > tbody > tr > td > form")
	for i := range forms.Nodes {
		form := forms.Eq(i)
		parentClassName := form.Parent().AttrOr("class", "")
		switch parentClassName {
		case "left", "right":
			action, exist := form.Attr("action")
			if !exist {
				return nil, errors.New("one of ichiran_idou's children does not have attribute 'action'")
			}

			u, err := url.Parse(action)
			if err != nil {
				return nil, err
			}

			if parentClassName == "left" {
				table.NextStartIndex = u.Query().Get("s")
			} else {
				table.PrevStartIndex = u.Query().Get("s")
			}

		default:
			return nil, errors.New("class of one of ichiran_idou's children is neither 'left' or 'right'")
		}
	}

	return table, nil
}
