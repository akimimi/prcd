package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestProjectJsonUnmarshal(t *testing.T) {
	proj := Project{}
	filename := "samples/project.json"
	if file, e := ioutil.ReadFile(filename); e != nil {
		panic(e)
	} else {
		json.Unmarshal(file, &proj)
	}
	if proj.Name != "mingdao" {
		t.Errorf("Project unmarshal failed, expected %s, actual %s", "mingdao", proj.Name)
	}
}

func TestPullRequestHookJsonUnmarshal(t *testing.T) {
	prHook := PullRequestHook{}
	filename := "samples/pull_request.json"
	if file, e := ioutil.ReadFile(filename); e != nil {
		panic(e)
	} else {
		json.Unmarshal(file, &prHook)
	}
	if prHook.HookName != "merge_request_hooks" {
		t.Errorf("Pull request hook unmarshal failed, expected %s, actual %s", "merge_request_hooks", prHook.HookName)
	}
}

func TestHookBranchEnvironment(t *testing.T) {
	testData := map[string]string{
		"master":              "production",
		"master7":             "production",
		"release":             "production",
		"release/version-3.0": "production",
		"develop":             "debug",
		"develop7":            "debug",
	}
	for branch, expected := range testData {
		if hookBranchEnvironment(branch) != expected {
			t.Errorf("Environment failed for %s, expected %s, actual %s",
				branch, expected, hookBranchEnvironment(branch))
		}
	}
}

func TestPullRequestHookAgent_Parse(t *testing.T) {
	agent := PullRequestHookAgent{}
	filename := "samples/pull_request.json"
	if file, e := ioutil.ReadFile(filename); e != nil {
		panic(e)
	} else {
		agent.Parse(file)
	}
	if !agent.isParsed {
		t.Error("Pull request parse failed!")
	}
}

func TestPullRequestHookAgent_Name(t *testing.T) {
	agent := PullRequestHookAgent{}
	if agent.Name() != "PullRequestHookAgent" {
		t.Error("PullRequestHookAgent name is not correct")
	}
}

func TestPullRequestHookAgent_HookProject_HookBranch(t *testing.T) {
	filename := "samples/pull_request.json"
	agent := PullRequestHookAgent{}
	if file, e := ioutil.ReadFile(filename); e == nil {
		agent.Parse(file)
	}
	if agent.HookProject() != "mingdao" {
		t.Errorf("Pull Request project is not correct, expected %s, actual %s", "mingdao", agent.HookProject())
	}
	if agent.HookBranch() != "master" {
		t.Errorf("Pull Request branch is not correct, expected %s, actual %s", "master", agent.HookBranch())
	}
}

func TestPullRequestHookAgent_HookProject_HookBranch_Failed(t *testing.T) {
	filename := "samples/pull_request.json"
	agent := PullRequestHookAgent{}
	if file, e := ioutil.ReadFile(filename); e == nil {
		agent.Parse(file)
	}
	agent.isParsed = false
	if agent.HookProject() != "" {
		t.Errorf("Pull Request project is not correct, expected %s, actual %s", "", agent.HookProject())
	}
	if agent.HookBranch() != "" {
		t.Errorf("Pull Request branch is not correct, expected %s, actual %s", "", agent.HookBranch())
	}
}

func TestPullRequestHookAgent_CanTriggerEvent(t *testing.T) {
	filename := "samples/pull_request.json"
	agent := PullRequestHookAgent{}
	if file, e := ioutil.ReadFile(filename); e == nil {
		agent.Parse(file)
	}
	agent.prHook.PullRequest.State = "open"
	if agent.CanTriggerEvent() {
		t.Errorf("Pull Request in %s state should not trigger events.", agent.prHook.PullRequest.State)
	}
	agent.prHook.PullRequest.State = "closed"
	if agent.CanTriggerEvent() {
		t.Errorf("Pull Request in %s state should not trigger events.", agent.prHook.PullRequest.State)
	}
	agent.prHook.PullRequest.State = "merged"
	if !agent.CanTriggerEvent() {
		t.Errorf("Pull Request in %s state should trigger events.", agent.prHook.PullRequest.State)
	}
}

func TestPullRequestHookAgent_Environment(t *testing.T) {
	filename := "samples/pull_request.json"
	agent := PullRequestHookAgent{}
	if file, e := ioutil.ReadFile(filename); e == nil {
		agent.Parse(file)
	}
	testData := map[string]string{
		"master":              "production",
		"master7":             "production",
		"release":             "production",
		"release/version-3.0": "production",
		"develop":             "debug",
		"develop7":            "debug",
	}
	for branch, expected := range testData {
		agent.prHook.PullRequest.Base.Ref = branch
		if agent.Environment() != expected {
			t.Errorf("Environment failed for %s, expected %s, actual %s",
				branch, expected, agent.Environment())
		}
	}
}

func TestPushTagHookAgentJsonUnmarshal(t *testing.T) {
	hook := PushTagHook{}
	filename := "samples/push_tag.json"
	if file, e := ioutil.ReadFile(filename); e != nil {
		panic(e)
	} else {
		json.Unmarshal(file, &hook)
	}
	if hook.HookName != "push_hooks" {
		t.Errorf("Pull request hook unmarshal failed, expected %s, actual %s", "push_hooks", hook.HookName)
	}
}

func TestPushTagHookAgent_Parse(t *testing.T) {
	agent := PushTagHookAgent{}
	filename := "samples/push_tag.json"
	if file, e := ioutil.ReadFile(filename); e != nil {
		panic(e)
	} else {
		agent.Parse(file)
	}
	if !agent.isParsed || agent.pushHook.HookName != "push_hooks" {
		t.Error("Push tag parse failed!")
	}
}

