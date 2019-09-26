package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/utils"
)

// Provider for VMware VMC Console APIs. Returns terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vmc_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://vmc.vmware.com/vmc/api",
			},
			"csp_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://console.cloud.vmware.com",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmc_sddc": resourceSddc(),
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
	connector, err := utils.NewVmcConnector(refreshToken, "", "")
	if err != nil {
		return connector,fmt.Errorf("Error creating connector : %v ", err)
	}
   return connector, nil
}
