package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/sddcs/publicips"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/runtime/protocol/client"
	"os"
	"testing"
)

func TestAccResourceVmcPublicIP_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSddcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVmcPublicIPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcPublicIPExists("vmc_publicips.publicip_1"),
				),
			},
		},
	})
}

func testCheckVmcPublicIPExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		allocationID := rs.Primary.Attributes["id"]
		sddcID := rs.Primary.Attributes["sddc_id"]
		orgID := rs.Primary.Attributes["org_id"]
		connector := testAccProvider.Meta().(client.Connector)
		publicIPClient := publicips.NewPublicipsClientImpl(connector)

		publicIP, err := publicIPClient.Get(orgID, sddcID, allocationID)
		if err != nil {
			return fmt.Errorf("Bad: Get on publicIP: %s", err)
		}

		if *publicIP.AllocationId != allocationID {
			return fmt.Errorf("Bad: Public IP %q does not exist", allocationID)
		}

		fmt.Printf("Public IP created successfully with id %s ", allocationID)
		return nil
	}
}

/*
func testCheckVmcSddcDestroy(s *terraform.State) error {

	connector := testAccProvider.Meta().(client.Connector)
	sddcClient := sddcs.NewSddcsClientImpl(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_sddc" {
			continue
		}

		sddcID := rs.Primary.Attributes["id"]
		orgID := rs.Primary.Attributes["org_id"]
		task, err := sddcClient.Delete(orgID, sddcID, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("Error while deleting sddc %s, %s", sddcID, err)
		}
		err = WaitForTask(connector, orgID, task.Id)
		if err != nil {
			return fmt.Errorf("Error while waiting for task %q: %v", task.Id, err)
		}
	}

	return nil
}*/

func testAccVmcPublicIPConfigBasic() string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = %q
	
	# refresh_token = "ac5140ea-1749-4355-a892-56cff4893be0"
	 csp_url       = "https://console-stg.cloud.vmware.com"
    vmc_url = "https://stg.skyscraper.vmware.com"
}
	
data "vmc_org" "my_org" {
	id = "05e0a625-3293-41bb-a01f-35e762781c2a"
}

resource "vmc_publicips" "publicip_1" {
	org_id = "${data.vmc_org.my_org.id}"
	sddc_id = "85de9b00-a442-448f-9ea8-bb10d8b377b4"
	names     = ["srege-test-VM-999"]
	host_count = 3
	private_ips = ["10.105.167.133"]
}
`,
		os.Getenv("REFRESH_TOKEN"),
	)
}
