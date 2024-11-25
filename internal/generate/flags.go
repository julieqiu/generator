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
	"flag"
	"fmt"
)

type config struct {
	cmd        string
	dir        string
	googleapis string
	language   string
	output     string
	patterns   []string
}

func parseFlags(args []string) (*config, error) {
	cfg := &config{}
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flags.StringVar(&cfg.dir, "C", "", "change to `dir` before running generator")
	flags.StringVar(&cfg.language, "language", "", "specify from cpp, csharp, go, java, node, php, python, ruby, rust")
	flags.StringVar(&cfg.googleapis, "googleapis", "/Users/julieqiu/code/googleapis/googleapis", "location of googleapis `dir`")
	flags.StringVar(&cfg.output, "output", "/tmp/cloudsdkgenerator", "`dir` to write generated client library output")

	// We don't want to print the whole usage message on each flags
	// error, so we set to a no-op and do the printing ourselves.
	flags.Usage = func() {}
	usage := func() {
		fmt.Fprint(flags.Output(), `Generator generates client libraries for Google APIs.

Usage:

  generator [command] [flags]

Commands:

  generate           Generate a client library using the provided files
  lint               Run linters to validate the API specification and configuration files

Flags:

`)
		flags.PrintDefaults()
		fmt.Fprintf(flags.Output(), "\n\n")
	}

	cfg.cmd = flag.Arg(0)
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			usage() // print usage only on help
		}
		return nil, err
	}
	if cfg.cmd == "generate" {
		if cfg.language == "" || cfg.googleapis == "" {
			usage() // print usage only on help
			return nil, fmt.Errorf("missing flags")
		}
	}
	return cfg, nil
}
