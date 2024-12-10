// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package container

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

var errNotImplemented = errors.New("not implemented")

func Generate(ctx context.Context, apiRoot, apiPath, output, generatorInput string) error {
	return errNotImplemented
}

func Clean(ctx context.Context, repoRoot, apiPath string) error {
	return errNotImplemented
}

func Build(ctx context.Context, repoRoot, apiPath string) error {
	return errNotImplemented
}

func Configure(ctx context.Context, apiRoot, apiPath, generatorInput string) error {
	return errNotImplemented
}

const dotnetImageTag = "picard"

func runDocker(googleapisDir, languageDir, api string) error {
	args := []string{
		"run",
		"-v", fmt.Sprintf("%s:/apis", googleapisDir),
		"-v", fmt.Sprintf("%s:/output", languageDir),
		dotnetImageTag,
		"--command=update",
		"--api-root=/apis",
		fmt.Sprintf("--api=%s", api),
		"--output-root=/output",
	}
	return runCommand("docker", args...)
}

func runCommand(c string, args ...string) error {
	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	slog.Info(strings.Repeat("-", 80))
	slog.Info(cmd.String())
	slog.Info(strings.Repeat("-", 80))
	return cmd.Run()
}
