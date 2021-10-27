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
	username, apiToken := notifier.UserName, notifier.UserApiToken
	if notifier.JenkinsProject.HasJenkinsConfig() {
		username, apiToken = notifier.JenkinsProject.Username, notifier.JenkinsProject.UserApiToken
	}
	req.SetBasicAuth(username, apiToken)
	resp, err := http.DefaultClient.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		logs.Info("Notified to project ", notifier.JenkinsProject.Name)
		return nil
	} else {
		if err == nil {
			return errors.New("Notify Status is " + resp.Status)
		} else {
			return err
		}
	}
}

func (notifier *JenkinsNotifier) notifyUrl() string {
	host := notifier.JenkinsHost
	url := notifier.JenkinsUrl
	if notifier.JenkinsProject.HasJenkinsConfig() {
		host = notifier.JenkinsProject.Host
		url = notifier.JenkinsProject.Url
	}
	url = strings.Replace(url, "<project>", notifier.JenkinsProject.Name, 1)
	url = strings.Replace(url, "<token>", notifier.JenkinsProject.Token, 1)
	return host + url
}
