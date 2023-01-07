/*
 * Copyright (c) 2023 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

//go:build mage

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/cockroachdb/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/magefile/mage/mg"
)

type (
	Version mg.Namespace
	buildType string
)

const (
	devBuild buildType = "dev"
	alphaBuild buildType = "alpha"
	betaBuild buildType = "beta"
	prodBuild buildType = ""
)

func newVersion() (*semver.Version, error) {
	var buildNum int
	var err error
	if isCi() {
		buildNum, err = getBuildNumber()
	} else {
		buildNum = 0
	}
	now := time.Now().UTC()

	version := semver.New(uint64(now.Year()), uint64(now.Month()), 0, "dev", fmt.Sprintf("build%d", buildNum))

	return version, err
}

// generate a new version from scratch with a potential build number
func (Version) New() error {
	version, err := newVersion()
	fmt.Println(version.String())
	return err
}

// generate the next version number for teamcity
func (Version) Buildset() error {

	// the version is already set, so we can just bump the release channel
	existingVersion := os.Getenv("BUILD_TAG")
	if existingVersion != "nil" {
		currentVersion, err := semver.NewVersion(existingVersion)
		if err != nil {
			return errors.Wrap(err, "failed to parse existing version")
		}

		switch buildType(currentVersion.Prerelease()) {
		case devBuild:
			updatedVersion, err := currentVersion.SetPrerelease(string(alphaBuild))
			if err != nil {
				return errors.Wrap(err, "failed to set prerelease")
			}
			tellTeamCity(&updatedVersion)
			return nil
		case alphaBuild:
			updatedVersion, err := currentVersion.SetPrerelease(string(betaBuild))
			if err != nil {
				return errors.Wrap(err, "failed to set prerelease")
			}
			tellTeamCity(&updatedVersion)
			return nil
		case betaBuild:
			//
			updatedVersion, err := currentVersion.SetPrerelease(string(prodBuild))
			if err != nil {
				return errors.Wrap(err, "failed to clear prerelease")
			}
			finalVersion, err := updatedVersion.SetMetadata("")
			if err != nil {
				return errors.Wrap(err, "failed to clear metadata")
			}
			tellTeamCity(&finalVersion)
			return nil
		default:
			break
		}
	}

	localRepo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Printf("error: %s", err)
	}

	tagIter, err := localRepo.Tags()
	if err != nil {
		return err
	}

	// we want to ensure we have the latest version as git sees it.
	var mostRecentVersion *semver.Version
	err = tagIter.ForEach(func(reference *plumbing.Reference) error {
		val := strings.Split(reference.Name().String(), "/")

		parsedVersion, err := semver.NewVersion(val[len(val)-1])
		if err != nil {
			return errors.Wrapf(err, "failed to parse new version %s", val)
		}

		// since it's nil, go ahead and move along since there's nothing to compare to
		if mostRecentVersion == nil {
			mostRecentVersion = parsedVersion
			return nil
		}

		if mostRecentVersion.LessThan(parsedVersion) {
			mostRecentVersion = parsedVersion
			return nil
		}

		return err
	})

	currentVersion, err := newVersion()
	if err != nil {
		return err
	}

	finalVersion := *currentVersion
	// if it's the same month, just bump the patch version
	if currentVersion.Major() == mostRecentVersion.Major() && currentVersion.Minor() == mostRecentVersion.Minor() {
		targetPatch := mostRecentVersion.Patch()+1
		for finalVersion.Patch() < targetPatch {
			finalVersion = finalVersion.IncPatch()
		}
	}
	finalVersion, _ = finalVersion.SetMetadata(currentVersion.Metadata())
	finalVersion, _ = finalVersion.SetPrerelease(currentVersion.Prerelease())

	if isCi() {
		tellTeamCity(&finalVersion)
	}

	return err
}

func tellTeamCity(v *semver.Version) {
	// set the build and env parameters
	fmt.Printf("##teamcity[setParameter name='build.BUILD_TAG' value='%s']\n", v.String())
	fmt.Printf("##teamcity[setParameter name='env.BUILD_TAG' value='%s']\n", v.String())

	// clear the metadata for docker
	f, _ := v.SetMetadata("")
	fmt.Printf("##teamcity[setParameter name='env.DOCKER_BUILD_TAG' value='%s']\n", f.String())
}

func getBuildNumber() (int, error) {
	val := os.Getenv("BUILD_NUMBER")
	return strconv.Atoi(val)
}

func isCi() bool {
	val := os.Getenv("BUILD_NUMBER")
	if len(val) > 0 {
		return true
	}
	return false
}
