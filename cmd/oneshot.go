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
	"fmt"
	"os"
	"syscall"

	"github.com/JAORMX/selinuxd/pkg/daemon"
	"github.com/JAORMX/selinuxd/pkg/semodule/semanage"
	"github.com/spf13/cobra"
)

// oneshotCmd represents the oneshot command
var oneshotCmd = &cobra.Command{
	Use:   "oneshot",
	Short: "install SELinux policies in the designated directory",
	Long:  `This does a one-shot installation of SELinux policies.`,
	Run:   oneshotCmdFunc,
}

// nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(oneshotCmd)
}

func oneshotCmdFunc(rootCmd *cobra.Command, _ []string) {
	logger, err := getLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		syscall.Exit(1)
	}

	sh, err := semanage.NewSemanageHandler(logger)
	if err != nil {
		logger.Error(err, "Creating semanage handler")
	}
	defer sh.Close()

	policyops := make(chan daemon.PolicyAction)

	logger.Info("Running oneshot command")

	go func() {
		if err := daemon.InstallPoliciesInDir(defaultModulePath, policyops); err != nil {
			logger.Error(err, "Installing policies in module directory")
		}
		close(policyops)
	}()

	daemon.InstallPolicies(defaultModulePath, sh, policyops, logger)

	logger.Info("Done installing policies in directory")
}