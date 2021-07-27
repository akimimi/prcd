VCS Pull Request to trigger Continuous Deployment
================================================================

[![Go Reference](https://pkg.go.dev/badge/github.com/akimimi/prcd.svg)](https://pkg.go.dev/github.com/akimimi/prcd)
[![Build Status](https://travis-ci.com/akimimi/prcd.svg?branch=master)](https://travis-ci.com/akimimi/prcd)
[![Coverage Status](https://coveralls.io/repos/github/akimimi/prcd/badge.svg?branch=master)](https://coveralls.io/github/akimimi/prcd?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/akimimi/prcd)](https://goreportcard.com/report/github.com/akimimi/prcd)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This service receives pull request or push tag web hooks, and trigger 
continuous deployment(CD) in Jenkins. 

User can configure pull request to CD project map in projects.yaml. 