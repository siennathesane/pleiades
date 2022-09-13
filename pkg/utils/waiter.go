/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package utils

import (
	"os"
	"time"
)

func Wait(wait time.Duration) {
	val := os.Getenv("CI_JOB_ID")
	if val != "" {
		length := wait.Milliseconds()
		longerTime := float64(length) * 2.25
		time.Sleep(time.Duration(longerTime) * time.Millisecond)
		return
	}
	time.Sleep(wait)
}

func Timeout(wait time.Duration) time.Duration {
	val := os.Getenv("CI_JOB_ID")
	if val != "" {
		length := wait.Milliseconds()
		longerTime := float64(length) * 2.25
		return time.Duration(longerTime) * time.Millisecond
	}
	return wait
}