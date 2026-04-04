// Copyright 2021 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

// Package pprompt provides a facility to prompt a user for a password
// securely (i.e. without echoing the password) from an interactive
// terminal.
//
// This is a separate package to ensure that the CLI code that uses
// this does not indirectly depend on other features from the
// 'security' package.
package pprompt

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// PromptForPassword prompts for a password.
// This is meant to be used when using a password.
func PromptForPassword() (string, error) {
	fmt.Print("Enter password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	// Make sure stdout moves on to the next line.
	fmt.Print("\n")

	return string(password), nil
}
