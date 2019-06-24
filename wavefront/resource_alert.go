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
			"alert_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default: wavefront.AlertTypeClassic,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": {
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: trimSpaces,
			},
			"threshold_conditions": {
				Type:      schema.TypeMap,
				Optional:  true,
			},
			"threshold_targets": {
				Type:      schema.TypeMap,
				Optional:  true,
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
			"notification_resend_frequency_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
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

func mapToThresholdConditions(m map[string]string) (*wavefront.ThresholdConditions, error) {
	var t wavefront.ThresholdConditions
	for sev, ts := range m {
		switch sev {
		case "severe":
			t.Severe = trimSpaces(ts)
		case "warn":
			t.Warn = trimSpaces(ts)
		case "info":
			t.Info = trimSpaces(ts)
		case "smoke":
			t.Smoke = trimSpaces(ts)
		default:
			return nil, fmt.Errorf("invalid severity: %s", sev)
		}
	}

	return &t, nil
}

func mapToThresholdTargets(m map[string]string) (*wavefront.ThresholdTargets, error) {
	var t wavefront.ThresholdTargets
	for sev, ts := range m {
		switch sev {
		case "severe":
			t.Severe = trimSpaces(ts)
		case "warn":
			t.Warn = trimSpaces(ts)
		case "info":
			t.Info = trimSpaces(ts)
		case "smoke":
			t.Smoke = trimSpaces(ts)
		default:
			return nil, fmt.Errorf("invalid severity: %s", sev)
		}
	}

	return &t, nil
}

func resourceAlertCreate(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	var tags []string
	for _, tag := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, tag.(string))
	}

	// Do some validation here? That we either have a condition or threshold_condition
	// Also need to set the type accordinglt

	a := &wavefront.Alert{
		Name:                               d.Get("name").(string),
		AdditionalInfo:                     trimSpaces(d.Get("additional_information").(string)),
		DisplayExpression:                  trimSpaces(d.Get("display_expression").(string)),
		Minutes:                            d.Get("minutes").(int),
		ResolveAfterMinutes:                d.Get("resolve_after_minutes").(int),
		NotificationResendFrequencyMinutes: d.Get("notification_resend_frequency_minutes").(int),
		Tags:                               tags,
	}

	if d.Get("alert_type") == wavefront.AlertTypeThreshold {   //DO THIS!! YOUVE CHANGED IT
		a.Type = wavefront.AlertTypeThreshold
		conditions, err := mapToThresholdConditions(d.Get("threshold_conditions").(map[string]string))
		if err != nil {
			return err
		}
		a.ThresholdConditions = *conditions

		targets, err := mapToThresholdTargets(d.Get("threshold_conditions").(map[string]string))
		if err != nil {
			return err
		}
		a.ThresholdTargets = *targets
	} else if d.Get("condition") != nil {
		// This is a CLASSIC alert
		a.Type = wavefront.AlertTypeClassic
		a.Condition = trimSpaces(d.Get("condition").(string))
		if d.Get("severity") == nil {
			return fmt.Errorf("severity must be supplied with a classic alert")
		}
		a.Severity =  d.Get("severity").(string)
		a.Target =  d.Get("target").(string)
	} else {
		return fmt.Errorf("either condition or threshold_conditions is required")
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
	d.Set("condition", trimSpaces(tmpAlert.Condition))
	d.Set("additional_information", trimSpaces(tmpAlert.AdditionalInfo))
	d.Set("display_expression", trimSpaces(tmpAlert.DisplayExpression))
	d.Set("minutes", tmpAlert.Minutes)
	d.Set("resolve_after_minutes", tmpAlert.ResolveAfterMinutes)
	d.Set("notification_resend_frequency_minutes", tmpAlert.NotificationResendFrequencyMinutes)
	d.Set("severity", tmpAlert.Severity)
	d.Set("tags", tmpAlert.Tags)
	d.Set("threshold_conditions", tmpAlert.ThresholdConditions)

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
	a.Condition = trimSpaces(d.Get("condition").(string))
	a.AdditionalInfo = trimSpaces(d.Get("additional_information").(string))
	a.DisplayExpression = trimSpaces(d.Get("display_expression").(string))
	a.Minutes = d.Get("minutes").(int)
	a.ResolveAfterMinutes = d.Get("resolve_after_minutes").(int)
	a.NotificationResendFrequencyMinutes = d.Get("notification_resend_frequency_minutes").(int)
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
