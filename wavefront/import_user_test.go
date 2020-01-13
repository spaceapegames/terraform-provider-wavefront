package wavefront_plugin

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccUser_importBasic(t *testing.T) {
	resourceName := "wavefront_user.basic"
	var record wavefront.User

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontUserImporter_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontUserExists("wavefront_user.basic", &record),
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

func testAccCheckWavefrontUserImporter_basic() string {
	return fmt.Sprintf(`
resource "wavefront_user" "basic" {
	email  = "test+tftesting@example.com"
	groups = [
		"agent_management",
		"alerts_management",
	]
}`)
}
