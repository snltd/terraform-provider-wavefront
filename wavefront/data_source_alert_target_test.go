package wavefront

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTargetDataNonExistentTarget(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      nonExistentTarget,
				ExpectError: regexp.MustCompile(".*did not find 'any type' alert target 'no-such-thing-as-this' in Wavefront"),
			},
		},
	})
}

func TestAccTargetDataWrongTarget(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTarget,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("wavefront_alert_target.res", "id", regexp.MustCompile("^[a-zA-Z0-9]{16}$")),
					resource.TestCheckResourceAttr("wavefront_alert_target.res", "method", "PAGERDUTY"),
				),
			},
			{
				Config:      correctTargetWithWrongMethod,
				ExpectError: regexp.MustCompile(".*did not find 'WEBHOOK' alert target 'Terraform Test Target' in Wavefront"),
			},
		},
	})
}

func TestAccTargetDataNoMethodType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testTarget,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("wavefront_alert_target.res", "id", regexp.MustCompile("^[a-zA-Z0-9]{16}$")),
					resource.TestCheckResourceAttr("wavefront_alert_target.res", "method", "PAGERDUTY"),
				),
			},
			{
				Config: correctTarget,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.wavefront_alert_target.dat", "id", regexp.MustCompile("^[a-zA-Z0-9]{16}$")),
					resource.TestCheckResourceAttr("data.wavefront_alert_target.dat", "template", "{}"),
					resource.TestCheckResourceAttr("data.wavefront_alert_target.dat", "description", "Test target"),
					resource.TestCheckResourceAttr("data.wavefront_alert_target.dat", "recipient", "12345678910111213141516171819202"),
				),
			},
		},
	})
}

func TestAccTargetDataMethodType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTarget,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("wavefront_alert_target.res", "id", regexp.MustCompile("^[a-zA-Z0-9]{16}$")),
					resource.TestCheckResourceAttr("wavefront_alert_target.res", "method", "PAGERDUTY"),
				),
			},
			{
				Config: correctTargetWithMethod,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.wavefront_alert_target.dat", "id", regexp.MustCompile("^[a-zA-Z0-9]{16}$")),
					resource.TestCheckResourceAttr("data.wavefront_alert_target.dat", "template", "{}"),
					resource.TestCheckResourceAttr("data.wavefront_alert_target.dat", "description", "Test target"),
					resource.TestCheckResourceAttr("data.wavefront_alert_target.dat", "recipient", "12345678910111213141516171819202"),
				),
			},
		},
	})
}

const (
	nonExistentTarget = `
data wavefront_alert_target dat {
 title = "no-such-thing-as-this"
}
`
	correctTarget = `
data wavefront_alert_target dat {
	title = "Terraform Test Target"
}`

	correctTargetWithMethod = `
data wavefront_alert_target dat {
	title = "Terraform Test Target"
	method = "PAGERDUTY"
}`

	correctTargetWithWrongMethod = `
data wavefront_alert_target dat {
	title = "Terraform Test Target"
	method = "WEBHOOK"
}`

	testTarget = `
resource wavefront_alert_target res {
  name        = "Terraform Test Target"
  description = "Test target"
  method      = "PAGERDUTY"
  recipient   = "12345678910111213141516171819202"
  template    = "{}"
  triggers    = [
    "ALERT_OPENED",
  	"ALERT_RESOLVED"
  ]
}`
)
