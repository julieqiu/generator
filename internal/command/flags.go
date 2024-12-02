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

package command

import (
	"context"
	"flag"
	"fmt"
)

var Cmd = &Command{
	Run:   run,
	Short: `Generator is a tool for generating Google client libraries.`,
}

var Generate = &Command{
	Run:   runGenerate,
	Short: `generate generates client libraries`,
}

// A Command is an implementation of the generator command,
// like generator generate or generator lint.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(ctx context.Context, cmd *Command, args []string) error

	// Short is a short description of the command.
	Short string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

var (
	flagLanguage string
	flagAPI      string
)

func usage(flags *flag.FlagSet) {
	flags.Usage = func() {
		fmt.Fprint(flags.Output(), `Generator generates client libraries for Google APIs.

Usage:

  generator [commands] [flags]

Flags:

`)
		flags.PrintDefaults()
		fmt.Fprintf(flags.Output(), "\n\n")
	}
}

func addAPIFlag(flags *flag.FlagSet) {
	flags.Func("api", "name of API inside googleapis", func(l string) error {
		return nil
	})
}

func addLanguageFlag(flags *flag.FlagSet) {
	flags.Func("language", "the generated language", func(l string) error {
		if l != "cpp" {
			return fmt.Errorf("not implemented")
		}
		if _, ok := languages[l]; !ok {
			return fmt.Errorf("invalid -language flag specified: %q", l)
		}

		flagLanguage = l
		return nil
	})
}

var languages = map[string]bool{
	"cpp":    true,
	"dotnet": true,
	"go":     true,
	"java":   true,
	"node":   true,
	"php":    true,
	"python": true,
	"ruby":   true,
	"rust":   true,
}
