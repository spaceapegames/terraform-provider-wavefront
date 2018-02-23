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

// Construct a Wavefront Alert
func buildAlert(d *schema.ResourceData) (*wavefront.Alert, error) {

	var tags []string
	for _, tag := range d.Get("tags").([]interface{}) {
		tags = append(tags, tag.(string))
	}

	return &wavefront.Alert{
		Name:                d.Get("name").(string),
		Target:              d.Get("target").(string),
		Condition:           d.Get("condition").(string),
		AdditionalInfo:      d.Get("additional_information").(string),
		DisplayExpression:   d.Get("display_expression").(string),
		Minutes:             d.Get("minutes").(int),
		ResolveAfterMinutes: d.Get("resolve_after_minutes").(int),
		Severity:            d.Get("severity").(string),
		Tags:                tags,
	}, nil
}

// Create the alert on Wavefront
func resourceAlertCreate(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()
	alert, err := buildAlert(d)
	if err != nil {
		return fmt.Errorf("Failed to parse Alert, %s", err)
	}

	err = alerts.Create(alert)
	if err != nil {
		return fmt.Errorf("Error Creating Alert %s. %s", d.Get("name"), err)
	}
	d.SetId(alert.ID)

	return nil
}

// Read a Wavefront Alert
func resourceAlertRead(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	alert := wavefront.Alert{
		ID: d.Id(),
	}

	// search for an dashboard with our id. We should receive 1 (Exact Match) or 0 (No Match)
	err := alerts.Get(&alert)
	if err != nil {
		// alert no longer exists
		d.SetId("")
	}

	// Use the Wavefront ID as the Terraform ID
	d.SetId(alert.ID)
	d.Set("name", alert.Name)
	d.Set("target", alert.Target)
	d.Set("condition", alert.Condition)
	d.Set("additional_information", alert.AdditionalInfo)
	d.Set("display_expression", alert.DisplayExpression)
	d.Set("minutes", alert.Minutes)
	d.Set("resolve_after_minutes", alert.ResolveAfterMinutes)
	d.Set("severity", alert.Severity)
	d.Set("tags", alert.Tags)

	return nil
}

// Update the alert on Wavefront
func resourceAlertUpdate(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()
	alert, err := buildAlert(d)
	if err != nil {
		return fmt.Errorf("Failed to parse Alert, %s", err)
	}

	err = alerts.Update(alert)
	if err != nil {
		return fmt.Errorf("Error Updating Alert %s. %s", d.Get("name"), err)
	}
	return nil
}

// Delete the alert on Wavefront
func resourceAlertDelete(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()
	alert := wavefront.Alert{
		ID: d.Id(),
	}

	err := alerts.Get(&alert)
	if err != nil {
		return fmt.Errorf("Error finding Wavefront Alert %s. %s", d.Id(), err)
	}

	// Delete the Alert
	err = alerts.Delete(&alert)
	if err != nil {
		return fmt.Errorf("Failed to delete Alert %s. %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
