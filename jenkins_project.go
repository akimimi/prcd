package main

import config_loader "github.com/akimimi/config-loader"

// JenkinsProject defines a structure for jenkins project.
type JenkinsProject struct {
	Name  string
	Token string
}

// JenkinsProjectConfig defines the structure for jenkins configure.
type JenkinsProjectConfig struct {
	Environment    string `json:"environment" yaml:"environment"`
	VcsProject     string `json:"vcs_project" yaml:"vcs_project"`
	Branch         string `json:"branch" yaml:"branch"`
	JenkinsProject string `json:"jenkins_project" yaml:"jenkins_project"`
	JenkinsToken   string `json:"jenkins_token" yaml:"jenkins_token"`
}

var jenkinsProjectConfigGrp map[string]JenkinsProjectConfig

func loadJenkinsProjectConfig(filename string) {
	config_loader.LoadByFile(filename, &jenkinsProjectConfigGrp)
}

func matchJenkinsProject(environment, project, branch string) JenkinsProject {
	for _, config := range jenkinsProjectConfigGrp {
		if config.Environment == environment && config.VcsProject == project && config.Branch == branch {
			return JenkinsProject{
				Name:  config.JenkinsProject,
				Token: config.JenkinsToken,
			}
		}
	}
	return JenkinsProject{}
}
