package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

const (
	appNamePlaceHolder = "_app_name_"
)

func filterGitFiles(filePath string, isDir bool) bool {
	if isDir {
		return false
	}
	if strings.Contains(filePath, ".git/") {
		log.Debugf("Walk: ignore file %s.", "./git/xxx")
		return false
	}

	if strings.HasSuffix(filePath, "README.md") {
		log.Debugf("Walk: ignore file %s.", "README.md")
		return false
	}
	return true
}

func getRepoFileNameFunc(appName, repoName string) file.DirFileNameFunc {
	return func(filePath, srcPath string) string {
		relativePath, _ := filepath.Rel(srcPath, filePath)
		if strings.HasPrefix(relativePath, repoName) {
			repoNamePath := fmt.Sprintf("%s/", repoName)
			relativePath = strings.Replace(relativePath, repoNamePath, "", 1)
		}
		replacedFileName := file.ReplaceAppNameInPathStr(relativePath, appNamePlaceHolder, appName)
		return strings.TrimSuffix(replacedFileName, ".tpl")
	}
}

func processRepoFileFunc(appName string, vars map[string]interface{}) file.DirFileProcessFunc {
	return func(filePath string) ([]byte, error) {
		var fileContent []byte
		log.Debugf("Walk: found file: %s.", filePath)
		// if file endswith tpl, render this file, else copy this file directly
		if strings.Contains(filePath, "tpl") {
			fileContentStr, err := template.New().FromLocalFile(filePath).SetDefaultRender(appName, vars).Render()
			if err != nil {
				log.Warnf("repo render file failed: %s", err)
				return fileContent, err
			}
			return []byte(fileContentStr), nil
		}
		return os.ReadFile(filePath)
	}
}
