// +build mage

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
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

	mg.Deps(func() error {
		cmd := exec.Command("brew", "install", "capnp")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		return err
	})

	mg.Deps(func() error {
		cmd := exec.Command("go", "install", "capnproto.org/go/capnp/v3/capnpc-go@latest")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		return err
	})

	mg.Deps(func() error {
		os.Setenv("GO111MODULE", "off")
		cmd := exec.Command("go", "get", "-u", "capnproto.org/go/capnp/v3/")
		err := cmd.Run()
		os.Unsetenv("GO111MODULE")
		if err != nil {
			fmt.Println(err)
		}
		return err
	})

	cmd := exec.Command("go", "get", "-v", "./...")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return cmd.Run()
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
	cmd := exec.Command("go", "clean", "-modcache")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	return err
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