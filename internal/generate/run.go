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

	"github.com/googleapis/generator/internal/gitrepo"
)

func Run(ctx context.Context, arg ...string) error {
	cfg := &config{}
	cfg, err := parseFlags(cfg, arg)
	if err != nil {
		return err
	}

	for _, repo := range []struct {
		dir, url string
	}{
		{
			"/tmp/generator-googleapis",
			"https://github.com/googleapis/googleapis",
		},
		{
			"/tmp/generator-google-cloud-dotnet",
			"https://github.com/googleapis/google-cloud-dotnet",
		},
	} {
		slog.Info(fmt.Sprintf("Cloning %q to %q", repo.url, repo.dir))
		_, err := gitrepo.CloneOrOpen(ctx, repo.dir, repo.url)
		if err != nil {
			return err
		}
	}
	args := []string{
		"run", "-v", "tmp-apis:/tmp/generator-googleapis",
		"-v", "tmp-dotnet:/output",
		"picard", "--command=update",
		"--api-root=/tmp/generator-googleapis",
		"--api=$api",
		"--output=/output",
	}
	if err := runCommand("docker", args...); err != nil {
		return err
	}
	return nil
}

func runCommand(c string, args ...string) error {
	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	slog.Info(strings.Repeat("-", 80))
	slog.Info(cmd.String())
	return cmd.Run()
}
