package wavefront_plugin

import (
	"fmt"
	"strings"

	"github.com/MikeMcMahon/go-wavefront"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:             schema.TypeString,
				Optional:         true,
				Default:          wavefront.AlertTypeClassic,
				DiffSuppressFunc: suppressCase,
			},
			"target": {
				Type:             schema.TypeString,
				Optional:         true,
				StateFunc:        trimSpaces,
				DiffSuppressFunc: suppressSpaces,
			},
			"condition": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressSpaces,
			},
			"threshold_conditions": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"threshold_targets": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"additional_information": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressSpaces,
			},
			"display_expression": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressSpaces,
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
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressCase,
			},
			"tags": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlertCreate(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	var tags []string
	for _, tag := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, tag.(string))
	}

	a := &wavefront.Alert{
		Name:                               d.Get("name").(string),
		AdditionalInfo:                     trimSpaces(d.Get("additional_information")),
		DisplayExpression:                  trimSpaces(d.Get("display_expression")),
		Minutes:                            d.Get("minutes").(int),
		ResolveAfterMinutes:                d.Get("resolve_after_minutes").(int),
		NotificationResendFrequencyMinutes: d.Get("notification_resend_frequency_minutes").(int),
		Tags:                               tags,
	}

	err := validateAlertConditions(a, d)
	if err != nil {
		return err
	}

	// Create the alert on Wavefront
	err = alerts.Create(a)
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
			return nil
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
	d.Set("alert_type", tmpAlert.AlertType)
	d.Set("threshold_conditions", tmpAlert.Conditions)
	d.Set("threshold_targets", tmpAlert.Targets)

	return nil
}

func resourceAlertUpdate(d *schema.ResourceData, m interface{}) error {
	alerts := m.(*wavefrontClient).client.Alerts()

	alertID := d.Id()
	tmpAlert := wavefront.Alert{ID: &alertID}
	err := alerts.Get(&tmpAlert)

	d.SetId("")
	if err != nil {
		d.SetId("")
		return nil
	}

	var tags []string
	for _, tag := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, tag.(string))
	}

	a := tmpAlert
	a.Name = d.Get("name").(string)
	a.AdditionalInfo = trimSpaces(d.Get("additional_information").(string))
	a.DisplayExpression = trimSpaces(d.Get("display_expression").(string))
	a.Minutes = d.Get("minutes").(int)
	a.ResolveAfterMinutes = d.Get("resolve_after_minutes").(int)
	a.NotificationResendFrequencyMinutes = d.Get("notification_resend_frequency_minutes").(int)
	a.Tags = tags

	err = validateAlertConditions(&a, d)
	if err != nil {
		return err
	}

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

func validateAlertConditions(a *wavefront.Alert, d *schema.ResourceData) error {
	alertType := strings.ToUpper(d.Get("alert_type").(string))
	if alertType == wavefront.AlertTypeThreshold {
		a.AlertType = wavefront.AlertTypeThreshold
		if conditions, ok := d.GetOk("threshold_conditions"); ok {
			a.Conditions = trimSpacesMap(conditions.(map[string]interface{}))
			err := validateThresholdLevels(a.Conditions)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("threshold_conditions must be supplied for threshold alerts")
		}

		if targets, ok := d.GetOk("threshold_targets"); ok {
			a.Targets = trimSpacesMap(targets.(map[string]interface{}))
			return validateThresholdLevels(a.Targets)
		}

	} else if alertType == wavefront.AlertTypeClassic {
		a.AlertType = wavefront.AlertTypeClassic

		if d.Get("condition") == "" {
			return fmt.Errorf("condition must be supplied for classic alerts")
		}
		a.Condition = trimSpaces(d.Get("condition").(string))

		if d.Get("severity") == "" {
			return fmt.Errorf("severity must be supplied for classic alerts")
		}
		a.Severity = d.Get("severity").(string)
		a.Target = d.Get("target").(string)
	} else {
		return fmt.Errorf("alert_type must be CLASSIC or THRESHOLD")
	}

	return nil
}

func validateThresholdLevels(m map[string]string) error {
	for key := range m {
		ok := false
		for _, level := range []string{"severe", "warn", "info", "smoke"} {
			if key == level {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("invalid severity: %s", key)
		}
	}
	return nil
}
