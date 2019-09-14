package resource_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	resource "github.com/ahume/github-deployment-resource"
)

func file(path, contents string) {
	Ω(ioutil.WriteFile(path, []byte(contents), 0644)).Should(Succeed())
}

var _ = Describe("Resources", func() {

	var (
		p resource.OutRequest

		sourceDir string
	)

	BeforeEach(func() {
		var err error

		sourceDir, err = ioutil.TempDir("", "github-deployment")
		Ω(err).ShouldNot(HaveOccurred())

		// Um.
		os.Args = []string{"out", sourceDir}
	})

	AfterEach(func() {
		Ω(os.RemoveAll(sourceDir)).Should(Succeed())
	})

	Context("Params is unmarshalled", func() {

		BeforeEach(func() {
			p = resource.NewOutRequest()
		})

		It("adds default type", func() {
			r := bytes.NewReader([]byte(`{
				"params": {}
				}`))
			_ = json.NewDecoder(r).Decode(&p)
			Ω(*p.Params.Type).Should(Equal("status"))
		})

		It("gets values from strings", func() {
			r := bytes.NewReader([]byte(`{
				"params": {
					"type": "deployment",
					"ref": "ref-string",
					"state": "state-string",
					"task": "task-string",
					"environment": "environment-string",
					"description": "description-string"
					}
				}`))
			err := json.NewDecoder(r).Decode(&p)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(*p.Params.Type).Should(Equal("deployment"))
			Ω(*p.Params.Ref).Should(Equal("ref-string"))
			Ω(*p.Params.State).Should(Equal("state-string"))
			Ω(*p.Params.Task).Should(Equal("task-string"))
			Ω(*p.Params.Environment).Should(Equal("environment-string"))
			Ω(*p.Params.Description).Should(Equal("description-string"))
		})

		It("gets values from files", func() {
			idPath := filepath.Join(sourceDir, "id")
			refPath := filepath.Join(sourceDir, "ref")
			statePath := filepath.Join(sourceDir, "state")
			taskPath := filepath.Join(sourceDir, "task")
			environmentPath := filepath.Join(sourceDir, "environment")
			descriptionPath := filepath.Join(sourceDir, "description")

			file(idPath, "id-from-file")
			file(refPath, "ref-from-file")
			file(statePath, "state-from-file")
			file(taskPath, "task-from-file")
			file(environmentPath, "environment-from-file")
			file(descriptionPath, "description-from-file")

			r := bytes.NewReader([]byte(`{
				"params": {
					"type": "deployment",
					"id": {
						"file": "id"
					},
					"ref": {
						"file": "ref"
					},
					"state": {
						"file": "state"
					},
					"task": {
						"file": "task"
					},
					"environment": {
						"file": "environment"
					},
					"description": {
						"file": "description"
					}
				}
			}`))
			err := json.NewDecoder(r).Decode(&p)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(*p.Params.Type).Should(Equal("deployment"))
			Ω(*p.Params.ID).Should(Equal("id-from-file"))
			Ω(*p.Params.Ref).Should(Equal("ref-from-file"))
			Ω(*p.Params.State).Should(Equal("state-from-file"))
			Ω(*p.Params.Task).Should(Equal("task-from-file"))
			Ω(*p.Params.Environment).Should(Equal("environment-from-file"))
			Ω(*p.Params.Description).Should(Equal("description-from-file"))
		})

		It("gets raw payload", func() {
			r := bytes.NewReader([]byte(`{
				"params": {
					"type": "deployment",
					"payload": {
						"one":"two",
						"three":"four"
					}
				}
			}`))
			err := json.NewDecoder(r).Decode(&p)
			payload := *p.Params.Payload
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*p.Params.Type).Should(Equal("deployment"))
			Ω(payload["one"]).Should(Equal("two"))
			Ω(payload["three"]).Should(Equal("four"))
		})

		It("merges raw payload into file payload", func() {
			payloadPath := filepath.Join(sourceDir, "payload")
			file(payloadPath, `{"three":"four"}`)

			r := bytes.NewReader([]byte(`{
				"params": {
					"type": "deployment",
					"payload": {"one":"two"},
					"payload_path": "payload"
				}
			}`))
			err := json.NewDecoder(r).Decode(&p)
			payload := *p.Params.Payload
			Ω(err).ShouldNot(HaveOccurred())
			Ω(*p.Params.Type).Should(Equal("deployment"))
			Ω(payload["one"]).Should(Equal("two"))
			Ω(payload["three"]).Should(Equal("four"))
		})
	})
})
