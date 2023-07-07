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

package huaweicloud

import (
	"errors"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/zclconf/go-cty/cty"
)

// AK & SK will be got from env on init()
const AK = "AK"
const SK = "SK"

// Default region
const Region = "ru-moscow-1"

type HuaweiCloudProvider struct { //nolint
	terraformutils.Provider
	Region string
	AK     string
	SK     string
}

func (p *HuaweiCloudProvider) Init(args []string) error {
	if SBC_ACCESS_KEY, ok := os.LookupEnv("SBC_ACCESS_KEY"); ok {
		p.AK = SBC_ACCESS_KEY
		log.Println("SBC_ACCESS_KEY: ", p.AK)
	} else {
		log.Printf("%s not set\n", SBC_ACCESS_KEY)
	}
	//p.AK = os.Getenv("SBC_ACCESS_KEY")
	if SBC_SECRET_KEY, ok := os.LookupEnv("SBC_SECRET_KEY"); ok {
		p.SK = SBC_SECRET_KEY
		log.Println("SBC_SECRET_KEY: ", p.SK)
	} else {
		log.Printf("%s not set\n", SBC_SECRET_KEY)
	}
	//p.SK = os.Getenv("SBC_SECRET_KEY")
	p.Region = Region
	log.Println("Region: ", p.Region)
	if p.AK == "" || p.SK == "" {
		panic("AK or SK not set in ENV variables, please set up env")
	}

	log.Println("Provider.config: ", p.GetConfig())

	return nil
}

func (p *HuaweiCloudProvider) GetName() string {
	return "huaweicloud"
}

func (p *HuaweiCloudProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{}
}

func (HuaweiCloudProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{}
}

func (p *HuaweiCloudProvider) GetSupportedService() map[string]terraformutils.ServiceGenerator {
	return map[string]terraformutils.ServiceGenerator{
		"evs": &EVSGenerator{},
		// InstanceGenerator подтягивается автоматом из пакета:
		"ecs": &ECSGenerator{},
	}
}

// Configure provider:
func (p *HuaweiCloudProvider) GetConfig() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		//	"AccessKey": cty.StringVal(p.AK),
		//	"SecretKey": cty.StringVal(p.SK),
		"region": cty.StringVal(p.Region),
	})
}

func (p *HuaweiCloudProvider) InitService(serviceName string, verbose bool) error {
	var isSupported bool

	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New("Huawei Cloud: " + serviceName + " not supported service")
	}

	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetVerbose(verbose)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]interface{}{
		"AK":     p.AK,
		"SK":     p.SK,
		"Region": p.Region,
	})
	return nil
}
