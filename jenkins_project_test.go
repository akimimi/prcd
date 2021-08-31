package main

import "testing"

func TestJenkinsProjectConfigParsing(t *testing.T) {
	loadJenkinsProjectConfig("config/projects.sample.yaml")
	if _, ok := jenkinsProjectConfigGrp["dev-backend"]; !ok {
		t.Error("Jenkins project config dev-backend parse error.")
		t.Fail()
	}
	if _, ok := jenkinsProjectConfigGrp["dev-backend-php7"]; !ok {
		t.Error("Jenkins project config dev-backend-php7 parse error.")
		t.Fail()
	}
	if jenkinsProjectConfigGrp["dev-backend"].Environment != "debug" ||
		jenkinsProjectConfigGrp["dev-backend"].VcsProject != "mimixiche-backend" ||
		jenkinsProjectConfigGrp["dev-backend"].Branch != "develop" ||
		jenkinsProjectConfigGrp["dev-backend"].JenkinsProject != "dev-jenkins-project" ||
		jenkinsProjectConfigGrp["dev-backend"].JenkinsToken != "abcdefg1234" {
		t.Error("Jenkins config dev-backend parse failed!")
	}
	if jenkinsProjectConfigGrp["dev-backend-php7"].Environment != "debug" ||
		jenkinsProjectConfigGrp["dev-backend-php7"].VcsProject != "mimixiche-backend" ||
		jenkinsProjectConfigGrp["dev-backend-php7"].Branch != "develop7" ||
		jenkinsProjectConfigGrp["dev-backend-php7"].JenkinsProject != "dev-jenkins-project-php7" ||
		jenkinsProjectConfigGrp["dev-backend-php7"].JenkinsToken != "abcdefg1234" {
		t.Error("Jenkins config dev-backend parse failed!")
	}
	if jenkinsProjectConfigGrp["release-backend"].Environment != "production" ||
		jenkinsProjectConfigGrp["release-backend"].VcsProject != "mimixiche-backend" ||
		jenkinsProjectConfigGrp["release-backend"].Branch != "release" ||
		jenkinsProjectConfigGrp["release-backend"].JenkinsProject != "production-backend-release" ||
		jenkinsProjectConfigGrp["release-backend"].JenkinsToken != "abcdefg1234" {
		t.Error("Jenkins config dev-backend parse failed!")
	}
}

func TestMatchJenkinsProject(t *testing.T) {
	loadJenkinsProjectConfig("config/projects.sample.yaml")
	env, project, branch := "debug", "mimixiche-backend", "develop"
	jenkinsProject := matchJenkinsProject(env, project, branch)
	expected := "dev-jenkins-project"
	if jenkinsProject.Name != expected {
		t.Errorf("Jenkins project error, expected %s, actual %s.", expected, jenkinsProject.Name)
	}

	env, project, branch = "debug", "mimixiche-backend", "develop7"
	jenkinsProject = matchJenkinsProject(env, project, branch)
	expected = "dev-jenkins-project-php7"
	if jenkinsProject.Name != expected {
		t.Errorf("Jenkins project error, expected %s, actual %s.", expected, jenkinsProject.Name)
	}

	env, project, branch = "production", "mimixiche-backend", "release"
	jenkinsProject = matchJenkinsProject(env, project, branch)
	expected = "production-backend-release"
	if jenkinsProject.Name != expected {
		t.Errorf("Jenkins project error, expected %s, actual %s.", expected, jenkinsProject.Name)
	}

	env, project, branch = "master", "mimixiche-backend", "develop"
	jenkinsProject = matchJenkinsProject(env, project, branch)
	expected = ""
	if jenkinsProject.Name != expected {
		t.Errorf("Jenkins project error, expected %s, actual %s.", expected, jenkinsProject.Name)
	}
}
