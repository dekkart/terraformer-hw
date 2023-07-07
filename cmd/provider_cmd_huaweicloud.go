// Copyright 2019 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	hw_terraforming "github.com/dekkart/terraformer-hw/providers/huaweicloud"
	"github.com/spf13/cobra"
)

func newCmdHuaweiCloudImporter(options ImportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "huaweicloud",
		Short: "Import current state to Terraform configuration from Huawei Cloud",
		Long:  "Import current state to Terraform configuration from Huawei Cloud",
		RunE: func(cmd *cobra.Command, args []string) error {

			provider := newHuaweiCloudProvider()
			log.Println("Provider name: ", provider.GetName())
			log.Println("Provider data: ", provider.GetConfig())
			err := Import(provider, options, []string{})
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.AddCommand(listCmd(newHuaweiCloudProvider()))
	baseProviderFlags(cmd.PersistentFlags(), &options, "ecs, evs", "")

	return cmd
}

func newHuaweiCloudProvider() terraformutils.ProviderGenerator {
	return &hw_terraforming.HuaweiCloudProvider{}
}
