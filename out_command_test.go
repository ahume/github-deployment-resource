package resource_test

import (
  "io/ioutil"
  "os"
  "path/filepath"
  "time"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "github.com/ahume/go-github/github"

  "github.com/ahume/github-deployment-resource"
  "github.com/ahume/github-deployment-resource/fakes"
)

func file(path, contents string) {
  Ω(ioutil.WriteFile(path, []byte(contents), 0644)).Should(Succeed())
}

var _ = Describe("Status Out Command", func() {
  var (
    command      *resource.OutCommand
    githubClient *fakes.FakeGitHub

    sourcesDir string

    request resource.OutRequest
  )

  BeforeEach(func() {
    var err error

    githubClient = &fakes.FakeGitHub{}
    command = resource.NewOutCommand(githubClient, ioutil.Discard)

    sourcesDir, err = ioutil.TempDir("", "github-deployment")
    Ω(err).ShouldNot(HaveOccurred())
  })

  AfterEach(func() {
    Ω(os.RemoveAll(sourcesDir)).Should(Succeed())
  })

  Context("when creating a new deployment status", func() {
    BeforeEach(func() {
      idPath := filepath.Join(sourcesDir, "id")
      statePath := filepath.Join(sourcesDir, "state")

      file(idPath, "1234")
      file(statePath, "success")

      request = resource.OutRequest{
        Params: resource.OutParams{
          IDPath: "id",
          StatePath:  "state",
        },
      }

      githubClient.CreateDeploymentStatusReturns(&github.DeploymentStatus{
        ID: github.Int(12),
        State: github.String("success"),
        CreatedAt: &github.Timestamp{time.Date(2016, 01, 20, 15, 15, 15, 0, time.UTC)},
      }, nil)
    })

    It("creates a new status", func() {
      _, err := command.Run(sourcesDir, request)
      Ω(err).ShouldNot(HaveOccurred())

      Ω(githubClient.CreateDeploymentStatusCallCount()).Should(Equal(1))
      id, status := githubClient.CreateDeploymentStatusArgsForCall(0)

      Ω(id).Should(Equal(*github.Int(1234)))
      Ω(status.State).Should(Equal(github.String("success")))
    })

    It("returns some metadata", func() {
      outResponse, err := command.Run(sourcesDir, request)
      Ω(err).ShouldNot(HaveOccurred())

      Ω(outResponse.Metadata).Should(ConsistOf(
        resource.MetadataPair{Name: "id", Value: "12"},
        resource.MetadataPair{Name: "state", Value: "success"},
        resource.MetadataPair{Name: "created_at", Value: "2016-01-20 15:15:15"},
      ))
    })

    It("returns the version number of the deployment, not the status", func() {
      outResponse, err := command.Run(sourcesDir, request)
      Ω(err).ShouldNot(HaveOccurred())

      Ω(outResponse.Version).Should(Equal(
        resource.Version{
          ID: "1234",
        },
      ))
    })
  })
})
