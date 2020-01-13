package wavefront_plugin

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"sort"
	"testing"
)

func TestAccWavefrontUser_BasicUser(t *testing.T) {
	var record wavefront.User

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontUser_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontUserExists("wavefront_user.basic", &record),
					testAccCheckWavefrontUserAttributes(&record, []string{"agent_management", "alerts_management"}, []string{}),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_user.basic", "id", "test+tftesting@example.com"),
					resource.TestCheckResourceAttr(
						"wavefront_user.basic", "groups.#", "2"),
				),
			},
		},
	})
}

func TestAccWavefrontUser_BasicUserChangeGroups(t *testing.T) {
	var record wavefront.User

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontUser_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontUserExists("wavefront_user.basic", &record),
					testAccCheckWavefrontUserAttributes(&record, []string{"agent_management", "alerts_management"}, []string{}),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_user.basic", "id", "test+tftesting@example.com"),
					resource.TestCheckResourceAttr(
						"wavefront_user.basic", "groups.#", "2"),
				),
			},
			{
				Config: testAccCheckWavefrontUser_changeGroups(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontUserExists("wavefront_user.basic", &record),
					testAccCheckWavefrontUserAttributes(&record, []string{"agent_management", "events_management"}, []string{}),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_user.basic", "id", "test+tftesting@example.com"),
					resource.TestCheckResourceAttr(
						"wavefront_user.basic", "groups.#", "2"),
				),
			},
		},
	})
}

func TestAccWavefrontUser_BasicUserChangeEmail(t *testing.T) {
	var record wavefront.User

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontUser_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontUserExists("wavefront_user.basic", &record),
					testAccCheckWavefrontUserAttributes(&record, []string{"agent_management", "alerts_management"}, []string{}),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_user.basic", "id", "test+tftesting@example.com"),
					resource.TestCheckResourceAttr(
						"wavefront_user.basic", "groups.#", "2"),
				),
			},
			{
				Config: testAccCheckWavefrontUser_changeEmail(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontUserExists("wavefront_user.basic", &record),
					testAccCheckWavefrontUserAttributes(&record, []string{"agent_management", "alerts_management"}, []string{}),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_user.basic", "id", "test+tftesting2@example.com"),
					resource.TestCheckResourceAttr(
						"wavefront_user.basic", "groups.#", "2"),
				),
			},
		},
	})
}

func testAccCheckWavefrontUserDestroy(s *terraform.State) error {

	users := testAccProvider.Meta().(*wavefrontClient).client.Users()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "wavefront_user" {
			continue
		}

		results, err := users.Find(
			[]*wavefront.SearchCondition{
				{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("error finding Wavefront User. %s", err)
		}
		if len(results) > 0 {
			return fmt.Errorf("user still exists")
		}
	}

	return nil
}

func testAccCheckWavefrontUserExists(n string, user *wavefront.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		users := testAccProvider.Meta().(*wavefrontClient).client.Users()

		results, err := users.Find(
			[]*wavefront.SearchCondition{
				{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("Error finding Wavefront User %s", err)
		}
		// resource has been deleted out of band. So unset ID
		if len(results) != 1 {
			return fmt.Errorf("No Users Found")
		}
		if *results[0].ID != rs.Primary.ID {
			return fmt.Errorf("User not found")
		}

		*user = *results[0]

		return nil
	}
}

func testAccCheckWavefrontUserAttributes(user *wavefront.User, permissions []string, groups []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, v := range permissions {
			if sort.SearchStrings(user.Permissions, v) == len(user.Permissions) {
				return fmt.Errorf("permission not found or present on user. %s", v)
			}
		}

		for _, v := range groups {
			found := false
			for _, g := range user.Groups.UserGroups {
				if *g.ID == v {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("group not found or present on user. %s", v)
			}
		}
		return nil
	}
}

func testAccCheckWavefrontUser_basic() string {
	return fmt.Sprintf(`
resource "wavefront_user" "basic" {
	email  = "test+tftesting@example.com"
	groups = [
		"agent_management",
		"alerts_management",
	]
}`)
}

func testAccCheckWavefrontUser_changeGroups() string {
	return fmt.Sprintf(`
resource "wavefront_user" "basic" {
	email  = "test+tftesting@example.com"
	groups = [
		"agent_management",
		"events_management",
	]
}`)
}

func testAccCheckWavefrontUser_changeEmail() string {
	return fmt.Sprintf(`
resource "wavefront_user" "basic" {
	email  = "test+tftesting2@example.com"
	groups = [
		"agent_management",
		"alerts_management",
	]
}`)
}
