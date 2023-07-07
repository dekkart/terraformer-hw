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

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/httphandler"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	ecsModel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	ecsRegion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/region"
)

type ECSGenerator struct {
	HuaweiCloudService
}

func (g *ECSGenerator) loadInstances() (*[]ecsModel.ServerDetail, error) {
	AK := g.Args["AK"].(string)
	SK := g.Args["SK"].(string)
	Region := g.Args["Region"].(string)
	fmt.Println("Load instances")
	// Добавить цикл для count если серверов много
	auth := basic.NewCredentialsBuilder().
		// Access Key and Secret Access Key used for authentication
		WithAk(AK).
		WithSk(SK).
		// If ProjectId is not filled in, the SDK will automatically call the IAM service to query the project id corresponding to the region.
		//	WithProjectId("ru-moscow-1").
		// Configure the SDK built-in IAM service endpoint, default is https://iam.myhuaweicloud.com
		WithIamEndpointOverride("https://iam.ru-moscow-1.hc.sbercloud.ru/v3").
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
	client := ecs.NewEcsClient(
		ecs.EcsClientBuilder().
			// Enpoint will be added automaticly from ~/.huaweicloud/regions.yaml
			//WithEndpoint("https://ecs.ru-moscow-1.hc.sbercloud.ru").
			// Configure region, it will cause a panic if the region does not exist
			WithRegion(ecsRegion.ValueOf(Region)).
			// Configure authentication
			WithCredential(auth).

			// Configure HTTP
			WithHttpConfig(httpConfig).
			//WithHttpConfig(config.DefaultHttpConfig()).
			Build())

	// Create a request
	request := &ecsModel.ListServersDetailsRequest{}
	// Configure the number of records (Servers) on each page
	limit := int32(100)
	request.Limit = &limit

	// Send the request and get the response
	response, err := client.ListServersDetails(request)
	// Handle error and print response
	if err == nil {
		//fmt.Printf("%+v\n", response.Count)
		fmt.Println(*response.Count)
	} else {
		fmt.Println(err)
	}
	data := *response.Servers
	for i, s := range data {
		fmt.Println(i, s.Id)
		//instances = append(instances, &s...)
	}

	return response.Servers, nil

}

// InitResources Gets the list of all ECS instance ids and generates resources
func (g *ECSGenerator) InitResources() error {
	fmt.Println("InitResources")
	result, err := g.loadInstances()
	if err != nil {
		return err
	}
	g.Resources = g.createResources(result)

	return nil
}

// Генерация ресурсов в HCL
func (g *ECSGenerator) createResources(instances *[]ecsModel.ServerDetail) []terraformutils.Resource {
	fmt.Println("createResources")
	var resources []terraformutils.Resource
	for _, instance := range *instances {
		resources = append(resources, terraformutils.NewSimpleResource(
			instance.Id,
			instance.Name,
			"sbercloud_compute_instance",
			"sbercloud",
			[]string{}))
	}
	return resources
}
