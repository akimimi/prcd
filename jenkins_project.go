package main

import cl "github.com/akimimi/config-loader"

// JenkinsProject defines a structure for jenkins project.
type JenkinsProject struct {
	Name         string
	Token        string
	Host         string
	Url          string
	Username     string
	UserApiToken string
}

// HasJenkinsConfig returns True if the project is configured as a dependent project
// which does not use the default Jenkins Host and Url configuration.
func (p *JenkinsProject) HasJenkinsConfig() bool {
	return p.Host != "" && p.Url != "" && p.Username != "" && p.UserApiToken != ""
}

// JenkinsProjectConfig defines the structure for jenkins configure.
type JenkinsProjectConfig struct {
	Environment    string `json:"environment" yaml:"environment"`
	VcsProject     string `json:"vcs_project" yaml:"vcs_project"`
	Branch         string `json:"branch" yaml:"branch"`
	JenkinsProject string `json:"jenkins_project" yaml:"jenkins_project"`
	JenkinsToken   string `json:"jenkins_token" yaml:"jenkins_token"`

	// The following parameters are not forced, default values will be used if one of the following parameter is empty.
	JenkinsHost         string `json:"jenkins_host" yaml:"jenkins_host"`
	JenkinsUrl          string `json:"jenkins_url" yaml:"jenkins_url"`
	JenkinsUsername     string `json:"jenkins_username" yaml:"jenkins_username"`
	JenkinsUserApiToken string `json:"jenkins_user_api_token" yaml:"jenkins_user_api_token"`
}

var jenkinsProjectConfigGrp map[string]JenkinsProjectConfig

func loadJenkinsProjectConfig(filename string) {
	cl.LoadByFile(filename, &jenkinsProjectConfigGrp)
}

func matchJenkinsProject(environment, project, branch string) JenkinsProject {
	for _, config := range jenkinsProjectConfigGrp {
		if config.Environment == environment && config.VcsProject == project && config.Branch == branch {
			return JenkinsProject{
				Name:         config.JenkinsProject,
				Token:        config.JenkinsToken,
				Host:         config.JenkinsHost,
				Url:          config.JenkinsUrl,
				Username:     config.JenkinsUsername,
				UserApiToken: config.JenkinsUserApiToken,
			}
		}
	}
	return JenkinsProject{}
}
