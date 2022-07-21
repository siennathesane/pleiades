//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

var (
	homebrewTargets = []string{
		"capnp",
	}
)

type Install mg.Namespace

// install pleiades to your local directory
func (Install) Local() error {
	mg.Deps(Install.Deps, Build.Compile)
	fmt.Println("installing...")
	return os.Rename("build/pleiades", "/usr/bin/pleiades")
}

// install the binary to a homebrew location
func (Install) Homebrew(path string) error {
	mg.Deps(Build.Compile)
	fmt.Println("installing to homebrew...")
	return os.Rename("build/pleiades", path)
}

// install necessary tools and dependencies to develop pleiades
func (Install) Deps() error {
	fmt.Println("installing tools...")

	// each of these should be their own dep :shrug:
	mg.Deps(func() error {
		for idx := range homebrewTargets {
			if err := sh.RunWithV(nil, "brew", "install", homebrewTargets[idx]); err != nil {
				return err
			}
		}
		return nil
	})

	mg.Deps(func() error {
		if err := sh.RunWithV(nil, "go", "install", "github.com/spf13/cobra-cli@latest"); err != nil {
			return err
		}
		return nil
	})

	mg.Deps(func() error {
		fmt.Println("installing capn' proto compiler")
		return sh.RunWithV(nil, "brew", "install", "capnp")
	})

	mg.Deps(func() error {
		fmt.Println("installing capn' proto go compiler plugin")
		return sh.RunWithV(nil, "go", "install", "capnproto.org/go/capnp/v3/capnpc-go@latest")
	})

	mg.Deps(func() error {
		fmt.Println("installing cap'n proto golang compiler cli")
		return sh.RunWithV(map[string]string{
			//"GO111MODULE": "off",
		}, "go", "get", "-u", "capnproto.org/go/capnp/v3")
	})

	mg.Deps(func() error {
		fmt.Println("getting pleiades deps")
		return sh.RunWithV(nil, "go", "get", "-u", "./...")
	})

	return nil
}

type Build mg.Namespace

// compile pleiades with the local build information
func (Build) Compile() error {
	fmt.Println("compiling...")
	return sh.RunWithV(nil, "go", "build", fmt.Sprintf("-ldflags=%s", ldflags()), "-o", "build/pleiades", "./main.go")
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

type Clean mg.Namespace

// clear the local build directory
func (Clean) Cache() error {
	fmt.Println("Cleaning...")
	return os.RemoveAll("build")
}

// clear all tools and dependencies
func (Clean) All() error {
	fmt.Println("removing build directory...")
	err := os.RemoveAll("build")
	if err != nil {
		return err
	}

	fmt.Println("cleaning mod cache...")
	if err := sh.RunWithV(nil, "go", "clean", "-modcache"); err != nil {
		return err
	}

	fmt.Println("removing homebrew tools")
	for idx := range homebrewTargets {
		if err := sh.RunWithV(nil, "brew", "remove", homebrewTargets[idx]); err != nil {
			return err
		}
	}
	return nil
}

type Gen mg.Namespace

// generate all schemas
func (Gen) All() {
	mg.SerialDeps(Gen.Host, Gen.Database)
}

// compiles the database schemas and generates the go code
func (Gen) Database() error {
	gopath := os.Getenv("GOPATH")

	fmt.Println("generating database protocols")
	files, err := filepath.Glob("protocols/v1/database/*.capnp")
	if err != nil {
		return err
	}

	args := []string{"compile", fmt.Sprintf("-I%s/src/capnproto.org/go/capnp/std", gopath), "-ogo:pkg"}
	args = append(args, files...)
	return sh.RunWithV(nil, "capnp", args...)
}

// compiles the host schemas and generates the go code
func (Gen) Host() error {
	gopath := os.Getenv("GOPATH")

	fmt.Println("generating host protocols")
	files, err := filepath.Glob("protocols/v1/host/*.capnp")
	if err != nil {
		return err
	}

	args := []string{"compile", fmt.Sprintf("-I%s/src/capnproto.org/go/capnp/std", gopath), "-ogo:pkg"}
	args = append(args, files...)
	return sh.RunWithV(nil, "capnp", args...)
}
