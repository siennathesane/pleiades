//go:build mage
package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

// compile pleiades with the local build information
func (Build) Compile() error {
	fmt.Println("compiling...")
	return compileWithPath("build/pleiades")
}

// compile pleiades with the local build information
func compileWithPath(path string) error {
	return sh.RunWithV(nil, "go", "build", fmt.Sprintf("-ldflags=%s", ldflags()), "-o", path, "./main.go")
}

func ldflags() string {
	fmt.Println("generating ldflags...")

	writeComma := func(sb *strings.Builder) {
		if sb.Len() > 0 {
			sb.WriteString(" ")
		}
	}

	headReceiver := make(chan string)
	dirtyHead := make(chan bool)
	go func(hr chan string) {
		fmt.Println("getting git head...")
		localRepo, err := git.PlainOpen(".")
		if err != nil {
			hr <- ""
			return
		}

		head, err := localRepo.Head()
		if err != nil {
			hr <- ""
			return
		}
		fmt.Printf("got git head: %s\n", head.Hash().String())
		hr <- head.Hash().String()

		worktreeStatus, err := localRepo.Worktree()
		if err != nil {
			hr <- ""
			return
		}

		status, err := worktreeStatus.Status()
		if err != nil {
			hr <- ""
			return
		}

		if status.IsClean() {
			dirtyHead <- false
			return
		}
		dirtyHead <- true
	}(headReceiver)

	sb := strings.Builder{}

	sb.WriteString("-X '")
	sb.WriteString("github.com/mxplusb/pleiades/pkg.GoVersion=")
	sb.WriteString(runtime.Version())
	sb.WriteString("'")
	writeComma(&sb)

	now := time.Now().Format(time.RFC3339)
	fmt.Printf("using build time: %s\n", now)

	sb.WriteString("-X '")
	sb.WriteString("github.com/mxplusb/pleiades/pkg.BuildTime=")
	sb.WriteString(now)
	sb.WriteString("'")
	writeComma(&sb)

	localHash := <-headReceiver
	shortHead := localHash[len(localHash)-7:]
	fmt.Printf("using git hash: %s\n", shortHead)
	sb.WriteString("-X '")
	sb.WriteString("github.com/mxplusb/pleiades/pkg.Sha=")
	sb.WriteString(shortHead)
	sb.WriteString("'")

	headIsDirty := <-dirtyHead
	fmt.Printf("is head dirty: %v\n", headIsDirty)
	sb.WriteString("-X '")
	sb.WriteString("github.com/mxplusb/pleiades/pkg.Dirty=")
	sb.WriteString(strconv.FormatBool(headIsDirty))
	sb.WriteString("'")

	close(headReceiver)

	fmt.Printf("using ldflags: %s\n", sb.String())

	return sb.String()
}

// clean rebuild of pleiades
func (Build) Rebuild() error {
	fmt.Println("cleaning...")
	err := sh.Rm("build")
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "clean")
	err = cmd.Run()
	if err != nil {
		return err
	}
	mg.Deps(Build.Compile)
	return nil
}

// lint the repo
func (Build) Vet() error {
	fmt.Println("running linter")
	return sh.RunWithV(nil, "go", "vet", "./...")
}