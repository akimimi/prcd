package main

import (
	"github.com/gogap/errors"
	"github.com/gogap/logs"
	"io"
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
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取一小段响应体用于诊断（Jenkins 触发成功一般是 201 Created + Location: /queue/item/...）。
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	bodySnippet := strings.TrimSpace(string(bodyBytes))
	location := resp.Header.Get("Location")

	// Jenkins 触发构建一般返回 201 Created（带 Location 指向 queue item），
	// 老的判定只接受 200 OK，会把 201 当成失败、把任意 200 页面当成成功，这里改为接受所有 2xx。
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logs.Info("Notified to project ", notifier.JenkinsProject.Name,
			" status=", resp.Status, " location=", location, " body=", bodySnippet)
		return nil
	}
	return errors.New("Notify failed: project=" + notifier.JenkinsProject.Name +
		" status=" + resp.Status + " body=" + bodySnippet)
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
