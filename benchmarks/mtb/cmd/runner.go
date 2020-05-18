/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type args struct {
	values []string
}

func (a *args) addIfNonEmpty(flagName, value string) {
	if value != "" {
		a.values = append(a.values, fmt.Sprintf("--%s=%s", flagName, value))
	}
}

func (a *args) addBool(flagName string, value bool) {
	a.values = append(a.values, fmt.Sprintf("--%s=%t", flagName, value))
}

func (a *args) addInt(flagName string, value int) {
	a.values = append(a.values, fmt.Sprintf("--%s=%d", flagName, value))
}

func runCommand(a *args) error {

	ginkgoPath, err := findBinary("ginkgo")
	if err != nil {
		return err
	}

	e2eTest, err := findBinary("tests.test")
	if err != nil {
		return err
	}

	a.values = append(a.values, []string{
		e2eTest,
	}...)

	var extraArgs []string

	ginkgoArgs := append(a.values, extraArgs...)
	cmd := exec.Command(ginkgoPath, ginkgoArgs...)

	var stdoutBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)

	errNew := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", errNew)
	}

	return nil
}

func findBinary(name string) (string, error) {
	gopath := os.Getenv("GOPATH")

	if gopath == "" {
		return "", fmt.Errorf("$GOPATH not set in the env variables")
	}

	locations := []string{
		filepath.Join(gopath, "bin", name),
		filepath.Join(gopath, "src", "sigs.k8s.io", "multi-tenancy", "benchmarks", name),
		filepath.Join(gopath, "src", "sigs.k8s.io", "multi-tenancy", "benchmarks", "e2e", "tests", name),
	}

	newestLocation := ""
	var newestModTime time.Time
	for _, loc := range locations {
		stat, err := os.Stat(loc)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", fmt.Errorf("error from stat %s: %v", loc, err)
		}
		if newestLocation == "" || stat.ModTime().After(newestModTime) {
			newestModTime = stat.ModTime()
			newestLocation = loc
		}
	}

	if newestLocation == "" {
		log.Printf("could not find %s, looked in %s", name, locations)
		return "", fmt.Errorf("could not find %s", name)
	}

	return newestLocation, nil
}
