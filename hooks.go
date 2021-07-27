package main

// Project is the struct for a repository in VCS
type Project struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

// Branch is the struct for a branch data in VCS
type Branch struct {
	Label string  `json:"label"`
	Ref   string  `json:"ref"`
	Sha   string  `json:"sha"`
	Repo  Project `json:"repo"`
}

// PullRequest is the struct for a pull request record in VCS
type PullRequest struct {
	Id        int    `json:"id"`
	State     string `json:"state"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Base      Branch `json:"base"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// BasicHook contains the common parameters for a VCS webhook.
type BasicHook struct {
	HookName string `json:"hook_name"`
	HookId   int    `json:"hook_id,omitempty"`
	HookUrl  string `json:"hook_url,omitempty"`
}

// PullRequestHook is the pull request webhook struct.
type PullRequestHook struct {
	BasicHook   `json:",inline"`
	PullRequest PullRequest `json:"pull_request"`
}

// PushTagHook is the push and tag webhook struct.
type PushTagHook struct {
	BasicHook `json:",inline"`
	Ref       string  `json:"ref"`
	Project   Project `json:"repository"`
}
