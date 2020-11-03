package wavefront

import (
	"fmt"

	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlertTarget() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlertTargetRead,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"recipient": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"triggers": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"template": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlertTargetRead(d *schema.ResourceData, m interface{}) error {
	targets := m.(*wavefrontClient).client.Targets()

	title, _ := d.GetOk("title")

	conds := []*wavefront.SearchCondition{{
		Key:            "title",
		Value:          title.(string),
		MatchingMethod: "EXACT",
	}}

	method, ok := d.GetOk("method")

	if ok {
		conds = append(conds,
			&wavefront.SearchCondition{
				Key:            "method",
				Value:          method.(string),
				MatchingMethod: "EXACT",
			})
	} else {
		method = "any type"
	}

	results, err := targets.Find(conds)

	if err != nil {
		return fmt.Errorf("error finding '%s' alert target '%s' in Wavefront: %s", method, title, err)
	}

	if len(results) == 0 {
		return fmt.Errorf("did not find '%s' alert target '%s' in Wavefront", method, title)
	}

	if len(results) > 1 {
		return fmt.Errorf("found multiple '%s' alert targets '%s' in Wavefront", method, title)
	}

	target := results[0]
	d.SetId(*target.ID)
	d.Set("recipient", target.Recipient)
	d.Set("description", target.Description)
	d.Set("triggers", target.Triggers)
	d.Set("template", target.Template)

	return nil
}
