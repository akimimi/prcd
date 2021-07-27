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

	notifier.JenkinsHost = "http://www.mimixiche.com"
	notifier.JenkinsUrl = "/<project>/notify?token=<token>"
	if err := notifier.Notify(); err != nil {
		t.Errorf("Notify failed with %s", err)
	}
}
