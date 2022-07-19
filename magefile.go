// +build mage

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/mage"
)

var (
	homebrewTargets = []string{
		"capnp",
	}
)

// compile pleiades with the local build information
func Build() error {
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "build/pleiades", "main.go")
	return cmd.Run()
}

// install pleiades to your local directory
func Install() error {
	mg.Deps(InstallDeps, Build)
	fmt.Println("installing...")
	return os.Rename("build/pleiades", "/usr/bin/pleiades")
}

// install necessary tools and dependencies to develop pleiades
func InstallDeps() error {
	fmt.Println("installing tools...")

	// each of these should be their own dep :shrug:
	mg.Deps(func () error {
		for idx := range homebrewTargets {
			if err := sh.RunWithV(nil, "brew", "install", homebrewTargets[idx]); err != nil {
				return err
			}
		}
		return nil
	})

	mg.Deps(func() error {
		fmt.Println("installing capn' proto compiler")
		return run("brew install capnp")
	})

	mg.Deps(func() error {
		fmt.Println("installing capn' proto go compiler plugin")
		return sh.RunWithV(nil, "go", "install", "capnproto.org/go/capnp/v3/capnpc-go@latest")
	})

	mg.Deps(func() error {
		fmt.Println("installing cap'n proto golang compiler cli")
		return sh.RunWithV(map[string]string{
			"GO111MODULE": "off",
		}, "go", "get", "-u", "capnproto.org/go/capnp/v3/")
	})

	mg.Deps(func() error {
		fmt.Println("getting pleiades deps")
		return run("go get -v ./...")
	})

	return nil
}

// quickly recompile pleiades
func Rebuild() error {
	fmt.Println("cleaning...")
	os.RemoveAll("build")
	cmd := exec.Command("go", "clean")
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("building...")
	cmd = exec.Command("go", "build", "-o", "build/pleiades", "main.go")
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return cmd.Run()
}

// clear the local build directory
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("build")
}

// Reset your local state
func Reset() error {
	fmt.Println("removing build directory...")
	os.RemoveAll("build")

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

// compiles the local capnp schemas and generates the go code
func Generate() error {
	fmt.Println("generating host protocols")
	gopath := os.Getenv("GOPATH")

	files, err := filepath.Glob("protocols/v1/host/*.capnp")
	if err != nil {
		return err
	}

	args := []string{"compile", fmt.Sprintf("-I%s/src/capnproto.org/go/capnp/std", gopath), "-ogo:pkg"}
	cmd := exec.Command("capnp", append(args, files...)...)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	if err != nil {
		fmt.Println(stderrBuf.String())
	}
	return err
}

func run(shellCmd string) error {
	cmd := exec.Command("bash", "-c", shellCmd)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	if err != nil {
		fmt.Println(stderrBuf.String())
	}
	return err
}
