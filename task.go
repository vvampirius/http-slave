package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

var (
	ErrEmptyUrl = errors.New(`URL is empty`)
	ErrBadStatus = errors.New(`Status is not 200`)
	ErrNotJson = errors.New(`Is not JSON`)
)

type Task struct {
	Command []string		`json:"command"`
	RespondUrl string		`json:"respond_url"`
	ImmediatelyNext bool	`json:"immediately_next"`
}

func GetTask(url string) (*Task, error) {
	if url == `` {
		ErrorLog.Println(ErrEmptyUrl)
		return nil, ErrEmptyUrl
	}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ErrorLog.Println(err.Error())
		return nil, err
	}
	client := http.Client{ Timeout: time.Second * 3	}
	response, err := client.Do(request)
	if err != nil {
		ErrorLog.Println(err.Error())
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		ErrorLog.Println(response.Status)
		return nil, ErrBadStatus
	}
	if contentType := response.Header.Get(`Content-Type`); contentType != `application/json` {
		ErrorLog.Println(contentType)
		return nil, ErrNotJson
	}
	task := Task{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&task); err != nil {
		ErrorLog.Println(err.Error())
		return nil, err
	}
	return &task, nil
}

func ExecTask(task *Task) ([]byte, int, time.Duration, error) {
	start := time.Now()
	cmd := exec.Command(task.Command[0], task.Command[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		ErrorLog.Println(err.Error())
		return nil, 0, 0, err
	}
	spent := time.Now().Sub(start)
	return stdoutStderr, cmd.ProcessState.ExitCode(), spent, nil
}

func RespondTask(respondUrl string, data []byte, exitCode int, spent time.Duration) {
	buffer := bytes.NewBuffer(data)
	request, err := http.NewRequest(http.MethodPost, respondUrl, buffer)
	if err != nil {
		ErrorLog.Println(err.Error())
		return
	}
	request.Header.Add(`EXIT_CODE`, fmt.Sprintf(`%d`, exitCode))
	request.Header.Add(`SPENT`, spent.String())
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ErrorLog.Println(err.Error())
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		ErrorLog.Println(respondUrl, response.Status)
	}
}