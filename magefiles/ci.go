//go:build mage
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type CI mg.Namespace

func (CI) Setup() error {
	fmt.Println("installing fly cli")
	return installConcourseCli()
}

func (CI) Validate(pipelineName string) error {
	fmt.Printf("validating %s pipeline", pipelineName)
	return sh.RunWithV(nil, "fly")
}

type concourseVersionInfo struct {
	ExternalURL  string `json:"external_url"`
	FeatureFlags struct {
		AcrossStep           bool `json:"across_step"`
		BuildRerun           bool `json:"build_rerun"`
		CacheStreamedVolumes bool `json:"cache_streamed_volumes"`
		GlobalResources      bool `json:"global_resources"`
		PipelineInstances    bool `json:"pipeline_instances"`
		RedactSecrets        bool `json:"redact_secrets"`
		ResourceCausality    bool `json:"resource_causality"`
	} `json:"feature_flags"`
	Version       string `json:"version"`
	WorkerVersion string `json:"worker_version"`
}

func installConcourseCli() error {
	fmt.Println("querying concourse for version info")
	client := http.DefaultClient

	earl, err := url.Parse("https://ci.github.com/mxplusb/pleiades/pkg/api/v1/info")
	if err != nil {
		return nil
	}

	req := http.Request{
		Method: "GET",
		URL: earl,
	}

	resp, err := client.Do(&req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var concourseVersion concourseVersionInfo
	if err := json.Unmarshal(body, &concourseVersion); err != nil {
		return err
	}

	fmt.Printf("found version v%s\n", concourseVersion.Version)
	fmt.Printf("downloading fly source version v%s\n", concourseVersion.Version)

	concourseDirName := "concourse-tmp"

	_, err = os.Stat(concourseDirName)
	if !os.IsNotExist(err) {
		sh.Rm(concourseDirName)
	}

	_, err = git.PlainClone(concourseDirName, false, &git.CloneOptions{
		URL:               "https://github.com/concourse/concourse",
		ReferenceName: plumbing.NewTagReferenceName(fmt.Sprintf("v%s", concourseVersion.Version)),
	})
	if err != nil {
		return err
	}

	os.Chdir(concourseDirName)

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println("downloading dependencies")
	err = sh.RunWithV(nil, "go", "get", "-v", "./...")

	// change to the fly folder
	os.Chdir("fly")

	// this only works if you run it from the root directory
	targetPath := filepath.Join(currDir, "..", "build")
	fmt.Printf("target build dir: %s\n", targetPath)

	fmt.Println("compiling fly cli")
	err = sh.RunWithV(nil, "go", "build", "-o", fmt.Sprintf("%s/fly", targetPath), "main.go")
	if err != nil {
		return err
	}

	err = os.Chdir("../..")
	if err != nil {
		return err
	}

	return os.Remove(concourseDirName)
}