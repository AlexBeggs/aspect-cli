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

package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

const (
	BoolFlagTrue  = "true"
	BoolFlagFalse = "false"
	BoolFlagYes   = "yes"
	BoolFlagNo    = "no"
	BoolFlag1     = "1"
	BoolFlag0     = "0"
)

const (
	// The prefix that Bazel & Aspect use for negative flags such as --nohome_rc && --aspect:nohome_config
	NoFlagPrefix = "no"
)

// Prefixes a flag name with "no" to determine the Bazel negative flag from a flag name. Also takes
// into consideration the 'aspect:' prefix so that a flag such as 'aspect:foo' get a negative flag
// name of 'aspect:nofoo'. For example, `nohome_rc` is the negative of `home_rc` and
// `aspect:nohome_config` is the negative of 'aspect:home_config'
func NoFlagName(name string) string {
	if strings.HasPrefix(name, AspectFlagPrefix) {
		return fmt.Sprintf("%s%s%s", AspectFlagPrefix, NoFlagPrefix, name[len(AspectFlagPrefix):])
	} else {
		return NoFlagPrefix + name
	}
}

// RegisterNoableBool registers a boolean flag that supports Bazel option parsing.
func RegisterNoableBool(flags *pflag.FlagSet, name string, value bool, usage string) *bool {
	return RegisterNoableBoolP(flags, name, "", value, usage)
}

// RegisterNoableBoolP registers a boolean flag that supports Bazel option parsing with a shorthand.
// https://bazel.build/reference/command-line-reference#option-syntax
//
// By default the noable variants are -[no]foo which mimics how Bazel handles boolean flags.
// For Aspect CLI flags, which start with 'aspect:', the noable variants are --aspect:[no]foo.
//
// This implementation normalizes any user-provided values before processing. Hence,
// `--foo=yes` is the same as `--foo=YES`.
//
// Examples:
//
//	--foo
//	--nofoo
//	--foo=yes
//	--foo=no
//	--foo=1
//	--foo=0
func RegisterNoableBoolP(
	flags *pflag.FlagSet,
	name string,
	shorthand string,
	defaultValue bool,
	usage string) *bool {

	result := defaultValue

	trueNB := &noableBool{value: &result, valueWhenTrue: true}
	flag := &pflag.Flag{
		Name:      name,
		Shorthand: shorthand,
		Usage:     usage,
		Value:     trueNB,
		DefValue:  trueNB.String(),
		// The value that will be passed to Set() if no other values are specified.
		NoOptDefVal: BoolFlagTrue,
	}
	flags.AddFlag(flag)

	falseNB := &noableBool{value: &result, valueWhenTrue: false}
	noFlag := &pflag.Flag{
		Name:      NoFlagName(name),
		Shorthand: "",
		Usage:     usage,
		Value:     falseNB,
		DefValue:  falseNB.String(),
		// The value that will be passed to Set() if no other values are specified.
		NoOptDefVal: BoolFlagTrue,
	}
	flags.AddFlag(noFlag)

	return &result
}

func boolStr(value bool) string {
	return fmt.Sprintf("%t", value)
}

type noableBool struct {
	// The address of the actual value.
	value *bool
	// The value that should be set when the flag is set to true.
	valueWhenTrue bool
}

func (nb *noableBool) Type() string {
	return "bool"
}

func (nb *noableBool) String() string {
	return boolStr(*nb.value)
}

func (nb *noableBool) Set(value string) error {
	normalizedValue := strings.ToLower(value)
	var inValue bool
	// If this is the noXXX flag
	if !nb.valueWhenTrue {
		// The only allowed value for a noXXX flag is true
		if normalizedValue == BoolFlagTrue {
			inValue = true
		} else {
			return fmt.Errorf("invalid no flag value '%s'", value)
		}
	} else {
		switch normalizedValue {
		case BoolFlagTrue, BoolFlagYes, BoolFlag1:
			inValue = true
		case BoolFlagFalse, BoolFlagNo, BoolFlag0:
			inValue = false
		default:
			return fmt.Errorf("invalid bool value '%s'", value)
		}
	}

	// Set the actual boolean value
	if inValue {
		*nb.value = nb.valueWhenTrue
	} else {
		*nb.value = !nb.valueWhenTrue
	}
	return nil
}
