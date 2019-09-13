package resource_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/google/go-github/v28/github"

	resource "github.com/ahume/github-deployment-resource"
	"github.com/ahume/github-deployment-resource/fakes"
)

var _ = Describe("In Command", func() {
	var (
		command      *resource.InCommand
		githubClient *fakes.FakeGitHub

		inRequest resource.InRequest

		inResponse resource.InResponse
		inErr      error

		tmpDir  string
		destDir string
	)

	BeforeEach(func() {
		var err error

		githubClient = &fakes.FakeGitHub{}
		command = resource.NewInCommand(githubClient, ioutil.Discard)

		tmpDir, err = ioutil.TempDir("", "github-deployment")
		Ω(err).ShouldNot(HaveOccurred())

		destDir = filepath.Join(tmpDir, "destination")

		inRequest = resource.InRequest{}
	})

	AfterEach(func() {
		Ω(os.RemoveAll(tmpDir)).Should(Succeed())
	})

	buildDeployment := func(ID int64, env string, task string) *github.Deployment {
		return &github.Deployment{
			ID:          github.Int64(ID),
			Environment: github.String(env),
			Task:        github.String(task),
			Ref:         github.String("master"),
			SHA:         github.String("12345"),
			Description: github.String("One more"),
			Creator: &github.User{
				Login: github.String("Something"),
			},
			CreatedAt: &github.Timestamp{time.Date(2016, 01, 20, 15, 15, 15, 0, time.UTC)},
		}
	}

	buildDeploymentStatus := func(ID int64, state string) *github.DeploymentStatus {
		return &github.DeploymentStatus{
			ID:        github.Int64(ID),
			State:     github.String(state),
			CreatedAt: &github.Timestamp{time.Date(2016, 01, 20, 15, 15, 15, 0, time.UTC)},
		}
	}

	Context("when there is a deployment found", func() {
		disaster := errors.New("no deployment")

		BeforeEach(func() {
			githubClient.GetDeploymentReturns(&github.Deployment{}, disaster)

			inRequest.Version = resource.Version{
				ID: "1",
			}
		})

		It("returns an appropriate error", func() {
			inResponse, inErr = command.Run(destDir, inRequest)

			Expect(inErr).To(Equal(disaster))
		})
	})

	Context("when there is a deployment found", func() {
		BeforeEach(func() {
			githubClient.GetDeploymentReturns(buildDeployment(1, "production", "deploy"), nil)
			githubClient.ListDeploymentStatusesReturns([]*github.DeploymentStatus{
				buildDeploymentStatus(1, "success"),
			}, nil)

			inRequest.Version = resource.Version{
				ID: "1",
			}
		})

		It("creates the correct data files", func() {
			inResponse, inErr = command.Run(destDir, inRequest)

			contents, err := ioutil.ReadFile(path.Join(destDir, "id"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(contents)).Should(Equal("1"))

			contents, err = ioutil.ReadFile(path.Join(destDir, "ref"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(contents)).Should(Equal("master"))

			contents, err = ioutil.ReadFile(path.Join(destDir, "sha"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(contents)).Should(Equal("12345"))

			contents, err = ioutil.ReadFile(path.Join(destDir, "task"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(contents)).Should(Equal("deploy"))

			contents, err = ioutil.ReadFile(path.Join(destDir, "environment"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(contents)).Should(Equal("production"))

			contents, err = ioutil.ReadFile(path.Join(destDir, "description"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(string(contents)).Should(Equal("One more"))
		})

		It("outputs the correct metadata", func() {
			inResponse, inErr = command.Run(destDir, inRequest)

			Ω(inResponse.Metadata).Should(ConsistOf(
				resource.MetadataPair{Name: "id", Value: "1"},
				resource.MetadataPair{Name: "ref", Value: "master"},
				resource.MetadataPair{Name: "sha", Value: "12345"},
				resource.MetadataPair{Name: "task", Value: "deploy"},
				resource.MetadataPair{Name: "description", Value: "One more"},
				resource.MetadataPair{Name: "environment", Value: "production"},
				resource.MetadataPair{Name: "creator", Value: "Something"},
				resource.MetadataPair{Name: "created_at", Value: "2016-01-20 15:15:15"},
				resource.MetadataPair{Name: "status_id", Value: "1"},
				resource.MetadataPair{Name: "status", Value: "success"},
				resource.MetadataPair{Name: "status_created_at", Value: "2016-01-20 15:15:15"},
				resource.MetadataPair{Name: "status_count", Value: "1"},
			))
		})
	})
})
