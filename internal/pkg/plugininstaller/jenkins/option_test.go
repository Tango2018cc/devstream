package jenkins

import (
	"context"
	"errors"
	"fmt"

	"github.com/bndr/gojenkins"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/jenkins/dingtalk"
)

// create mock client
var testError = errors.New("test")

type mockSuccessJenkinsClient struct {
}

func (m *mockSuccessJenkinsClient) CreateOrUpdateJob(string, string) (job *gojenkins.Job, created bool, err error) {
	return nil, true, nil
}

func (m *mockSuccessJenkinsClient) ExecuteScript(string) (string, error) {
	return "", nil
}
func (m *mockSuccessJenkinsClient) GetJob(context.Context, string, ...string) (*gojenkins.Job, error) {
	return nil, nil
}
func (m *mockSuccessJenkinsClient) DeleteJob(context.Context, string) (bool, error) {
	return true, nil
}
func (m *mockSuccessJenkinsClient) InstallPluginsIfNotExists(plugin []*jenkins.JenkinsPlugin, enableRestart bool) error {
	return nil
}
func (m *mockSuccessJenkinsClient) CreateGiltabCredential(string, string) error {
	return nil
}

func (m *mockSuccessJenkinsClient) CreateSSHKeyCredential(id, userName, privateKey string) error {
	return nil
}

func (m *mockSuccessJenkinsClient) ConfigCasc(cascScript string) error {
	return nil
}

func (m *mockSuccessJenkinsClient) ApplyDingTalkBot(dingtalk.BotConfig) error {
	return nil
}

type mockErrorJenkinsClient struct {
}

func (m *mockErrorJenkinsClient) CreateOrUpdateJob(string, string) (job *gojenkins.Job, created bool, err error) {
	return nil, true, testError
}

func (m *mockErrorJenkinsClient) ExecuteScript(string) (string, error) {
	return "", testError
}
func (m *mockErrorJenkinsClient) GetJob(context.Context, string, ...string) (*gojenkins.Job, error) {
	return nil, testError
}
func (m *mockErrorJenkinsClient) DeleteJob(context.Context, string) (bool, error) {
	return false, testError
}
func (m *mockErrorJenkinsClient) InstallPluginsIfNotExists(plugin []*jenkins.JenkinsPlugin, enableRestart bool) error {
	return testError
}
func (m *mockErrorJenkinsClient) CreateGiltabCredential(string, string) error {
	return testError
}
func (m *mockErrorJenkinsClient) ConfigCasc(string) error {
	return testError
}

func (m *mockErrorJenkinsClient) ApplyDingTalkBot(dingtalk.BotConfig) error {
	return testError
}

func (m *mockErrorJenkinsClient) CreateSSHKeyCredential(id, userName, privateKey string) error {
	return testError
}

