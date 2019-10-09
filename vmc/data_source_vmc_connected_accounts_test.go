package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceVmcConnectedAccounts_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcConnectedAccountsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_connected_accounts.my_accounts", "ids.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceVmcConnectedAccountsConfig() string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = %q
    csp_url       = "https://console-stg.cloud.vmware.com"
    vmc_url = "https://stg.skyscraper.vmware.com"
}
	
data "vmc_org" "my_org" {
	id = "05e0a625-3293-41bb-a01f-35e762781c2a"
}
	
data "vmc_connected_accounts" "my_accounts" {
	org_id = "${data.vmc_org.my_org.id}"
}
`,
		os.Getenv("REFRESH_TOKEN"),
	)
}
