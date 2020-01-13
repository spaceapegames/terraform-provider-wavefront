package wavefront_plugin

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"sort"
	"testing"
)

func TestAccWavefrontUserGroup_BasicUserGroup(t *testing.T) {
	var record wavefront.UserGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontUserGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontUserGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontUserGroupExists("wavefront_user_group.basic", &record),
					testAccCheckWavefrontUserGroupAttributes(&record, []string{"alerts_management", "events_management"}),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_user_group.basic", "name", "Basic User Group"),
					resource.TestCheckResourceAttr(
						"wavefront_user_group.basic", "description", "Basic User Group for Unit Tests"),
					resource.TestCheckResourceAttr("wavefront_user_group.basic", "permissions.#", "2"),
				),
			},
			{
				Config: testAccCheckWavefrontUserGroup_changed(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontUserGroupExists("wavefront_user_group.basic", &record),
					testAccCheckWavefrontUserGroupAttributes(&record, []string{"alerts_management", "agent_management"}),

					// Check against the state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_user_group.basic", "name", "Basic User Groups"),
					resource.TestCheckResourceAttr(
						"wavefront_user_group.basic", "description", "Basic User Groups for Unit Tests"),
					resource.TestCheckResourceAttr("wavefront_user_group.basic", "permissions.#", "2"),
				),
			},
		},
	})
}

func testAccCheckWavefrontUserGroupExists(n string, userGroup *wavefront.UserGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		userGroups := testAccProvider.Meta().(*wavefrontClient).client.UserGroups()

		results, err := userGroups.Find(
			[]*wavefront.SearchCondition{
				{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("Error finding Wavefront User Group %s", err)
		}
		// resource has been deleted out of band. So unset ID
		if len(results) != 1 {
			return fmt.Errorf("No User Groups Found")
		}
		if *results[0].ID != rs.Primary.ID {
			return fmt.Errorf("User Group not found")
		}

		*userGroup = *results[0]

		return nil
	}
}

func testAccCheckWavefrontUserGroupDestroy(s *terraform.State) error {

	userGroups := testAccProvider.Meta().(*wavefrontClient).client.UserGroups()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "wavefront_user_group" {
			continue
		}

		results, err := userGroups.Find(
			[]*wavefront.SearchCondition{
				{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("error finding Wavefront User Group. %s", err)
		}
		if len(results) > 0 {
			return fmt.Errorf("user group still exists")
		}
	}

	return nil
}

func testAccCheckWavefrontUserGroupAttributes(userGroup *wavefront.UserGroup, permissions []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, v := range permissions {
			sort.Strings(userGroup.Permissions)
			if sort.SearchStrings(userGroup.Permissions, v) == len(userGroup.Permissions) {
				return fmt.Errorf("permission not found or present on User Group. %s", v)
			}
		}

		if !(userGroup.Name == "Basic User Group" || userGroup.Name == "Basic User Groups") {
			return fmt.Errorf("unexpected User Group name encountered. %s", userGroup.Name)
		}

		if !(userGroup.Description == "Basic User Group for Unit Tests" ||
			userGroup.Description == "Basic User Groups for Unit Tests") {
			return fmt.Errorf("unexpected User Group description encountered. %s", userGroup.Description)
		}

		return nil
	}
}

func testAccCheckWavefrontUserGroup_basic() string {
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

func testAccCheckWavefrontUserGroup_changed() string {
	return fmt.Sprintf(`
resource "wavefront_user_group" "basic" {
  name        = "Basic User Groups"
  description = "Basic User Groups for Unit Tests"
  permissions = [
    "alerts_management",
	"agent_management",
  ]
}
`)
}
