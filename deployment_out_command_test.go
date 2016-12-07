package resource_test

import (
  "io/ioutil"
  "os"
  "path/filepath"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "github.com/ahume/go-github/github"

  "github.com/ahume/github-deployment-resource"
  "github.com/ahume/github-deployment-resource/fakes"
)

var _ = Describe("Deployment Out Command", func() {
  var (
    command      *resource.DeploymentOutCommand
    githubClient *fakes.FakeGitHub

    sourcesDir string

    request resource.OutRequest
  )

  BeforeEach(func() {
    var err error

    githubClient = &fakes.FakeGitHub{}
    command = resource.NewDeploymentOutCommand(githubClient, ioutil.Discard)

    sourcesDir, err = ioutil.TempDir("", "github-deployment")
    Ω(err).ShouldNot(HaveOccurred())
  })

  AfterEach(func() {
    Ω(os.RemoveAll(sourcesDir)).Should(Succeed())
  })

  Context("when creating a new deployment", func() {
    BeforeEach(func() {

      request = resource.OutRequest{
        Params: resource.OutParams{
          Ref: "refparam",
          Task: "taskparam",
          Description: "descparam",
          Environment: "envparam",
        },
      }

      githubClient.CreateDeploymentReturns(&github.Deployment{
        ID: github.Int(1),
        Ref: github.String("refparam"),
        SHA: github.String("1234"),
        Task: github.String("taskparam"),
        Description: github.String("descparam"),
        Environment: github.String("envparam"),
        Creator: &github.User{
          Login: github.String("theboss"),
        },
      }, nil)
    })

    It("creates a new deployment", func() {
      _, err := command.Run(sourcesDir, request)
      Ω(err).ShouldNot(HaveOccurred())

      Ω(githubClient.CreateDeploymentCallCount()).Should(Equal(1))
      deployment := githubClient.CreateDeploymentArgsForCall(0)

      Ω(deployment.Ref).Should(Equal(github.String("refparam")))
      Ω(deployment.Task).Should(Equal(github.String("taskparam")))
      Ω(deployment.Description).Should(Equal(github.String("descparam")))
      Ω(deployment.Environment).Should(Equal(github.String("envparam")))
    })

    It("returns some metadata", func() {
      outResponse, err := command.Run(sourcesDir, request)
      Ω(err).ShouldNot(HaveOccurred())

      Ω(outResponse.Metadata).Should(ConsistOf(
        resource.MetadataPair{Name: "id", Value: "1"},
        resource.MetadataPair{Name: "ref", Value: "refparam"},
        resource.MetadataPair{Name: "sha", Value: "1234"},
        resource.MetadataPair{Name: "task", Value: "taskparam"},
        resource.MetadataPair{Name: "environment", Value: "envparam"},
        resource.MetadataPair{Name: "description", Value: "descparam"},
        resource.MetadataPair{Name: "creator", Value: "theboss"},
      ))
    })

    It("returns the new version number", func() {
      outResponse, err := command.Run(sourcesDir, request)
      Ω(err).ShouldNot(HaveOccurred())

      Ω(outResponse.Version).Should(Equal(
        resource.Version{
          ID: "1",
        },
      ))
    })

    Context("when output paths are provided use them", func() {
      BeforeEach(func() {
        refPath := filepath.Join(sourcesDir, "ref")
        taskPath := filepath.Join(sourcesDir, "task")
        payloadPath := filepath.Join(sourcesDir, "payload")
        envPath := filepath.Join(sourcesDir, "environment")
        descPath := filepath.Join(sourcesDir, "description")

        file(refPath, "ref-from-file")
        file(taskPath, "task-from-file")
        file(payloadPath, "{\"nice\": \"one\"}")
        file(envPath, "env-from-file")
        file(descPath, "description-from-file")

        request = resource.OutRequest{
          Params: resource.OutParams{
            RefPath: "ref",
            TaskPath: "task",
            PayloadPath: "payload",
            EnvironmentPath: "environment",
            DescriptionPath: "description",
          },
        }

        githubClient.CreateDeploymentReturns(&github.Deployment{
          ID: github.Int(1),
          Ref: github.String("ref-from-file"),
          SHA: github.String("1234"),
          Task: github.String("task-from-file"),
          Environment: github.String("env-from-file"),
          Creator: &github.User{
            Login: github.String("theboss"),
          },
        }, nil)
      })

      It("creates a new deployment", func() {
        _, err := command.Run(sourcesDir, request)
        Ω(err).ShouldNot(HaveOccurred())

        Ω(githubClient.CreateDeploymentCallCount()).Should(Equal(1))
        deployment := githubClient.CreateDeploymentArgsForCall(0)

        Ω(deployment.Ref).Should(Equal(github.String("ref-from-file")))
        Ω(deployment.Task).Should(Equal(github.String("task-from-file")))
        Ω(deployment.Environment).Should(Equal(github.String("env-from-file")))
        Ω(deployment.Payload).Should(Equal(github.String("{\"nice\": \"one\"}")))
        Ω(deployment.Description).Should(Equal(github.String("description-from-file")))
      })

      It("returns some metadata", func() {
        outResponse, err := command.Run(sourcesDir, request)
        Ω(err).ShouldNot(HaveOccurred())

        Ω(outResponse.Metadata).Should(ConsistOf(
          resource.MetadataPair{Name: "id", Value: "1"},
          resource.MetadataPair{Name: "ref", Value: "ref-from-file"},
          resource.MetadataPair{Name: "sha", Value: "1234"},
          resource.MetadataPair{Name: "task", Value: "task-from-file"},
          resource.MetadataPair{Name: "environment", Value: "env-from-file"},
          resource.MetadataPair{Name: "creator", Value: "theboss"},
        ))
      })
    })
  })
})
