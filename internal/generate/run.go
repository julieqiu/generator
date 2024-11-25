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

package generate

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func Run(ctx context.Context, arg ...string) error {
	cfg, err := parseFlags(arg)
	if err != nil {
		return err
	}
	return generate(cfg)
}

func generate(cfg *config) error {
	fmt.Println(cfg.googleapis)
	targets, err := fetchTargets(cfg.googleapis, cfg.language)
	if err != nil {
		return err
	}
	for _, target := range targets {
		if err := bazelBuild(cfg.googleapis, target); err != nil {
			return err
		}
		if err := untar(cfg.googleapis, cfg.output, target); err != nil {
			return err
		}
	}
	if err := runGoPostprocessor(cfg.googleapis, cfg.output); err != nil {
		return err
	}
	return nil
}

func fetchTargets(googleapisDir, language string) ([]string, error) {
	cmd := exec.Command("bazelisk", "query", `filter("-(go)$", kind("rule", //...:*))`)
	cmd.Dir = googleapisDir
	cmd.Stderr = os.Stderr

	slog.Info(strings.Repeat("-", 80))
	slog.Info(cmd.Dir, googleapisDir)
	slog.Info(cmd.String())
	slog.Info(strings.Repeat("-", 80))

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	targets := strings.Fields(string(output))
	for _, t := range targets {
		slog.Info(t)
	}

	// Confirm that bazel can fetch remote build dependencies before building
	// with -k.  Otherwise, we can't distinguish a build failure due to a bad proto
	// vs. a build failure due to transient network issue.
	return targets, runCommand(googleapisDir, "bazelisk", append([]string{"fetch"}, targets...)...)
}

func bazelBuild(googleapisDir, target string) error {
	// Invoke bazel build. Limiting job count helps to avoid memory error b/376777535.
	return runCommand(googleapisDir, "bazelisk", "build", "--jobs=8", "-k", target)
}

func untar(googleapisDir, outputDir, target string) error {
	parts := strings.SplitN(target, ":", 2)
	parts[0] = strings.TrimPrefix(parts[0], "//")
	tarFile := fmt.Sprintf("%s/bazel-bin/%s/%s.tar.gz", googleapisDir, parts[0], parts[1])
	return runCommand(outputDir, "tar", "-xf", tarFile)
}

func runGoPostprocessor(googleapisDir, outputDir string) error {
	if _, err := os.Create(fmt.Sprintf("%s/cloud.google.com/go/internal/.repo-metadata-full.json", outputDir)); err != nil {
		return err
	}
	return runCommand(".", "go", "run", "./postprocessor",
		"--client-root", fmt.Sprintf("%s/cloud.google.com/go", outputDir),
		"--googleapis-dir", googleapisDir,
		"--dirs", fmt.Sprintf("%s/cloud.google.com/go", outputDir),
		"--pr-file", "prfile.txt",
	)
}

func runCommand(dir, c string, args ...string) error {
	cmd := exec.Command(c, args...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	slog.Info(strings.Repeat("-", 80))
	slog.Info(cmd.Dir)
	slog.Info(cmd.String())
	slog.Info(strings.Repeat("-", 80))

	return cmd.Run()
}
