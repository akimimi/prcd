package main

import (
	"github.com/gogap/errors"
	"github.com/gogap/logs"
	"net/http"
	"strings"
)

// JenkinsNotifier defines a notify struct which contains CD host, url, project and user information.
type JenkinsNotifier struct {
	JenkinsHost    string
	JenkinsUrl     string
	JenkinsProject JenkinsProject
	UserName       string
	UserApiToken   string
}

// Notify executes notify based on CD information in the struct.
func (notifier *JenkinsNotifier) Notify() error {
	if notifier.JenkinsProject.Name == "" || notifier.JenkinsProject.Token == "" {
		return errors.New("Jenkins Project config is not correct.")
	}

	url := notifier.notifyUrl()
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(notifier.UserName, notifier.UserApiToken)
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		defer resp.Body.Close()
		logs.Info("Notified to project ", notifier.JenkinsProject.Name)
		return nil
	} else {
		return err
	}
}

func (notifier *JenkinsNotifier) notifyUrl() string {
	url := strings.Replace(notifier.JenkinsUrl, "<project>", notifier.JenkinsProject.Name, 1)
	url = strings.Replace(url, "<token>", notifier.JenkinsProject.Token, 1)
	return notifier.JenkinsHost + url
}
