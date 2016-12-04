package resource_test

import (
	"github.com/ahume/go-github/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGithubDeploymentResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GithubDeploymentResource Suite")
}

func newDeployment(id int) *github.Deployment {
	return &github.Deployment{
		ID: github.Int(id),
	}
}
