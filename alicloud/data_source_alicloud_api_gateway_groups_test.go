package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudApigatewayGroupsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudApiGatewayGroupDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_api_gateway_groups.data_apigatway_groups"),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_groups.data_apigatway_groups", "groups.0.name", "tf_testAccGroupDataSource"),
					resource.TestCheckResourceAttr("data.alicloud_api_gateway_groups.data_apigatway_groups", "groups.0.description", "tf_testAcc api gateway description"),
				),
			},
		},
	})
}

const testAccCheckAlicloudApiGatewayGroupDataSource = `

variable "apigateway_group_name_test" {
  default = "tf_testAccGroupDataSource"
}

variable "apigateway_group_description_test" {
  default = "tf_testAcc api gateway description"
}

resource "alicloud_api_gateway_group" "apiGroupTest" {
  name = "${var.apigateway_group_name_test}"
  description = "${var.apigateway_group_description_test}"
}

data "alicloud_api_gateway_groups" "data_apigatway_groups"{
  name_regex = "${alicloud_api_gateway_group.apiGroupTest.name}"
}

`
