package wavefront_plugin

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spaceapegames/go-wavefront"
)

func resourceTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceTargetCreate,
		Read:   resourceTargetRead,
		Update: resourceTargetUpdate,
		Delete: resourceTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"triggers": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// 'method' must be EMAIL, WEBHOOK or PAGERDUTY
			"method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"recipient": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			// EMAIL targets only
			"email_subject": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// WEBHOOK targets only
			"content_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// WEBHOOK targets only
			"custom_headers": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceTargetCreate(d *schema.ResourceData, m interface{}) error {
	targets := m.(*wavefrontClient).client.Targets()

	var triggers []string
	for _, trigger := range d.Get("triggers").([]interface{}) {
		triggers = append(triggers, trigger.(string))
	}

	customHeaders := make(map[string]string)
	for k, v := range d.Get("custom_headers").(map[string]interface{}) {
		customHeaders[k] = v.(string)
	}

	t := &wavefront.Target{
		Title:         d.Get("name").(string),
		Description:   d.Get("description").(string),
		Triggers:      triggers,
		Template:      d.Get("template").(string),
		Method:        d.Get("method").(string),
		Recipient:     d.Get("recipient").(string),
		EmailSubject:  d.Get("email_subject").(string),
		ContentType:   d.Get("content_type").(string),
		CustomHeaders: customHeaders,
	}

	// Create the Target on Wavefront
	err := targets.Create(t)
	if err != nil {
		return fmt.Errorf("Error Creating Target %s. %s", d.Get("name"), err)
	}

	d.SetId(*t.ID)

	return nil
}

func resourceTargetRead(d *schema.ResourceData, m interface{}) error {
	targets := m.(*wavefrontClient).client.Targets()

	targetID := d.Id()
	tmpTarget := wavefront.Target{ID: &targetID}
	err := targets.Get(&tmpTarget)
	if err != nil {
		return fmt.Errorf("Error finding Wavefront Target %s. %s", d.Id(), err)
	}

	// Use the Wavefront ID as the Terraform ID
	d.SetId(*tmpTarget.ID)
	d.Set("name", tmpTarget.Title)
	d.Set("description", tmpTarget.Description)
	d.Set("triggers", tmpTarget.Triggers)
	d.Set("template", tmpTarget.Template)
	d.Set("method", tmpTarget.Method)
	d.Set("recipient", tmpTarget.Recipient)
	d.Set("email_subject", tmpTarget.EmailSubject)
	d.Set("content_type", tmpTarget.ContentType)
	d.Set("custom_headers", tmpTarget.CustomHeaders)

	return nil
}

func resourceTargetUpdate(d *schema.ResourceData, m interface{}) error {
	targets := m.(*wavefrontClient).client.Targets()

	results, err := targets.Find(
		[]*wavefront.SearchCondition{
			{
				Key:            "id",
				Value:          d.Id(),
				MatchingMethod: "EXACT",
			},
		})
	if err != nil {
		return fmt.Errorf("Error finding Wavefront Target %s. %s", d.Id(), err)
	}

	var triggers []string
	for _, trigger := range d.Get("triggers").([]interface{}) {
		triggers = append(triggers, trigger.(string))
	}

	customHeaders := make(map[string]string)
	for k, v := range d.Get("custom_headers").(map[string]interface{}) {
		customHeaders[k] = v.(string)
	}

	t := results[0]
	t.Title = d.Get("name").(string)
	t.Description = d.Get("description").(string)
	t.Triggers = triggers
	t.Template = d.Get("template").(string)
	t.Method = d.Get("method").(string)
	t.Recipient = d.Get("recipient").(string)
	t.EmailSubject = d.Get("email_subject").(string)
	t.ContentType = d.Get("content_type").(string)
	t.CustomHeaders = customHeaders

	// Update the Target on Wavefront
	err = targets.Update(t)
	if err != nil {
		return fmt.Errorf("Error Updating Target %s. %s", d.Get("name"), err)
	}
	return nil
}

func resourceTargetDelete(d *schema.ResourceData, m interface{}) error {
	targets := m.(*wavefrontClient).client.Targets()

	results, err := targets.Find(
		[]*wavefront.SearchCondition{
			&wavefront.SearchCondition{
				Key:            "id",
				Value:          d.Id(),
				MatchingMethod: "EXACT",
			},
		})
	if err != nil {
		return fmt.Errorf("Error finding Wavefront Target %s. %s", d.Id(), err)
	}
	t := results[0]

	// Delete the Target
	err = targets.Delete(t)
	if err != nil {
		return fmt.Errorf("Failed to delete Target %s. %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
