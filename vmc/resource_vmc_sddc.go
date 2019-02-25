package vmc

import (
	"context"
	"fmt"
	"github.com/antihax/optional"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/vapi-sdk/vmc-go-sdk/vmc"
	"net/http"
)

func resourceSddc() *schema.Resource {
	return &schema.Resource{
		Create: resourceSddcCreate,
		Read:   resourceSddcRead,
		Update: resourceSddcUpdate,
		Delete: resourceSddcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"org_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of this resource",
			},
			"storage_capacity": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sddc_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_link_sddc_config": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
				Optional: true,
			},
			"vpc_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"num_host": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"sddc_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vxlan_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// TODO check the deprecation statement
			"account_link_config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			// TODO change default to AWS
			"provider_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "ZEROCLOUD",
			},
			"skip_creating_vxlan": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"sso_domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sddc_template_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployment_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SingleAZ",
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "US_WEST_2",
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSddcCreate(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	orgID := d.Get("org_id").(string)
	sddcName := d.Get("sddc_name").(string)
	// TODO add account linking
	// vpcCidr := d.Get("vpc_cidr").(string)
	numHost := d.Get("num_host").(int)
	// sddcType := d.Get("num_host").(string)
	// vxlanSubnet := d.Get("vxlan_subnet").(string)

	// // TODO
	// account_link_config := d.Get("account_link_config")

	// skipCreatingVxlan := d.Get("skip_creating_vxlan").(bool)
	// ssoDomain := d.Get("sso_domain").(string)
	// sddcTemplateId := d.Get("sddc_template_id").(string)

	providerType := d.Get("provider_type").(string)
	region := d.Get("region").(string)
	var awsSddcConfig = &vmc.AwsSddcConfig{
		Name:     sddcName,
		NumHosts: int32(numHost),
		Provider: providerType,
		Region:   region,
	}

	// Create a Sddc
	task, resp, err := vmcClient.SddcApi.OrgsOrgSddcsPost(context.Background(), orgID, *awsSddcConfig)
	if err != nil {
		return fmt.Errorf("Error while creating sddc %s: %v", sddcName, err)
	}

	// Wait until Sddc is created
	sddcID := task.ResourceId
	err = vmc.WaitForTask(vmcClient, orgID, task.Id)
	if err != nil {
		return fmt.Errorf("Error while waiting for task %s: %v", task.Id, err)
	}

	// Get Sddc detail
	sddc, resp, err := vmcClient.SddcApi.OrgsOrgSddcsSddcGet(context.Background(), orgID, sddcID)
	if err != nil {
		return fmt.Errorf("Error while getting sddc detail %s: %v", sddcID, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Sddc %s was not found", sddcID)
	}

	d.SetId(sddc.Id)
	d.Set("name", sddc.Name)
	d.Set("created", sddc.Created)

	return resourceSddcRead(d, m)
}

func resourceSddcRead(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)
	sddc, resp, err := vmcClient.SddcApi.OrgsOrgSddcsSddcGet(context.Background(), orgID, sddcID)
	if err != nil {
		return fmt.Errorf("Error while getting sddc detail %s: %v", sddcID, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	d.SetId(sddc.Id)
	d.Set("org_id", sddc.OrgId)
	d.Set("sddc_name", sddc.Name)
	d.Set("provider_type", sddc.Provider)
	d.Set("created", sddc.Created)
	return nil
}

func resourceSddcDelete(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)
	task, _, err := vmcClient.SddcApi.OrgsOrgSddcsSddcDelete(context.Background(), orgID, sddcID, nil)
	if err != nil {
		return fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err)
	}
	err = vmc.WaitForTask(vmcClient, orgID, task.Id)
	if err != nil {
		return fmt.Errorf("Error while waiting for task %s: %v", task.Id, err)
	}
	d.SetId("")
	return nil
}

func resourceSddcUpdate(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)
	oldTmp, newTmp := d.GetChange("num_host")
	oldNum := oldTmp.(int)
	newNum := newTmp.(int)

	action := "add"
	diffNum := newNum - oldNum

	if newNum < oldNum {
		action = "remove"
		diffNum = oldNum - newNum
	}

	esxConfig := vmc.EsxConfig{
		NumHosts: int32(diffNum),
	}

	actionString := optional.NewString(action)

	// API_CALL
	task, _, err := vmcClient.EsxApi.OrgsOrgSddcsSddcEsxsPost(context.Background(), orgID, sddcID, esxConfig, &vmc.OrgsOrgSddcsSddcEsxsPostOpts{Action: actionString})

	if err != nil {
		return fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err)
	}
	err = vmc.WaitForTask(vmcClient, orgID, task.Id)
	if err != nil {
		return fmt.Errorf("Error while waiting for task %s: %v", task.Id, err)
	}

	return resourceSddcRead(d, m)
}