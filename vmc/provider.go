package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/utils"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/tasks"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/runtime/protocol/client"
	"time"
)

// Provider for VMware VMC Console APIs. Returns terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vmc_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://vmc.vmware.com/vmc/api",
			},
			"csp_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://console.cloud.vmware.com",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmc_sddc": resourceSddc(),
			"vmc_publicips" : resourcePublicIP(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vmc_org":                dataSourceVmcOrg(),
			"vmc_connected_accounts": dataSourceVmcConnectedAccounts(),
			"vmc_customer_subnets":   dataSourceVmcCustomerSubnets(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	refreshToken := d.Get("refresh_token").(string)
	vmcURL := d.Get("vmc_url").(string)
	cspURL := d.Get("csp_url").(string)
	connector, err := utils.NewVmcConnector(refreshToken, vmcURL, cspURL)
	if err != nil {
		return connector, fmt.Errorf("Error creating connector : %v ", err)
	}
	return connector, nil
}

func WaitForTask(connector client.Connector, orgID string, taskID string) error {
	fmt.Printf("Wait for task %q to complete\n", taskID)
	tasksClient := tasks.NewTasksClientImpl(connector)

	for {

		task, err := tasksClient.Get(orgID, taskID)
		if err != nil {
			return fmt.Errorf("Error while getting task %s: %v", taskID, err)
		}

		if *task.Status == "STARTED" || *task.Status == "CANCELING" {
			waitInterval := 2 * time.Second
			fmt.Print(".")
			time.Sleep(waitInterval)
			continue
		}
		return nil
	}
}
