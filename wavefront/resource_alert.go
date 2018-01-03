package wavefront_plugin

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spaceapegames/go-wavefront"
)

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlertCreate,
		Read:   resourceAlertRead,
		Update: resourceAlertUpdate,
		Delete: resourceAlertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"target": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				StateFunc: trimSpaces,
			},
			"additional_information": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: trimSpaces,
			},
			"display_expression": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: trimSpaces,
			},
			"minutes": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"resolve_after_minutes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"severity": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func trimSpaces(d interface{}) string {
	return strings.TrimSpace(d.(string))
}

func resourceAlertCreate(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	var tags []string
	for _, tag := range d.Get("tags").([]interface{}) {
		tags = append(tags, tag.(string))
	}

	a := &wavefront.Alert{
		Name:                d.Get("name").(string),
		Target:              d.Get("target").(string),
		Condition:           d.Get("condition").(string),
		AdditionalInfo:      d.Get("additional_information").(string),
		DisplayExpression:   d.Get("display_expression").(string),
		Minutes:             d.Get("minutes").(int),
		ResolveAfterMinutes: d.Get("resolve_after_minutes").(int),
		Severity:            d.Get("severity").(string),
		Tags:                tags,
	}

	// Create the alert on Wavefront
	err := alerts.Create(a)
	if err != nil {
		return fmt.Errorf("Error Creating Alert %s. %s", d.Get("name"), err)
	}

	d.SetId(*a.ID)

	return nil
}

func resourceAlertRead(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	// search for an alert with our id. We should recieve 1 (Exact Match) or 0 (No Match)
	results, err := alerts.Find(
		[]*wavefront.SearchCondition{
			{
				Key:            "id",
				Value:          d.Id(),
				MatchingMethod: "EXACT",
			},
		})
	if err != nil {
		return fmt.Errorf("Error finding Wavefront Alert %s. %s", d.Id(), err)
	}
	// resource has been deleted out of band. So unset ID
	if len(results) != 1 {
		d.SetId("")
		return nil
	}

	// Use the Wavefront ID as the Terraform ID
	d.SetId(*results[0].ID)
	d.Set("name", results[0].Name)
	d.Set("target", results[0].Target)
	d.Set("condition", results[0].Condition)
	d.Set("additional_information", results[0].AdditionalInfo)
	d.Set("display_expression", results[0].DisplayExpression)
	d.Set("minutes", results[0].Minutes)
	d.Set("resolve_after_minutes", results[0].ResolveAfterMinutes)
	d.Set("severity", results[0].Severity)
	d.Set("tags", results[0].Tags)

	return nil
}

func resourceAlertUpdate(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	results, err := alerts.Find(
		[]*wavefront.SearchCondition{
			{
				Key:            "id",
				Value:          d.Id(),
				MatchingMethod: "EXACT",
			},
		})
	if err != nil {
		return fmt.Errorf("Error finding Wavefront Alert %s. %s", d.Id(), err)
	}

	var tags []string
	for _, tag := range d.Get("tags").([]interface{}) {
		tags = append(tags, tag.(string))
	}

	a := results[0]
	a.Name = d.Get("name").(string)
	a.Target = d.Get("target").(string)
	a.Condition = d.Get("condition").(string)
	a.AdditionalInfo = d.Get("additional_information").(string)
	a.DisplayExpression = d.Get("display_expression").(string)
	a.Minutes = d.Get("minutes").(int)
	a.ResolveAfterMinutes = d.Get("resolve_after_minutes").(int)
	a.Severity = d.Get("severity").(string)
	a.Tags = tags

	// Update the alert on Wavefront
	err = alerts.Update(a)
	if err != nil {
		return fmt.Errorf("Error Updating Alert %s. %s", d.Get("name"), err)
	}
	return nil
}

func resourceAlertDelete(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	results, err := alerts.Find(
		[]*wavefront.SearchCondition{
			&wavefront.SearchCondition{
				Key:            "id",
				Value:          d.Id(),
				MatchingMethod: "EXACT",
			},
		})
	if err != nil {
		return fmt.Errorf("Error finding Wavefront Alert %s. %s", d.Id(), err)
	}
	a := results[0]

	// Delete the Alert
	err = alerts.Delete(a)
	if err != nil {
		return fmt.Errorf("Failed to delete Alert %s. %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
