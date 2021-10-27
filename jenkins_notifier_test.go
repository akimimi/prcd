package main

import "testing"

func TestJenkinsNotifier_NotifyUrl(t *testing.T) {
	notifier := JenkinsNotifier{
		JenkinsHost: "http://notify.website.com",
		JenkinsUrl:  "/<project>/notify?token=<token>",
		JenkinsProject: JenkinsProject{
			Name:  "pro",
			Token: "abcd1234",
		},
		UserName:     "",
		UserApiToken: "",
	}
	expected := "http://notify.website.com/pro/notify?token=abcd1234"
	if notifier.notifyUrl() != expected {
		t.Errorf("Notify url error, expected %s, actual %s", expected, notifier.notifyUrl())
	}
}

func TestJenkinsNotifier_NotifyUrl_WithProjectConfig(t *testing.T) {
	notifier := JenkinsNotifier{
		JenkinsHost: "http://notify.website.com",
		JenkinsUrl:  "/<project>/notify?token=<token>",
		JenkinsProject: JenkinsProject{
			Name:         "pro",
			Token:        "abcd1234",
			Host:         "http://project-notify.website.com",
			Url:          "/<project>/project-notify?token=<token>",
			Username:     "akimimi",
			UserApiToken: "akimimi",
		},
		UserName:     "",
		UserApiToken: "",
	}
	expected := "http://project-notify.website.com/pro/project-notify?token=abcd1234"
	if notifier.notifyUrl() != expected {
		t.Errorf("Notify url error, expected %s, actual %s", expected, notifier.notifyUrl())
	}
}

func TestJenkinsNotifier_NotifyUrl_WithInvalidProjectConfig(t *testing.T) {
	notifier := JenkinsNotifier{
		JenkinsHost: "http://notify.website.com",
		JenkinsUrl:  "/<project>/notify?token=<token>",
		JenkinsProject: JenkinsProject{
			Name:         "pro",
			Token:        "abcd1234",
			Host:         "http://project-notify.website.com",
			Url:          "/<project>/project-notify?token=<token>",
			Username:     "",
			UserApiToken: "",
		},
		UserName:     "",
		UserApiToken: "",
	}
	expected := "http://notify.website.com/pro/notify?token=abcd1234"
	if notifier.notifyUrl() != expected {
		t.Errorf("Notify url error, expected %s, actual %s", expected, notifier.notifyUrl())
	}
}

func TestJenkinsNotifier_Notify(t *testing.T) {
	notifier := JenkinsNotifier{
		JenkinsHost:    "http://notify.website.com",
		JenkinsUrl:     "/<project>/notify?token=<token>",
		JenkinsProject: JenkinsProject{},
		UserName:       "",
		UserApiToken:   "",
	}
	if err := notifier.Notify(); err == nil {
		t.Error("Notify should failed.")
	}

	notifier.JenkinsHost = "http://notify.website.com"
	notifier.JenkinsProject = JenkinsProject{
		Name:  "pro",
		Token: "abcd1234",
	}
	if err := notifier.Notify(); err == nil {
		t.Error("Notify should failed.")
	}

	notifier.JenkinsHost = ":"
	notifier.JenkinsUrl = ""
	if err := notifier.Notify(); err == nil {
		t.Error("Notify should failed.")
	}

	notifier.JenkinsHost = "https://www.mimixiche.com"
	notifier.JenkinsUrl = "/<project>/notify?token=<token>"
	if err := notifier.Notify(); err != nil {
		t.Errorf("Notify failed with %s", err)
	}
}
