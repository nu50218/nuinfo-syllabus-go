package syllabus

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Client シラバスからデータを取得するクライアント
type Client struct {
	Endpoint string
	interval time.Duration
	mutex    *sync.Mutex
}

// NewClient クライアントを作って返す
func NewClient(endpoint string, interval time.Duration) *Client {
	return &Client{
		Endpoint: endpoint,
		mutex:    &sync.Mutex{},
		interval: interval,
	}
}

// GetFirstTable 開いたときに表示されるはじめの表を取得する
func (client *Client) GetFirstTable() (*Table, error) {
	return client.GetTable("")
}

// GetTable startIndex（URLクエリパラメータのs）を指定して表を取得する
func (client *Client) GetTable(startIndex string) (*Table, error) {
	u, err := url.Parse(client.Endpoint)
	if err != nil {
		return nil, err
	}

	if startIndex != "" {
		q := u.Query()
		q.Add("s", startIndex)
		u.RawQuery = q.Encode()
	}

	client.mutex.Lock()
	defer func() {
		time.Sleep(client.interval)
		client.mutex.Unlock()
	}()
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return makeTable(res.Body)
}

// GetAllConciseSubjects すべての表を取得して科目の列にして返す
func (client *Client) GetAllConciseSubjects() ([]*ConciseSubject, error) {
	subjects := []*ConciseSubject{}

	table, err := client.GetFirstTable()
	if err != nil {
		return []*ConciseSubject{}, err
	}
	subjects = append(subjects, table.Subjects...)

	for table.NextStartIndex != "" {
		if table, err = client.GetTable(table.NextStartIndex); err != nil {
			return []*ConciseSubject{}, err
		}
		subjects = append(subjects, table.Subjects...)
	}

	return subjects, nil
}

// GetSubject timetableCode（URLクエリパラメータのj、時間割コード）を指定してそのページを取得・パースして返す
func (client *Client) GetSubject(timetableCode string) (*Subject, error) {
	u, err := url.Parse(client.Endpoint)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("j", timetableCode)
	u.RawQuery = q.Encode()

	client.mutex.Lock()
	defer func() {
		time.Sleep(client.interval)
		client.mutex.Unlock()
	}()
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return makeSubject(u.String(), timetableCode, res.Body)
}
