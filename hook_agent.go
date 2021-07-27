package main

import (
	"encoding/json"
	"github.com/gogap/logs"
	"strings"
)

type HookAgent interface {
	Name() string
	Parse(b []byte) error
	CanTriggerEvent() bool
	HookBranch() string
	HookProject() string
	Environment() string
}

type PullRequestHookAgent struct {
	prHook   PullRequestHook
	isParsed bool
}

func (agent *PullRequestHookAgent) Name() string {
	return "PullRequestHookAgent"
}

func (agent *PullRequestHookAgent) Parse(b []byte) error {
	var e error
	agent.isParsed = false
	if e = json.Unmarshal(b, &agent.prHook); e == nil {
		agent.isParsed = true
		logs.Debug("PR:", agent.prHook.PullRequest.Title, "/", agent.prHook.PullRequest.Base.Repo.Name,
			"/", agent.prHook.PullRequest.Base.Ref, "/", agent.prHook.PullRequest.State)
	}
	return e
}

func (agent *PullRequestHookAgent) CanTriggerEvent() bool {
	return agent.prHook.PullRequest.State == "merged"
}

func (agent *PullRequestHookAgent) HookBranch() string {
	if !agent.isParsed {
		return ""
	}
	return agent.prHook.PullRequest.Base.Ref
}

func (agent *PullRequestHookAgent) HookProject() string {
	if !agent.isParsed {
		return ""
	}
	return agent.prHook.PullRequest.Base.Repo.Name
}

func (agent *PullRequestHookAgent) Environment() string {
	return hookBranchEnvironment(agent.HookBranch())
}

type PushTagHookAgent struct {
	pushHook PushTagHook
	isParsed bool
}

func (agent *PushTagHookAgent) Name() string {
	return "PushTagHookAgent"
}

func (agent *PushTagHookAgent) Parse(b []byte) error {
	var e error
	agent.isParsed = false
	if e = json.Unmarshal(b, &agent.pushHook); e == nil {
		agent.isParsed = true
		logs.Debug("TagPush:", agent.pushHook.Ref, "/", agent.pushHook.Project.Name, "/", agent.pushHook.Project.FullName)
	}
	return e
}

func (agent *PushTagHookAgent) CanTriggerEvent() bool {
	return false
}

func (agent *PushTagHookAgent) HookBranch() string {
	if !agent.isParsed {
		return ""
	}
	return agent.pushHook.Ref
}

func (agent *PushTagHookAgent) HookProject() string {
	if !agent.isParsed {
		return ""
	}
	return agent.pushHook.Project.Name
}

func (agent *PushTagHookAgent) Environment() string {
	return hookBranchEnvironment(agent.HookBranch())
}

type DefaultHookAgent struct {
	isParsed bool
}

func (agent *DefaultHookAgent) Name() string {
	return "DefaultHookAgent"
}

func (agent *DefaultHookAgent) Parse(b []byte) error {
	return nil
}

func (agent *DefaultHookAgent) CanTriggerEvent() bool {
	return false
}

func (agent *DefaultHookAgent) HookBranch() string {
	return "unknown"
}

func (agent *DefaultHookAgent) HookProject() string {
	return "unknown"
}

func (agent *DefaultHookAgent) Environment() string {
	return "unknown"
}

func hookBranchEnvironment(branch string) string {
	if strings.HasPrefix(branch, "master") || strings.HasPrefix(branch, "release") {
		return "production"
	} else {
		return "debug"
	}
}

func createHookAgentByName(name string) HookAgent {
	if name == "merge_request_hooks" {
		return &PullRequestHookAgent{}
	}
	if name == "tag_push_hooks" || name == "push_hooks" {
		return &PushTagHookAgent{}
	}
	return &DefaultHookAgent{}
}

func createNotifierByAgent(agent HookAgent) *JenkinsNotifier {
	project := matchJenkinsProject(agent.Environment(), agent.HookProject(), agent.HookBranch())
	notifier := JenkinsNotifier{
		JenkinsHost:    settings.jenkinsHost,
		JenkinsUrl:     settings.jenkinsNotifyUrl,
		JenkinsProject: project,
		UserName:       settings.jenkinsUserName,
		UserApiToken:   settings.jenkinsUserApiToken,
	}
	return &notifier
}
