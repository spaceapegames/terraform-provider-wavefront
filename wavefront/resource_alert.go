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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": {
				Type:      schema.TypeString,
				Required:  true,
				StateFunc: trimSpaces,
			},
			"additional_information": {
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: trimSpaces,
			},
			"display_expression": {
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: trimSpaces,
			},
			"minutes": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"resolve_after_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeSet,
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
	for _, tag := range d.Get("tags").(*schema.Set).List() {
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
		return fmt.Errorf("error creating Alert %s. %s", d.Get("name"), err)
	}

	d.SetId(*a.ID)

	return nil
}

func resourceAlertRead(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	alertID := d.Id()
	tmpAlert := wavefront.Alert{ID: &alertID}
	err := alerts.Get(&tmpAlert)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Wavefront Alert %s. %s", d.Id(), err)
		}
	}

	// Use the Wavefront ID as the Terraform ID
	d.SetId(*tmpAlert.ID)
	d.Set("name", tmpAlert.Name)
	d.Set("target", tmpAlert.Target)
	d.Set("condition", tmpAlert.Condition)
	d.Set("additional_information", tmpAlert.AdditionalInfo)
	d.Set("display_expression", tmpAlert.DisplayExpression)
	d.Set("minutes", tmpAlert.Minutes)
	d.Set("resolve_after_minutes", tmpAlert.ResolveAfterMinutes)
	d.Set("severity", tmpAlert.Severity)
	d.Set("tags", tmpAlert.Tags)

	return nil
}

func resourceAlertUpdate(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	alertID := d.Id()
	tmpAlert := wavefront.Alert{ID: &alertID}
	err := alerts.Get(&tmpAlert)

	if err != nil {
		return fmt.Errorf("Error finding Wavefront Alert %s. %s", d.Id(), err)
	}

	var tags []string
	for _, tag := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, tag.(string))
	}

	a := tmpAlert
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
	err = alerts.Update(&a)
	if err != nil {
		return fmt.Errorf("Error Updating Alert %s. %s", d.Get("name"), err)
	}
	return nil
}

func resourceAlertDelete(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	alertID := d.Id()
	tmpAlert := wavefront.Alert{ID: &alertID}
	err := alerts.Get(&tmpAlert)
	if err != nil {
		return fmt.Errorf("Error finding Wavefront Alert %s. %s", d.Id(), err)
	}
	a := tmpAlert

	// Delete the Alert
	err = alerts.Delete(&a)
	if err != nil {
		return fmt.Errorf("Failed to delete Alert %s. %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
