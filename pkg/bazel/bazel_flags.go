/*
 * Copyright 2022 Aspect Build Systems, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bazel

import (
	"fmt"
	"strings"

	"aspect.build/cli/pkg/aspect/root/flags"
	"github.com/spf13/cobra"
)

var (
	// Bazel flags specified here will be shown when running "aspect help".
	// By default bazel flags are hidden.
	documentedBazelFlags = []string{
		"keep_going",
		"expunge",
		"expunge_async",
		"show_make_env",
	}
)

var (
	// Bazel flags that expand to other flags. These are boolean flags that are not no-able. Currently
	// there is no way to detect these so we have to keep list up-to-date manually with the union of
	// these flags across all Bazel versions we support.
	// These were gathered by searching https://bazel.build/reference/command-line-reference for "Expands to:"
	expandoFlags = map[string]struct{}{
		"debug_app":                     {},
		"experimental_persistent_javac": {},
		"experimental_spawn_scheduler":  {},
		"expunge_async":                 {},
		"host_jvm_debug":                {},
		"java_debug":                    {},
		"long":                          {},
		"noincompatible_genquery_use_graphless_query": {},
		"noorder_results":                                 {},
		"null":                                            {},
		"order_results":                                   {},
		"persistent_android_dex_desugar":                  {},
		"persistent_android_resource_processor":           {},
		"persistent_multiplex_android_dex_desugar":        {},
		"persistent_multiplex_android_resource_processor": {},
		"persistent_multiplex_android_tools":              {},
		"remote_download_minimal":                         {},
		"remote_download_toplevel":                        {},
		"short":                                           {},
		"start_app":                                       {},
	}
)

// AddBazelFlags will process the configured cobra commands and add bazel
// flags to those commands.
func AddBazelFlags(cmd *cobra.Command) error {
	subCommands := make(map[string]*cobra.Command)

	for _, subCmd := range cmd.Commands() {
		subCmdName := strings.SplitN(subCmd.Use, " ", 2)[0]
		subCommands[subCmdName] = subCmd
	}

	bzl, err := FindFromWd()
	if err != nil {
		// We cannot run Bazel, but this just means we have no flags to add.
		// This will be the case when running aspect help from outside a workspace, for example.
		// If Bazel is really needed for the current command, an error will be handled somewhere else.
		return nil
	}
	bzlFlags, err := bzl.Flags()
	if err != nil {
		return fmt.Errorf("unable to determine available bazel flags: %w", err)
	}

	for flagName := range bzlFlags {
		flag := bzlFlags[flagName]
		flagAbbreviation := flag.GetAbbreviation()
		flagDoc := flag.GetDocumentation()

		for _, command := range flag.Commands {
			if command == "startup" {
				if flag.GetHasNegativeFlag() {
					flags.RegisterNoableBool(cmd.PersistentFlags(), flagName, false, flagDoc)
					markFlagAsHidden(cmd, flagName)
					markFlagAsHidden(cmd, flags.NoFlagName(flagName))
				} else if flag.GetAllowsMultiple() {
					var key = flags.MultiString{}
					cmd.PersistentFlags().VarP(&key, flagName, flagAbbreviation, flagDoc)
					markFlagAsHidden(cmd, flagName)
				} else {
					_, isExpando := expandoFlags[flagName]
					if isExpando {
						cmd.PersistentFlags().BoolP(flagName, flagAbbreviation, false, flagDoc)
					} else {
						cmd.PersistentFlags().StringP(flagName, flagAbbreviation, "", flagDoc)
					}
					markFlagAsHidden(cmd, flagName)
				}
			}
			if subcommand, ok := subCommands[command]; ok {
				subcommand.DisableFlagParsing = true // only want to disable flag parsing on actual bazel verbs
				if flag.GetHasNegativeFlag() {
					flags.RegisterNoableBoolP(subcommand.Flags(), flagName, flagAbbreviation, false, flagDoc)
					markFlagAsHidden(subcommand, flagName)
					markFlagAsHidden(subcommand, flags.NoFlagName(flagName))
				} else if flag.GetAllowsMultiple() {
					var key = flags.MultiString{}
					subcommand.Flags().VarP(&key, flagName, flagAbbreviation, flagDoc)
					markFlagAsHidden(subcommand, flagName)
				} else {
					_, isExpando := expandoFlags[flagName]
					if isExpando {
						subcommand.Flags().BoolP(flagName, flagAbbreviation, false, flagDoc)
					} else {
						subcommand.Flags().StringP(flagName, flagAbbreviation, "", flagDoc)
					}
					markFlagAsHidden(subcommand, flagName)
				}
			}
		}
	}

	return nil
}

func markFlagAsHidden(cmd *cobra.Command, flag string) {
	for _, documentedFlag := range documentedBazelFlags {
		if documentedFlag == flag {
			return
		}
	}

	cmd.Flags().MarkHidden(flag)
	cmd.PersistentFlags().MarkHidden(flag)
}
