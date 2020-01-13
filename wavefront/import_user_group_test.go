package wavefront_plugin

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccUserGroup_importBasic(t *testing.T) {
	resourceName := "wavefront_user_group.basic"
	var record wavefront.UserGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontUserGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontUserGroupImporter_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontUserGroupExists("wavefront_user_group.basic", &record),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckWavefrontUserGroupImporter_basic() string {
	return fmt.Sprintf(`
resource "wavefront_user_group" "basic" {
  name        = "Basic User Group"
  description = "Basic User Group for Unit Tests"
  permissions = [
    "alerts_management",
	"events_management"
  ]
}
`)
}
