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
	"fmt"
	"log"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/httphandler"

	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"
	vpcModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
	vpcRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/region"
)

type SGGenerator struct {
	HuaweiCloudService
}

func (g *SGGenerator) loadSG() (*[]vpcModel.SecurityGroup, error) {
	AK := g.Args["AK"].(string)
	SK := g.Args["SK"].(string)
	Region := g.Args["Region"].(string)
	log.Println("Load Security Groups")
	// Добавить цикл для count если SG больше 1000
	auth := basic.NewCredentialsBuilder().
		// Access Key and Secret Access Key used for authentication
		WithAk(AK).
		WithSk(SK).
		// If ProjectId is not filled in, the SDK will automatically call the IAM service to query the project id corresponding to the region.
		//	WithProjectId("ru-moscow-1").
		// Configure the SDK built-in IAM service endpoint, default is https://iam.myhuaweicloud.com
		//WithIamEndpointOverride("https://iam.ru-moscow-1.hc.sbercloud.ru/v3").
		Build()

	// Use default configuration
	httpConfig := config.DefaultHttpConfig()
	// Configure whether to ignore the SSL certificate verification, default is false
	httpConfig.WithIgnoreSSLVerification(true)
	// Configure HTTP handler for debugging, do not use in production environment
	/*
		requestHandler := func(request http.Request) {
			fmt.Println(request)
		}
		responseHandler := func(response http.Response) {
			fmt.Println(response)
		}
	*/
	// httpHandler := httphandler.NewHttpHandler().AddRequestHandler(requestHandler).AddResponseHandler(responseHandler)
	httpHandler := httphandler.NewHttpHandler()
	httpConfig.WithHttpHandler(httpHandler)

	// Create a service client
	client := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			// Enpoint will be added automaticly from ~/.huaweicloud/regions.yaml
			//WithEndpoint("https://ecs.ru-moscow-1.hc.sbercloud.ru").
			// Configure region, it will cause a panic if the region does not exist
			WithRegion(vpcRegion.ValueOf(Region)).
			// Configure authentication
			WithCredential(auth).

			// Configure HTTP
			WithHttpConfig(httpConfig).
			//WithHttpConfig(config.DefaultHttpConfig()).
			Build())

	// Create a request
	request := &vpcModel.ListSecurityGroupsRequest{}
	// Configure the number of records (disks) on each page
	limit := int32(defaultPageSize)
	request.Limit = &limit

	// Send the request and get the response
	// List up to 1000 SG now, will change later
	response, err := client.ListSecurityGroups(request)
	// Handle error and print response
	if err == nil {
		//fmt.Printf("%+v\n", response.Count)
		log.Println(&response.PageInfo.CurrentCount)
	} else {
		fmt.Println(err)
	}
	data := *response.SecurityGroups
	for i, s := range data {
		log.Println(i, s.Id)
	}

	return response.SecurityGroups, nil

}

// InitResources Gets the list of all Security Groups ids and generates resources
func (g *SGGenerator) InitResources() error {
	log.Println("Init Security Groups")
	result, err := g.loadSG()
	if err != nil {
		return err
	}
	g.Resources = g.createResources(result)

	return nil
}

// Генерация ресурсов в HCL
func (g *SGGenerator) createResources(SecurityGroups *[]vpcModel.SecurityGroup) []terraformutils.Resource {
	log.Println("create security groups")
	var resources []terraformutils.Resource
	for _, sg := range *SecurityGroups {
		resources = append(resources, terraformutils.NewSimpleResource(
			sg.Id,
			sg.Name,
			"huaweicloud_networking_secgroup",
			"huaweicloud",
			[]string{}))
	}
	return resources
}
