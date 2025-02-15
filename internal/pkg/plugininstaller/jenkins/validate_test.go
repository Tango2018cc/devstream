package jenkins

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetJobDefaultConfig func", func() {
	var (
		jenkinsUser, jenkinsPassword, jenkinsURL, jenkinsFilePath, projectURL string
		options                                                               map[string]interface{}
	)
	BeforeEach(func() {
		jenkinsUser = "test"
		jenkinsPassword = "testPassword"
		jenkinsURL = "http://test.jenkins.com/"
		projectURL = "https://test.gitlab.com/test/test_project"
		jenkinsFilePath = "http://raw.content.com/Jenkinsfile"
		err := os.Setenv("JENKINS_PASSWORD", jenkinsPassword)
		Expect(err).NotTo(HaveOccurred())
		options = map[string]interface{}{
			"jenkins": map[string]interface{}{
				"url":  jenkinsURL,
				"user": jenkinsUser,
			},
			"scm": map[string]interface{}{
				"cloneURL": projectURL,
			},
			"pipeline": map[string]interface{}{
				"jenkinsfilePath": jenkinsFilePath,
			},
		}
	})
	When("repo url is not valie", func() {
		BeforeEach(func() {
			options["scm"] = map[string]interface{}{
				"cloneURL": "not_valid_url/gg",
			}
		})
		It("should return err", func() {
			_, err := SetJobDefaultConfig(options)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("all input is valid", func() {
		BeforeEach(func() {
			options["scm"] = map[string]interface{}{
				"cloneURL": "git@54.71.232.26:30022:root/spring-demo.git",
				"apiURL":   "http://www.app.com",
			}
		})
		It("should set default value", func() {
			newOptions, err := SetJobDefaultConfig(options)
			Expect(err).Error().ShouldNot(HaveOccurred())
			opts, err := newJobOptions(newOptions)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(opts.CIConfig).ShouldNot(BeNil())
			Expect(opts.Pipeline.JobName).Should(Equal("spring-demo"))
			Expect(opts.BasicAuth).ShouldNot(BeNil())
			Expect(opts.BasicAuth.Username).Should(Equal(jenkinsUser))
			Expect(opts.BasicAuth.Password).Should(Equal(jenkinsPassword))
			Expect(opts.ProjectRepo).ShouldNot(BeNil())
			Expect(opts.ProjectRepo.Repo).Should(Equal("spring-demo"))
		})
	})
	AfterEach(func() {
		os.Unsetenv("JENKINS_PASSWORD")
	})
})

var _ = Describe("generateRandomSecretToken func", func() {
	It("should return random str", func() {
		token := generateRandomSecretToken()
		Expect(token).ShouldNot(BeEmpty())
	})
})

var _ = Describe("ValidateJobConfig func", func() {
	var (
		jenkinsUser, jenkinsURL, jenkinsFilePath, projectURL string
		options                                              map[string]interface{}
	)
	BeforeEach(func() {
		jenkinsUser = "test"
		jenkinsURL = "http://test.jenkins.com/"
		projectURL = "https://test.gitlab.com/test/test_project"
		jenkinsFilePath = "http://raw.content.com/Jenkinsfile"
		options = map[string]interface{}{
			"jenkins": map[string]interface{}{
				"url":  jenkinsURL,
				"user": jenkinsUser,
			},
			"scm": map[string]interface{}{
				"cloneURL": projectURL,
			},
			"pipeline": map[string]interface{}{
				"jenkinsfilePath": jenkinsFilePath,
			},
		}
	})
	When("Input field miss", func() {
		BeforeEach(func() {
			options = map[string]interface{}{
				"jenkins": map[string]interface{}{
					"url":  jenkinsURL,
					"user": jenkinsUser,
				},
				"scm": map[string]interface{}{
					"cloneURL": projectURL,
				},
			}
		})
		It("should return error", func() {
			_, err := ValidateJobConfig(options)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
})
