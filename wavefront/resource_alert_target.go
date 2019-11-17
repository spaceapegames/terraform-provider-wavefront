package wavefront_plugin

import (
	"fmt"
	"github.com/MikeMcMahon/go-wavefront"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"triggers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"template": {
				Type:     schema.TypeString,
				Required: true,
			},
			// 'method' must be EMAIL, WEBHOOK or PAGERDUTY
			"method": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"recipient": {
				Type:     schema.TypeString,
				Required: true,
			},
			"route": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target": {
							Type:     schema.TypeString,
							Required: true,
						},
						"filter": {
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},
			// EMAIL targets only
			"email_subject": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// EMAIL targets only
			"is_html_content": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// WEBHOOK targets only
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// WEBHOOK targets only
			"custom_headers": {
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

	alertRoutes := []wavefront.AlertRoute{}
	routes := d.Get("route").(*schema.Set).List()
	for _, route := range routes {
		r := route.(map[string]interface{})
		f := r["filter"].(map[string]interface{})
		alertRoutes = append(alertRoutes, wavefront.AlertRoute{
			Method: r["method"].(string),
			Target: r["target"].(string),
			Filter: f["key"].(string) + " " + f["value"].(string),
		})
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
		IsHtmlContent: d.Get("is_html_content").(bool),
		Routes:        alertRoutes,
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
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		} else {
			d.SetId("")
			return fmt.Errorf("error finding Wavefront Alert Target %s. %s", d.Id(), err)
		}
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
	d.Set("is_html_content", tmpTarget.IsHtmlContent)
	d.Set("custom_headers", tmpTarget.CustomHeaders)

	// Convert the routes from AlertRoute -> Terraform Friendly
	if tmpTarget.Routes != nil {
		var routes []interface{}
		for _, route := range tmpTarget.Routes {
			alertRoute := make(map[string]interface{})
			filterKV := strings.Split(route.Filter, " ")
			alertRoute["method"] = route.Method
			alertRoute["target"] = route.Target
			alertRoute["filter"] = map[string]interface{}{
				"key":   filterKV[0],
				"value": filterKV[1],
			}

			routes = append(routes, alertRoute)
		}
		d.Set("routes", routes)
	}

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

	if len(results) == 0 {
		return fmt.Errorf("error finding Wavefront Alert Target %s", d.Id())
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
	t.IsHtmlContent = d.Get("is_html_content").(bool)
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
