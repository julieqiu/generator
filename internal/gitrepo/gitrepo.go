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

// Package gitrepo provides operations on git repos.
package gitrepo

import (
	"context"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Repo struct {
	repo *git.Repository
}

// CloneOrOpen clones repoPath if it is an HTTP(S) URL, or opens it from the
// local disk otherwise.
func CloneOrOpen(ctx context.Context, repoPath string) (*Repo, error) {
	if strings.HasPrefix(repoPath, "http://") || strings.HasPrefix(repoPath, "https://") {
		return Clone(ctx, repoPath)
	}
	return Open(ctx, repoPath)
}

// Clone returns a bare repo by cloning the repo at repoURL.
func Clone(ctx context.Context, repoURL string) (*Repo, error) {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.HEAD,
		SingleBranch:  true,
		Depth:         1,
		Tags:          git.NoTags,
	})
	if err != nil {
		return nil, err
	}
	return &Repo{repo: repo}, nil
}

// Open returns a repo by opening the repo at the local path dirpath.
func Open(ctx context.Context, dirpath string) (*Repo, error) {
	repo, err := git.PlainOpen(dirpath)
	if err != nil {
		return nil, err
	}
	return &Repo{repo: repo}, nil
}
