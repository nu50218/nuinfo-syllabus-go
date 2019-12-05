package syllabus

import (
	"errors"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Subject 科目の詳細のページに相当
type Subject struct {
	// CourseTitle 科目名
	CourseTitle string
	// Update 更新日
	Update string
	// TimetableCode 時間割コード
	TimetableCode string
	// URL 科目詳細ページへのURL
	URL string
	// JapaneseValues 科目詳細ページの日本語の小見出しと内容のmap
	JapaneseValues map[string]string
	// EnglishValues 科目詳細ページの英語の小見出しと内容のmap
	EnglishValues map[string]string
}

const updatePrefix string = "更新日："

func makeSubject(url, timetableCode string, body io.Reader) (*Subject, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	subject := &Subject{
		URL:           url,
		TimetableCode: timetableCode,
	}

	// title
	subject.CourseTitle = doc.Find("h2#detail_midashi").Text()
	if subject.CourseTitle == "" {
		return nil, errors.New("title is empty")
	}

	// update
	updateText := doc.Find("table#detail_midashi > tbody > tr").Children().Last().Text()
	if !strings.Contains(updateText, updatePrefix) {
		return nil, errors.New("'更新日：' is not contained")
	}
	subject.Update = strings.TrimPrefix(updateText, "更新日：")

	// values
	parseDiv := func(div *goquery.Selection) map[string]string {
		values := make(map[string]string)
		currentKey := ""
		currentValue := ""
		children := div.Children()
		for i := range children.Nodes {
			sel := children.Eq(i)
			switch goquery.NodeName(sel) {
			case "h4":
				if currentKey != "" {
					values[currentKey] = currentValue
				}
				currentKey = sel.Text()
				currentValue = ""
			default:
				currentValue += sel.Text()
			}
		}
		return values
	}

	subject.JapaneseValues = parseDiv(doc.Find("div#japanese"))
	subject.EnglishValues = parseDiv(doc.Find("div#english"))

	if len(subject.JapaneseValues) == 0 {
		return nil, errors.New("could not get japanese values")
	}

	return subject, nil
}
