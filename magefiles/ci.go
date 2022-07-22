//go:build mage

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	concourseTargetUrl = "https://ci.a13s.io"
	concourseTarget = "-tp"
	acrossStepFlag     = "--enable-across-step"

	renderVarsShTemplate = `#!/bin/sh
set -e
[ -e {{ .Filename }}.yaml ] && rm -- {{ .Filename }}.yaml
for key in $(vault kv list -format=json concourse/main/{{ .Pipeline }} | jq -r '.[]'); do
	echo "rendering ${key}"
	vault kv get -format=yaml concourse/main/{{ .Pipeline }}/"${key}" | KEY=$key yq '{env(KEY): .data.data}' >> {{ .Filename }}.yaml
done`
)

var (
	knownPipelines = []string{"images", "pleiades", "minio"}
)

type CI mg.Namespace

// install and setup the concourse cli (fly)
func (CI) Setup() error {
	fmt.Println("installing fly cli")

	err := installConcourseCli()
	if err != nil {
		return err
	}

	err = concourseLogin()
	if err != nil {
		return err
	}

	return nil
}

// validate a pipeline's config
func (CI) Validate(pipelineName string) error {

	known := false
	for idx := range knownPipelines {
		if pipelineName == knownPipelines[idx] {
			known = true
		}
	}

	if !known {
		return fmt.Errorf("pipeline %s isn't supported in the build system", pipelineName)
	}

	return validatePipeline(pipelineName)
}

// render the pipeline variables
func (CI) Rendervars() error {

	// vault login -method=token token="${VAULT_READONLY_TOKEN}"
	readonlyToken := os.ExpandEnv("${VAULT_READONLY_TOKEN}")

	err := sh.RunWithV(nil, "vault", "login", "-no-print", "-method=token", fmt.Sprintf("token=%s", readonlyToken))
	if err != nil {
		return err
	}

	for idx := range knownPipelines {
		if err := renderPipelineVars(knownPipelines[idx]); err != nil {
			return err
		}
	}

	return nil
}

// run a oneoff pipeline task. $1 is the task name (not the filename) and $2 is any extra inputs in "key=value,key2=value2" format ("pleiades=." is included by default)
func (CI) Run(taskName string, extraInputs string) error {
	pleiadesInput := fmt.Sprintf("pleiades=.")
	otherInputs := strings.Split(extraInputs, ",")

	// make the task inputs
	inputs := make([]string, 0)
	inputs = append(inputs, "-i", pleiadesInput)
	for idx := range otherInputs {
		inputs = append(inputs, "-i", otherInputs[idx])
	}

	taskName = fmt.Sprintf("ci/tasks/%s.yaml", taskName)

	// clear out the all.yaml file, just in case
	allVarsLocation := "ci/vars/all.yaml"
	if err := sh.Rm(allVarsLocation); err != nil {
		return err
	}

	// render the vars into a single vars file
	files, err := filepath.Glob("ci/vars/*.yaml")
	if err != nil {
		return err
	}
	yqArgs := make([]string, 0)
	yqArgs = append(yqArgs, "ea", `. as $item ireduce ({}; . * $item )`)
	yqArgs = append(yqArgs, files...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("yq", yqArgs...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error when running yq command: %s", stderr.String())
	}

	if err := ioutil.WriteFile(allVarsLocation, stdout.Bytes(), 0644); err != nil {
		return err
	}

	args := make([]string, 0)
	args = append(args, concourseTarget, "execute", "-c", taskName, "-l", allVarsLocation)

	return sh.RunWithV(nil, "fly", args...)
}

func renderPipelineVars(targetPipeline string) error {
	renderScriptLocation := fmt.Sprintf("ci/vars/render-%s-vars.sh", targetPipeline)
	renderedVarsLocation := fmt.Sprintf("ci/vars/%s", targetPipeline)

	renderData := struct{
		Filename string
		Pipeline string
	}{
		Filename: renderedVarsLocation,
		Pipeline: targetPipeline,
	}

	tmpl := template.New("rendered")
	tmpl, err := tmpl.Parse(renderVarsShTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, renderData); err != nil {
		return err
	}

	fmt.Printf("rendering variables for %s pipeline\n", targetPipeline)

	err = ioutil.WriteFile(renderScriptLocation, buf.Bytes(), 0677)
	if err != nil {
		return fmt.Errorf("error while working on %s pipeline vars", targetPipeline)
	}
	defer sh.Rm(renderScriptLocation)

	err = sh.RunWithV(nil, "/bin/sh", renderScriptLocation)
	if err != nil {
		return err
	}

	return nil
}

func validatePipeline(pipeline string) error {
	renderPipelineVars(pipeline)

	imagesPipelineFile := fmt.Sprintf("ci/pipelines/%s.yaml", pipeline)
	varsFile := fmt.Sprintf("ci/vars/%s.yaml", pipeline)

	return sh.RunWithV(nil, "fly", concourseTarget, "validate-pipeline", "-c", imagesPipelineFile, "-l", varsFile, acrossStepFlag)
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
	_, err := os.Stat("build/fly")
	if !os.IsNotExist(err) {
		fmt.Println("the cli exists, skipping source build")
		fmt.Println("delete build/fly if there's something wrong with the binary")
		return nil
	}

	fmt.Println("querying concourse for version info")
	client := http.DefaultClient

	earl, err := url.Parse(concourseTargetUrl)
	if err != nil {
		return nil
	}

	req := http.Request{
		Method: "GET",
		URL:    earl,
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
		URL:           "https://github.com/concourse/concourse",
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
	err = sh.RunWithV(nil, "go", "build", "-v", "-o", fmt.Sprintf("%s/fly", targetPath), "main.go")
	if err != nil {
		return err
	}

	err = os.Chdir("../..")
	if err != nil {
		return err
	}

	return os.RemoveAll(concourseDirName)
}

func concourseLogin() error {
	fmt.Println("logging into concourse")
	err := sh.RunWithV(nil, "fly", concourseTarget, "login", "-c", concourseTargetUrl)
	if err != nil {
		return err
	}
	return nil
}
