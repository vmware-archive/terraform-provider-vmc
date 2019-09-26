package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/utils"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/com/vmware/vmc/orgs"
	"gitlab.eng.vmware.com/vapi-sdk/vmc-go-sdk/vmc"
)

func dataSourceVmcOrg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcOrgRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Unique ID of this resource",
				Required:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource",
				Computed:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The Name of this resource",
				Computed:    true,
			},
		},
	}
}

func dataSourceVmcOrgRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*vmc.Client)
	orgID := d.Get("id").(string)
	connector, err := utils.NewVmcConnector(client.RefreshToken, "", "")
	if err != nil {
		return fmt.Errorf("Error while reading org information for %s: %v", orgID, err)
	}

	orgClient := orgs.NewOrgsClientImpl(connector)
	org, err := orgClient.Get(orgID)

	if err != nil {
		return fmt.Errorf("Error while reading org information for %s: %v", orgID, err)
	}
	d.SetId(org.Id)
	d.Set("display_name", org.DisplayName)
	d.Set("name", org.Name)

	return nil
}
