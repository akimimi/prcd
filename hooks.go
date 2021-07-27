package main

type Project struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

type Branch struct {
	Label string  `json:"label"`
	Ref   string  `json:"ref"`
	Sha   string  `json:"sha"`
	Repo  Project `json:"repo"`
}

type PullRequest struct {
	Id        int    `json:"id"`
	State     string `json:"state"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Base      Branch `json:"base"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BasicHook struct {
	HookName string `json:"hook_name"`
	HookId   int    `json:"hook_id,omitempty"`
	HookUrl  string `json:"hook_url,omitempty"`
}

type PullRequestHook struct {
	BasicHook   `json:",inline"`
	PullRequest PullRequest `json:"pull_request"`
}

type PushTagHook struct {
	BasicHook `json:",inline"`
	Ref       string  `json:"ref"`
	Project   Project `json:"repository"`
}
