package main

import (
	"encoding/json"
	"github.com/gogap/logs"
	"strings"
)

// HookAgent defines a webhook agent interface. Struct of a webhook agent should satisfies the following interface.
type HookAgent interface {
	Name() string
	Parse(b []byte) error
	CanTriggerEvent() bool
	HookBranch() string
	HookProject() string
	Environment() string
}

// PullRequestHookAgent is the agent for pull request transfer.
type PullRequestHookAgent struct {
	prHook   PullRequestHook
	isParsed bool
}

// Name is the agent name implementation.
func (agent *PullRequestHookAgent) Name() string {
	return "PullRequestHookAgent"
}

// Parse unmarshal given bytes to agent.
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

// CanTriggerEvent determines whether an agent can trigger following events.
// The pull request webhook can trigger CD events only if the state is "merged".
func (agent *PullRequestHookAgent) CanTriggerEvent() bool {
	return agent.prHook.PullRequest.State == "merged"
}

// HookBranch returns the branch name of a pull request.
func (agent *PullRequestHookAgent) HookBranch() string {
	if !agent.isParsed {
		return ""
	}
	return agent.prHook.PullRequest.Base.Ref
}

// HookProject returns the project name of a pull request.
func (agent *PullRequestHookAgent) HookProject() string {
	if !agent.isParsed {
		return ""
	}
	return agent.prHook.PullRequest.Base.Repo.Name
}

// Environment returns "debug" or "production" based on the branch and project.
func (agent *PullRequestHookAgent) Environment() string {
	return hookBranchEnvironment(agent.HookBranch())
}

// PushTagHookAgent is the agent for pull request transfer.
type PushTagHookAgent struct {
	pushHook PushTagHook
	isParsed bool
}

// Name is the agent name implementation.
func (agent *PushTagHookAgent) Name() string {
	return "PushTagHookAgent"
}

// Parse unmarshal given bytes to agent.
func (agent *PushTagHookAgent) Parse(b []byte) error {
	var e error
	agent.isParsed = false
	if e = json.Unmarshal(b, &agent.pushHook); e == nil {
		agent.isParsed = true
		logs.Debug("TagPush:", agent.pushHook.Ref, "/", agent.pushHook.Project.Name, "/", agent.pushHook.Project.FullName)
	}
	return e
}

// CanTriggerEvent determines whether an agent can trigger following events.
// The push tag hook cannot trigger any events.
func (agent *PushTagHookAgent) CanTriggerEvent() bool {
	return false
}

// HookBranch returns the branch name of a push or the tag name.
func (agent *PushTagHookAgent) HookBranch() string {
	if !agent.isParsed {
		return ""
	}
	return agent.pushHook.Ref
}

// HookProject returns the project name of a push or tag.
func (agent *PushTagHookAgent) HookProject() string {
	if !agent.isParsed {
		return ""
	}
	return agent.pushHook.Project.Name
}

// Environment returns "debug" or "production" based on the branch and project.
func (agent *PushTagHookAgent) Environment() string {
	return hookBranchEnvironment(agent.HookBranch())
}

// DefaultHookAgent is a fake agent struct, a default agent cannot trigger any following events.
type DefaultHookAgent struct {
	isParsed bool
}

// Name is the agent name implementation.
func (agent *DefaultHookAgent) Name() string {
	return "DefaultHookAgent"
}

// Parse unmarshal given bytes to agent.
func (agent *DefaultHookAgent) Parse(_ []byte) error {
	return nil
}

// CanTriggerEvent determines whether an agent can trigger following events.
// The default agent cannot trigger any events.
func (agent *DefaultHookAgent) CanTriggerEvent() bool {
	return false
}

// HookBranch always returns "unknown" for a default agent.
func (agent *DefaultHookAgent) HookBranch() string {
	return "unknown"
}

// HookProject always returns "unknown" for a default agent.
func (agent *DefaultHookAgent) HookProject() string {
	return "unknown"
}

// Environment always returns "unknown" for a default agent.
func (agent *DefaultHookAgent) Environment() string {
	return "unknown"
}

func hookBranchEnvironment(branch string) string {
	if strings.HasPrefix(branch, "master") || strings.HasPrefix(branch, "release") {
		return "production"
	}
	return "debug"
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
