package wavefront_plugin

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Exists: resourceUserExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"user_groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"customer": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	users := m.(*wavefrontClient).client.Users()

	newUserRequest := &wavefront.NewUserRequest{
		EmailAddress: d.Get("email").(string),
	}

	err := resourceDecodeUserPermissions(d, newUserRequest)
	if err != nil {
		return fmt.Errorf("error extracting permisisons from terraform state. %s", err)
	}

	err = resourceDecodeUserGroups(d, newUserRequest)
	if err != nil {
		return fmt.Errorf("error extracting user groups from terraform state. %s", err)
	}

	user := &wavefront.User{}
	if err := users.Create(newUserRequest, user, true); err != nil {
		return fmt.Errorf("failed to create new user, %s", err)
	}

	d.SetId(*user.ID)

	return nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	users := m.(*wavefrontClient).client.Users()

	results, err := users.Find(
		[]*wavefront.SearchCondition{
			{
				Key:            "id",
				Value:          d.Id(),
				MatchingMethod: "EXACT",
			},
		})
	if err != nil {
		return fmt.Errorf("error finding Wavefront User %s. %s", d.Id(), err)
	}

	if len(results) == 0 {
		d.SetId("")
		return nil
	}

	user := results[0]

	d.Set("email", user.ID)
	d.Set("customer", user.Customer)
	d.Set("groups", user.Permissions)
	resourceEncodeUserGroups(d, user)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	users := m.(*wavefrontClient).client.Users()
	results, err := users.Find(
		[]*wavefront.SearchCondition{
			&wavefront.SearchCondition{
				Key:            "id",
				Value:          d.Id(),
				MatchingMethod: "EXACT",
			},
		})
	if err != nil {
		return fmt.Errorf("error finding Wavefront User %s. %s", d.Id(), err)
	}

	if len(results) == 0 {
		d.SetId("")
		return nil
	}

	u := results[0]
	emailAddress := d.Id()
	u.ID = &emailAddress

	err = resourceDecodeUserPermissions(d, u)
	if err != nil {
		return fmt.Errorf("error decoding permissions from state into the user %s. %s", d.Id(), err)
	}
	err = resourceDecodeUserGroups(d, u)
	if err != nil {
		return fmt.Errorf("error decoding user groups from state into the user %s. %s", d.Id(), err)
	}

	err = users.Update(u)
	if err != nil {
		return fmt.Errorf("error updating Wavefront User %s. %s", d.Id(), err)
	}

	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	users := m.(*wavefrontClient).client.Users()
	results, err := users.Find(
		[]*wavefront.SearchCondition{
			{
				Key:            "id",
				Value:          d.Id(),
				MatchingMethod: "EXACT",
			},
		})
	if err != nil {
		return fmt.Errorf("error finding Wavefront User %s. %s", d.Id(), err)
	}

	// Delete the user
	u := results[0]
	err = users.Delete(u)
	if err != nil {
		return fmt.Errorf("error deleting Wavefront User %s. %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceUserExists(d *schema.ResourceData, m interface{}) (bool, error) {
	users := m.(*wavefrontClient).client.Users()
	results, err := users.Find(
		[]*wavefront.SearchCondition{
			{
				Key:            "id",
				Value:          d.Id(),
				MatchingMethod: "EXACT",
			},
		})

	if err != nil {
		return false, fmt.Errorf("error finding Wavefront User %s. %s", d.Id(), err)
	}

	if len(results) == 0 {
		return false, nil
	}

	return true, nil
}

// Decodes the user groups from the state and assigns them to the User
func resourceDecodeUserGroups(d *schema.ResourceData, user interface{}) error {
	var userGroups *schema.Set
	if d.HasChange("user_groups") {
		_, n := d.GetChange("user_groups")

		// Largely fine if new is nil, likely means we're removing the user from all groups
		// Which default puts them back into the everyone group
		if n == nil {
			n = new(schema.Set)
		}
		userGroups = n.(*schema.Set)
	} else {
		userGroups = d.Get("user_groups").(*schema.Set)
	}

	var wfUserGroups []wavefront.UserGroup
	for _, ug := range userGroups.List() {
		if ug == nil {
			continue
		}
		ugId := ug.(string)
		wfUserGroups = append(wfUserGroups, wavefront.UserGroup{
			ID: &ugId,
		})
	}

	switch v := (user).(type) {
	case *wavefront.User:
		u := user.(*wavefront.User)
		u.Groups = wavefront.UserGroupsWrapper{UserGroups: wfUserGroups}
		user = u
	case *wavefront.NewUserRequest:
		u := user.(*wavefront.NewUserRequest)
		u.Groups = wavefront.UserGroupsWrapper{UserGroups: wfUserGroups}
		user = u
	default:
		return fmt.Errorf("unknown type attempted to cast %T", v)
	}

	return nil
}

// Encodes user groups from the User and assign them to the TF State
func resourceEncodeUserGroups(d *schema.ResourceData, user *wavefront.User) {
	var userGroups []string
	if len(user.Groups.UserGroups) > 0 {
		for _, g := range user.Groups.UserGroups {
			userGroups = append(userGroups, *g.ID)
		}
	}

	d.Set("user_groups", userGroups)
}

// Decodes the permissions (groups) from the state file and returns a []string of permissions
func resourceDecodeUserPermissions(d *schema.ResourceData, user interface{}) error {
	var permissions []string
	for _, permission := range d.Get("groups").(*schema.Set).List() {
		permissions = append(permissions, permission.(string))
	}

	switch v := user.(type) {
	case *wavefront.User:
		u := user.(*wavefront.User)
		u.Permissions = permissions
	case *wavefront.NewUserRequest:
		u := user.(*wavefront.NewUserRequest)
		u.Permissions = permissions
	default:
		return fmt.Errorf("unknown type attempted to cast %T", v)

	}

	return nil
}
