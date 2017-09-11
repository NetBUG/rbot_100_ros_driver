package client

import (
	"fmt"
	"github.com/CiscoZeus/zeus-analytics/registry/model"
	"github.com/franela/goreq"
	"net/http"
	//"strings"
)

// TODO func NewClient() *Client {}

type Client struct {
	Host string
	Port string
}

//TBD: Define App here, currenty takes defintion from model
type Job struct {
	AppID       string
	Name        string
	Description string
	Type        string
	Action      string
	Parameters  string
	Input       string
	Output      string
}
type Data struct {
	AppID       string
	Name        string
	Description string
	Type        string
	Store       string
	Parameters  string
}
type Schedule struct {
	AppID       string
	Name        string
	Description string
	Sched       string
}

func (w *Client) url() string {
	return fmt.Sprintf("http://%s:%s", w.Host, w.Port)
}
func (w *Client) GetAppList() (res *goreq.Response, err error) {
	url := w.url() + "/apps"
	res, err = goreq.Request{
		Method: http.MethodGet,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (w *Client) GetApp(id string) (res *goreq.Response, err error) {
	url := w.url() + "/apps/" + id

	res, err = goreq.Request{
		Method: http.MethodGet,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) CreateApp(name, description string) (res *goreq.Response, err error) {
	url := w.url() + "/apps"

	res, err = goreq.Request{
		Method: http.MethodPost,
		Uri:    url,
		Body:   model.App{Name: name, Description: description},
	}.Do()

	//str, err := res.Body.ToString();
	//fmt.Printf(" Resp: %s", str);

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) UpdateApp(id, name, description string) (res *goreq.Response, err error) {
	url := w.url() + "/apps/" + id
	res, err = goreq.Request{
		Method: http.MethodPost,
		Uri:    url,
		Body:   model.App{Name: name, Description: description},
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) DeleteApp(id string) (err error) {
	var req *http.Request

	url := w.url() + "/apps/" + id
	req, err = http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return
	}
	_, err = http.DefaultClient.Do(req)
	return
}

func (w *Client) CreateJob(appID string, body Job) (res *goreq.Response, err error) {
	url := w.url() + "/apps/" + appID + "/job"

	res, err = goreq.Request{
		Method: http.MethodPost,
		Uri:    url,
		Body:   body,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) GetJob(jobID string) (res *goreq.Response, err error) {
	url := w.url() + "/job/" + jobID

	res, err = goreq.Request{
		Method: http.MethodGet,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) UpdateJob(jobID string, body Job) (res *goreq.Response, err error) {
	url := w.url() + "/job/" + jobID

	res, err = goreq.Request{
		Method: http.MethodPost,
		Uri:    url,
		Body:   body,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (w *Client) DeleteJob(jobID string) (res *goreq.Response, err error) {
	url := w.url() + "/job/" + jobID

	res, err = goreq.Request{
		Method: http.MethodDelete,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (w *Client) GetJobList(appID string) (res *goreq.Response, err error) {
	url := w.url() + "/apps/" + appID + "/jobs"

	res, err = goreq.Request{
		Method: http.MethodGet,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) CreateData(appID string, body Data) (res *goreq.Response, err error) {
	url := w.url() + "/apps/" + appID + "/data"

	res, err = goreq.Request{
		Method: http.MethodPost,
		Uri:    url,
		Body:   body,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) GetData(dataID string) (res *goreq.Response, err error) {
	url := w.url() + "/data/" + dataID

	res, err = goreq.Request{
		Method: http.MethodGet,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) UpdateData(dataID string, body Data) (res *goreq.Response, err error) {
	url := w.url() + "/data/" + dataID

	res, err = goreq.Request{
		Method: http.MethodPost,
		Uri:    url,
		Body:   body,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (w *Client) DeleteData(dataID string) (res *goreq.Response, err error) {
	url := w.url() + "/data/" + dataID

	res, err = goreq.Request{
		Method: http.MethodDelete,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (w *Client) GetDataList(appID string) (res *goreq.Response, err error) {
	url := w.url() + "/apps/" + appID + "/data"

	res, err = goreq.Request{
		Method: http.MethodGet,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) CreateSchedule(appID string, body Schedule) (res *goreq.Response, err error) {
	url := w.url() + "/apps/" + appID + "/schedule"

	res, err = goreq.Request{
		Method: http.MethodPost,
		Uri:    url,
		Body:   body,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) GetSchedule(scheduleID string) (res *goreq.Response, err error) {
	url := w.url() + "/schedule/" + scheduleID

	res, err = goreq.Request{
		Method: http.MethodGet,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) UpdateSchedule(scheduleID string, body Schedule) (res *goreq.Response, err error) {
	url := w.url() + "/schedule/" + scheduleID

	res, err = goreq.Request{
		Method: http.MethodPost,
		Uri:    url,
		Body:   body,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (w *Client) DeleteSchedule(scheduleID string) (res *goreq.Response, err error) {
	url := w.url() + "/schedule/" + scheduleID

	res, err = goreq.Request{
		Method: http.MethodDelete,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (w *Client) GetScheduleList(appID string) (res *goreq.Response, err error) {
	url := w.url() + "/apps/" + appID + "/schedule"

	res, err = goreq.Request{
		Method: http.MethodGet,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Client) CloneApp(id string) (res *goreq.Response, err error) {
	url := w.url() + "/clone/" + id
	res, err = goreq.Request{
		Method: http.MethodPost,
		Uri:    url,
	}.Do()
	if err != nil {
		return nil, err
	}

	return res, nil
}