func TestPushTagHookAgent_Name(t *testing.T) {
	agent := PushTagHookAgent{}
	if agent.Name() != "PushTagHookAgent" {
		t.Error("PushTagHookAgent name is not correct")
	}
}

func TestPushTagHookAgent_HookProject_HookBranch(t *testing.T) {
	filename := "samples/push_tag.json"
	agent := PushTagHookAgent{}
	if file, e := ioutil.ReadFile(filename); e == nil {
		agent.Parse(file)
	}
	if agent.HookProject() != "Gitee" {
		t.Errorf("Pull Request project is not correct, expected %s, actual %s", "Gitee", agent.HookProject())
	}
	if agent.HookBranch() != "refs/heads/change_commitlint_config" {
		t.Errorf("Pull Request branch is not correct, expected %s, actual %s", "refs/heads/change_commitlint_config", agent.HookBranch())
	}
}

func TestPushTagHookAgent_HookProject_HookBranch_Failed(t *testing.T) {
	filename := "samples/push_tag.json"
	agent := PushTagHookAgent{}
	if file, e := ioutil.ReadFile(filename); e == nil {
		agent.Parse(file)
	}
	agent.isParsed = false
	if agent.HookProject() != "" {
		t.Errorf("Push tag project is not correct, expected %s, actual %s", "", agent.HookProject())
	}
	if agent.HookBranch() != "" {
		t.Errorf("Push tag branch is not correct, expected %s, actual %s", "", agent.HookBranch())
	}
}

func TestPushTagHookAgent_CanTriggerEvent(t *testing.T) {
	filename := "samples/push_tag.json"
	agent := PushTagHookAgent{}
	if file, e := ioutil.ReadFile(filename); e == nil {
		agent.Parse(file)
	}
	if agent.CanTriggerEvent() {
		t.Error("Push tag should not trigger events.")
	}
}

func TestPushTagHookAgent_Environment(t *testing.T) {
	filename := "samples/push_tag.json"
	agent := PushTagHookAgent{}
	if file, e := ioutil.ReadFile(filename); e == nil {
		agent.Parse(file)
	}
	testData := map[string]string{
		"master":              "production",
		"master7":             "production",
		"release":             "production",
		"release/version-3.0": "production",
		"develop":             "debug",
		"develop7":            "debug",
	}
	for branch, expected := range testData {
		agent.pushHook.Ref = branch
		if agent.Environment() != expected {
			t.Errorf("Environment failed for %s, expected %s, actual %s",
				branch, expected, agent.Environment())
		}
	}
}

func TestDefaultHookAgent_Parse(t *testing.T) {
	agent := DefaultHookAgent{}
	filename := "samples/push_tag.json"
	if file, e := ioutil.ReadFile(filename); e != nil {
		panic(e)
	} else {
		agent.Parse(file)
	}
	if agent.isParsed {
		t.Error("Default hook agent should not be parsed.")
	}
}

func TestDefaultHookAgent_Name(t *testing.T) {
	agent := DefaultHookAgent{}
	if agent.Name() != "DefaultHookAgent" {
		t.Error("DefaultHookAgent name is not correct")
	}
}

func TestDefaultHookAgent_HookProject_HookBranch(t *testing.T) {
	agent := DefaultHookAgent{}
	if agent.HookProject() != "unknown" {
		t.Errorf("Hook project is not correct, expected %s, actual %s", "unknown", agent.HookProject())
	}
	if agent.HookBranch() != "unknown" {
		t.Errorf("Hook branch is not correct, expected %s, actual %s", "unknown", agent.HookBranch())
	}
}

func TestDefaultHookAgent_CanTriggerEvent(t *testing.T) {
	agent := DefaultHookAgent{}
	if agent.CanTriggerEvent() {
		t.Error("Default hook should not trigger events.")
	}
}

func TestDefaultHookAgent_Environment(t *testing.T) {
	agent := DefaultHookAgent{}
	expected := "unknown"
	if agent.Environment() != expected {
		t.Errorf("Environment failed, expected %s, actual %s",
			expected, agent.Environment())
	}
}

func TestCreateHookAgentByName(t *testing.T) {
	agent := createHookAgentByName("merge_request_hooks")
	if agent.Name() != "PullRequestHookAgent" {
		t.Errorf("Agent created failed, expected %s, actual %s", "PullRequestHookAgent", agent.Name())
	}

	agent = createHookAgentByName("push_hooks")
	if agent.Name() != "PushTagHookAgent" {
		t.Errorf("Agent created failed, expected %s, actual %s", "PushTagHookAgent", agent.Name())
	}

	agent = createHookAgentByName("tag_push_hooks")
	if agent.Name() != "PushTagHookAgent" {
		t.Errorf("Agent created failed, expected %s, actual %s", "PushTagHookAgent", agent.Name())
	}

	agent = createHookAgentByName("")
	if agent.Name() != "DefaultHookAgent" {
		t.Errorf("Agent created failed, expected %s, actual %s", "DefaultHookAgent", agent.Name())
	}
}

func TestCreateNotifierByAgent(t *testing.T) {
	loadJenkinsProjectConfig("config/projects.sample.yaml")
	agent := PullRequestHookAgent{
		isParsed: true,
		prHook: PullRequestHook{
			PullRequest: PullRequest{
				Base: Branch{
					Ref: "develop",
					Repo: Project{
						Name: "mimixiche-backend",
					},
				},
			},
		},
	}
	notifier := createNotifierByAgent(&agent)
	if notifier.JenkinsProject.Name == "" {
		t.Error("Create notifier failed!")
	}
}
