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
					resource.TestCheckResourceAttr("data.vmc_connected_accounts.my_accounts", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.vmc_connected_accounts.my_accounts", "ids.0", "2968040b-5c14-373f-8353-79c3a28a673b"),
				),
			},
		},
	})
}

func testAccDataSourceVmcConnectedAccountsConfig() string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = %q
}
	
data "vmc_org" "my_org" {
	id = "54937bce-8119-4fae-84f5-e5e066ee90e6"
}
	
data "vmc_connected_accounts" "my_accounts" {
	org_id = "${data.vmc_org.my_org.id}"
}
`,
		os.Getenv("REFRESH_TOKEN"),
	)
}
