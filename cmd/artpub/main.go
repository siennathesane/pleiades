/*
 * Copyright (c) 2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"oras.land/oras/cmd/oras/root"
)

var (
	oses          = []string{"linux", "darwin", "windows"}
	architectures = []string{"amd64", "arm64"}
)

func main() {

	srcDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get working directory: %w\n", err)
		os.Exit(1)
	}

	val := os.Getenv("GITHUB_WORKSPACE")
	if val == "" {
		fmt.Printf("GITHUB_WORKSPACE not set\n")
		os.Exit(1)
	}

	err = os.Chdir(val)
	if err != nil {
		fmt.Printf("failed to change directory: %w\n", err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			fmt.Printf("failed to change directory: %w\n", err)
		}
	}(srcDir)

	art := Artifact{
		SchemaVersion: 2,
		Layers:        make([]Layer, 0),
		MediaType:     "application/vnd.oci.image.manifest.v1+json",
		Config: Config{
			MediaType: "application/vnd.oci.image.config.v1+json",
		},
		Annotations: Annotations{
			OrgOpencontainersImageSource:      "https://github.com/mxplusb/pleiades",
			OrgOpencontainersImageDescription: "Pleiades, a constellation mesh database by @mxplusb",
			OrgOpencontainersImageLicenses:    "PolyForm-Internal-Use-1.0.0",
		},
	}

	buildArtifacts := make([]string, 0)
	for _, opSys := range oses {
		for _, arch := range architectures {
			target := fmt.Sprintf("pleiades-%s-%s", opSys, arch)
			buildArtifacts = append(buildArtifacts, target)

			art.Layers = append(art.Layers, Layer{
				MediaType:    "application/vnd.oci.image.layer.v1.tar",
				Os:           opSys,
				Architecture: arch,
				Annotations: Annotations{
					OrgOpenContainersImageTitle:       target,
					OrgOpencontainersImageSource:      "https://github.com/mxplusb/pleiades",
					OrgOpencontainersImageDescription: "Pleiades, a constellation mesh database by @mxplusb",
					OrgOpencontainersImageLicenses:    "PolyForm-Internal-Use-1.0.0",
				},
			})

			_, err := copyFile(path.Join("build", target), target)
			if err != nil {
				fmt.Print(fmt.Errorf("failed to copy file: %w", err))
				os.Exit(1)
			}
		}
	}

	f, err := os.Create("config.json")
	if err != nil {
		fmt.Errorf("failed to create file: %w", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(art); err != nil {
		fmt.Errorf("failed to encode json: %w", err)
		os.Exit(1)
	}

	// change below array to play with your own cmd args:
	args := []string{"push", "--config", "config.json", "ghcr.io/mxplusb/pleiades:latest"}
	args = append(args, buildArtifacts...)
	cmd := root.New()
	cmd.SetArgs(args)
	err = cmd.Execute()
	if err != nil {
		fmt.Errorf("Failed to execute : %w", err)
		os.Exit(1)
	}

	for _, f := range buildArtifacts {
		if err := os.Remove(f); err != nil {
			fmt.Errorf("failed to remove file: %w", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}

// https://stackoverflow.com/a/67179604/4949938
func copyFile(in, out string) (int64, error) {
	i, e := os.Open(in)
	if e != nil { return 0, e }
	defer i.Close()
	o, e := os.Create(out)
	if e != nil { return 0, e }
	defer o.Close()
	return o.ReadFrom(i)
}

type Artifact struct {
	SchemaVersion int         `json:"schemaVersion"`
	Config        Config      `json:"config,omitempty"`
	MediaType     string      `json:"mediaType,omitempty"`
	Layers        []Layer     `json:"layers,omitempty"`
	Annotations   Annotations `json:"annotations,omitempty"`
}

type Annotations struct {
	OrgOpenContainersImageTitle       string `json:"org.opencontainers.image.title"`
	OrgOpencontainersImageSource      string `json:"org.opencontainers.image.source"`
	OrgOpencontainersImageDescription string `json:"org.opencontainers.image.description"`
	OrgOpencontainersImageLicenses    string `json:"org.opencontainers.image.licenses"`
}

type Config struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int    `json:"size"`
}

type Layer struct {
	MediaType    string      `json:"mediaType"`
	Annotations  Annotations `json:"annotations,omitempty"`
	Architecture string      `json:"architecture,omitempty"`
	Os           string      `json:"os,omitempty"`
}