var _ = Describe("JobOptions struct", func() {
	var (
		jenkinsURL, secretToken, jobName, projectURL, jenkinsFilePath, userName, password, repoOwner, repoName string
		jobOptions                                                                                             *JobOptions
		basicAuth                                                                                              *jenkins.BasicAuth
		projectRepo                                                                                            *common.Repo
		ciConfig                                                                                               *ci.CIConfig
		mockClient                                                                                             jenkins.JenkinsAPI
	)
	BeforeEach(func() {
		jenkinsURL = "http://test.com"
		userName = "test_user"
		password = "test_password"
		repoOwner = "owner"
		repoName = "repo"
		jobName = "test_folder/test_job"
		projectURL = "http://127.0.0.1:300/test/project"
		jenkinsFilePath = "http://raw.content.com/Jenkinsfile"
		basicAuth = &jenkins.BasicAuth{
			Password: password,
			Username: userName,
		}
		projectRepo = &common.Repo{
			Owner:   repoOwner,
			Repo:    repoName,
			Branch:  "test",
			BaseURL: "http://127.0.0.1:300",
		}
		ciConfig = &ci.CIConfig{
			Type:      "jenkins",
			RemoteURL: jenkinsFilePath,
		}
		secretToken = "secret"
		jobOptions = &JobOptions{
			Jenkins: Jenkins{
				URL:           jenkinsURL,
				User:          userName,
				Namespace:     "jenkins",
				EnableRestart: false,
			},
			SCM: SCM{
				CloneURL: projectURL,
				Branch:   "test",
			},
			Pipeline: Pipeline{
				JobName:         jobName,
				JenkinsfilePath: jenkinsFilePath,
				ImageRepo:       ImageRepo{},
			},
			BasicAuth:   basicAuth,
			ProjectRepo: projectRepo,
			CIConfig:    ciConfig,
			SecretToken: secretToken,
		}
	})
	Context("encode method", func() {
		It("should work noraml", func() {
			_, err := jobOptions.encode()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	Context("createOrUpdateJob method", func() {
		When("jenkins client return normal", func() {
			BeforeEach(func() {
				mockClient = &mockSuccessJenkinsClient{}
			})
			It("should work noraml", func() {
				err := jobOptions.createOrUpdateJob(mockClient)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("jenkins client return error", func() {
			BeforeEach(func() {
				mockClient = &mockErrorJenkinsClient{}
			})
			It("should return error", func() {
				err := jobOptions.createOrUpdateJob(mockClient)
				Expect(err).Should(HaveOccurred())
			})
		})
	})
	Context("buildWebhookInfo method", func() {
		It("should work normal", func() {
			webHookInfo := jobOptions.buildWebhookInfo()
			Expect(webHookInfo.Address).Should(Equal(fmt.Sprintf("%s/project/%s", jobOptions.Jenkins.URL, jobOptions.getJobPath())))
			Expect(webHookInfo.SecretToken).Should(Equal(secretToken))
		})
	})
	Context("installPlugins method", func() {
		When("jenkins client return error", func() {
			BeforeEach(func() {
				mockClient = &mockErrorJenkinsClient{}
			})
			It("should return error", func() {
				err := jobOptions.installPlugins(mockClient, []*jenkins.JenkinsPlugin{
					{Name: "test_plugin", Version: "123"},
				})
				Expect(err).Error().Should(HaveOccurred())
			})
		})
	})
	Context("createGitlabConnection method", func() {
		When("jenkins client return error", func() {
			BeforeEach(func() {
				mockClient = &mockErrorJenkinsClient{}
			})
			It("should return error", func() {
				err := jobOptions.createGitlabConnection(mockClient, "casc")
				Expect(err).Error().Should(HaveOccurred())
			})
		})
	})
	Context("deleteJob method", func() {
		When("jenkins client get job error", func() {
			BeforeEach(func() {
				mockClient = &mockErrorJenkinsClient{}
			})
			It("should return error", func() {
				err := jobOptions.deleteJob(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

	Context("buildCIConfig method", func() {
		When("jenkinsfilePath is local path", func() {
			BeforeEach(func() {
				jobOptions.Pipeline.JenkinsfilePath = "test/local"
			})
			It("should use localPath", func() {
				jobOptions.buildCIConfig()
				Expect(jobOptions.CIConfig.LocalPath).Should(Equal(jobOptions.Pipeline.JenkinsfilePath))
				Expect(jobOptions.CIConfig.RemoteURL).Should(BeEmpty())
			})
		})
		When("jenkinsfilePath is remote url", func() {
			BeforeEach(func() {
				jobOptions.Pipeline.JenkinsfilePath = "http://www.test.com/Jenkinsfile"
			})
			It("should use remote url", func() {
				jobOptions.buildCIConfig()
				Expect(jobOptions.CIConfig.LocalPath).Should(BeEmpty())
				Expect(jobOptions.CIConfig.RemoteURL).Should(Equal(jobOptions.CIConfig.RemoteURL))
				Expect(string(jobOptions.CIConfig.Type)).Should(Equal("jenkins"))
			})
		})
	})

})
