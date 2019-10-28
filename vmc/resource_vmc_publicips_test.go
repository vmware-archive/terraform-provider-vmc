package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/sddcs/publicips"
	"os"
	"testing"
)

func TestAccResourceVmcPublicIP_basic(t *testing.T) {
	VMName := "terraform_test_vm_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcPublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVmcPublicIPConfigBasic(VMName),
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

		sddcID := rs.Primary.Attributes["sddc_id"]
		orgID := rs.Primary.Attributes["org_id"]
		vmName := rs.Primary.Attributes["name"]
		connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
		connector := connectorWrapper.Connector
		publicIPClient := publicips.NewPublicipsClientImpl(connector)

		publicIPList, err := publicIPClient.List(orgID, sddcID)
		if err != nil {
			return fmt.Errorf("Bad: List on publicIP: %s", err)
		}
		var allocationID *string

		for i := range publicIPList {
			if *publicIPList[i].Name == vmName {
				allocationID = publicIPList[i].AllocationId
			}
		}

		publicIP, err := publicIPClient.Get(orgID, sddcID, *allocationID)
		if err != nil {
			return fmt.Errorf("Bad: Get on publicIP API: %s", err)
		}

		if *publicIP.Name != vmName {
			return fmt.Errorf("Bad: Public IP %q does not exist", *allocationID)
		}

		fmt.Printf("Public IP created successfully with id %s ", *allocationID)
		return nil
	}
}


func testCheckVmcPublicIPDestroy(s *terraform.State) error {

	connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
	connector := connectorWrapper.Connector
	publicIPClient := publicips.NewPublicipsClientImpl(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_sddc" {
			continue
		}

		allocationID := rs.Primary.Attributes["id"]
		orgID := rs.Primary.Attributes["org_id"]
		sddcID := rs.Primary.Attributes["sddc_id"]

		task, err := publicIPClient.Delete(orgID, sddcID,allocationID)
		if err != nil {
			return fmt.Errorf("Error while deleting sddc %s, %s", sddcID, err)
		}
		err = WaitForTask(connector, orgID, task.Id)
		if err != nil {
			return fmt.Errorf("Error while waiting for task %q: %v", task.Id, err)
		}
	}
	return nil
}

func testAccVmcPublicIPConfigBasic(name string) string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = %q
	
	# refresh_token = "ac5140ea-1749-4355-a892-56cff4893be0"
	 csp_url       = "https://console-stg.cloud.vmware.com"
    vmc_url = "https://stg.skyscraper.vmware.com"
}
	
data "vmc_org" "my_org" {
	id = %q
}

resource "vmc_publicips" "publicip_1" {
	org_id = "${data.vmc_org.my_org.id}"
	sddc_id = "30aa9e93-766d-498b-92aa-75f3b5304a7e"
	name     = %q
	private_ip = "10.105.167.133"
}
`,
		os.Getenv("REFRESH_TOKEN"),
		os.Getenv("ORG_ID"),
		name,
	)
}
