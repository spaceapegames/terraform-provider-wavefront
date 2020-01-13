package wavefront_plugin

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func resourceDerivedMetric() *schema.Resource {
	return &schema.Resource{
		Create: resourceDerivedMetricCreate,
		Read:   resourceDerivedMetricRead,
		Update: resourceDerivedMetricUpdate,
		Delete: resourceDerivedMetricDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"minutes": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: nil, // TODO - Validate to ensure > 0 minutes and < insane minutes
			},
			"additional_information": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDerivedMetricCreate(d *schema.ResourceData, m interface{}) error {
	derivedMetrics := m.(*wavefrontClient).client.DerivedMetrics()

	var tags []string
	for _, tag := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, tag.(string))
	}

	dm := &wavefront.DerivedMetric{
		Name:                  d.Get("name").(string),
		AdditionalInformation: d.Get("additional_information").(string),
		Query:                 d.Get("query").(string),
		Minutes:               d.Get("minutes").(int),
		Tags:                  wavefront.WFTags{CustomerTags: tags},
	}

	err := derivedMetrics.Create(dm)
	if err != nil {
		return fmt.Errorf("error creating Derived Metric %s. %s", d.Get("name"), err)
	}

	d.SetId(*dm.ID)

	return nil
}

func resourceDerivedMetricUpdate(d *schema.ResourceData, m interface{}) error {
	derivedMetrics := m.(*wavefrontClient).client.DerivedMetrics()

	derivedMetricId := d.Id()
	tmpDM := &wavefront.DerivedMetric{ID: &derivedMetricId}
	err := derivedMetrics.Get(tmpDM)

	if err != nil {
		d.SetId("")
		return nil
	}

	var tags []string
	for _, tag := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, tag.(string))
	}

	dm := tmpDM
	dm.Name = d.Get("name").(string)
	dm.Minutes = d.Get("minutes").(int)
	dm.AdditionalInformation = d.Get("additional_information").(string)
	dm.Query = d.Get("query").(string)
	dm.Tags = wavefront.WFTags{CustomerTags: tags}

	err = derivedMetrics.Update(dm)
	if err != nil {
		return fmt.Errorf("unable to update Wavefront Derived Metric %s, %s", derivedMetricId, err)
	}

	return nil
}

func resourceDerivedMetricRead(d *schema.ResourceData, m interface{}) error {
	derivedMetrics := m.(*wavefrontClient).client.DerivedMetrics()

	derivedMetricId := d.Id()
	tmpDM := &wavefront.DerivedMetric{ID: &derivedMetricId}
	err := derivedMetrics.Get(tmpDM)

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("unable to find Wavefront Derived Metric %s. %s", d.Id(), err)
		}
	}

	d.SetId(*tmpDM.ID)
	d.Set("name", tmpDM.Name)
	d.Set("minutes", tmpDM.Minutes)
	d.Set("additional_information", tmpDM.AdditionalInformation)
	d.Set("query", tmpDM.Query)
	d.Set("tags", tmpDM.Tags.CustomerTags)

	return nil
}

func resourceDerivedMetricDelete(d *schema.ResourceData, m interface{}) error {
	derivedMetrics := m.(*wavefrontClient).client.DerivedMetrics()

	derivedMetricId := d.Id()
	tmpDM := &wavefront.DerivedMetric{ID: &derivedMetricId}
	err := derivedMetrics.Get(tmpDM)

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("unable to find Wavefront Derived Metric %s. %s", d.Id(), err)
		}
	}

	err = derivedMetrics.Delete(tmpDM)
	if err != nil {
		return fmt.Errorf("error trying to delete Wavefront Derived Metric %s. %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
