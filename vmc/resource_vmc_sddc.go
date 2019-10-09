package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/model"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/sddcs"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/sddcs/esxs"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/tasks"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/runtime/protocol/client"
	"time"
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
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of this resource",
			},
			"storage_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"sddc_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_link_sddc_config": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer_subnet_ids": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								// Optional: true,
							},
							Optional: true,
						},
						"connected_account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				ForceNew: true,
			},
			"vpc_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"num_host": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"sddc_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vxlan_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// TODO check the deprecation statement
			"delay_account_link": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			// TODO change default to AWS
			"provider_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "ZEROCLOUD",
			},
			"skip_creating_vxlan": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"sso_domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "vmc.local",
			},
			"sddc_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "SingleAZ",
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "us-west-2",
			},
		},
	}
}

func resourceSddcCreate(d *schema.ResourceData, m interface{}) error {
	connector := m.(client.Connector)
	sddcClient := sddcs.NewSddcsClientImpl(connector)

	orgID := d.Get("org_id").(string)
	storageCapacity := d.Get("storage_capacity").(int)
	storageCapacityConverted := int64(storageCapacity)
	sddcName := d.Get("sddc_name").(string)
	vpcCidr := d.Get("vpc_cidr").(string)
	numHost := d.Get("num_host").(int)
	sddcType := d.Get("sddc_type").(string)
	var sddcTypePtr *string
	if sddcType != "" {
		sddcTypePtr = &sddcType
	}
	vxlanSubnet := d.Get("vxlan_subnet").(string)
	delayAccountLink := d.Get("delay_account_link").(bool)
	accountLinkConfig := &model.AccountLinkConfig{
		DelayAccountLink: &delayAccountLink,
	}
	providerType := d.Get("provider_type").(string)
	skipCreatingVxlan := d.Get("skip_creating_vxlan").(bool)
	ssoDomain := d.Get("sso_domain").(string)
	sddcTemplateID := d.Get("sddc_template_id").(string)
	deploymentType := d.Get("deployment_type").(string)
	region := d.Get("region").(string)
	accountLinkSddcConfig := expandAccountLinkSddcConfig(d.Get("account_link_sddc_config").([]interface{}))

	var awsSddcConfig = &model.AwsSddcConfig{
		StorageCapacity:       &storageCapacityConverted,
		Name:                  sddcName,
		VpcCidr:               &vpcCidr,
		NumHosts:              int64(numHost),
		SddcType:              sddcTypePtr,
		VxlanSubnet:           &vxlanSubnet,
		AccountLinkConfig:     accountLinkConfig,
		Provider:              providerType,
		SkipCreatingVxlan:     &skipCreatingVxlan,
		AccountLinkSddcConfig: accountLinkSddcConfig,
		SsoDomain:             &ssoDomain,
		SddcTemplateId:        &sddcTemplateID,
		DeploymentType:        &deploymentType,
		Region:                region,
	}

	// Create a Sddc
	task, err := sddcClient.Create(orgID, *awsSddcConfig, nil)
	if err != nil {
		return fmt.Errorf("Error while creating sddc %s: %v", sddcName, err)
	}

	// Wait until Sddc is created
	sddcID := task.ResourceId
	fmt.Println("Inside SDDC create ")
	fmt.Println(*sddcID)
	d.SetId(*sddcID)
	tasksClient := tasks.NewTasksClientImpl(connector)

	return resource.Retry(300*time.Minute, func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error describing instance: %s", err))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("Expected instance to be created but was in state %s", *task.Status))
		}
		return resource.NonRetryableError(resourceSddcRead(d, m))
	})
}

func resourceSddcRead(d *schema.ResourceData, m interface{}) error {
	connector := m.(client.Connector)
	sddcClient := sddcs.NewSddcsClientImpl(connector)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)
	sddc, err := sddcClient.Get(orgID, sddcID)
	if err != nil {
		return fmt.Errorf("Error while getting sddc detail %s: %v", sddcID, err)
	}

	d.SetId(sddc.Id)
	d.Set("name", sddc.Name)
	d.Set("updated", sddc.Updated)
	d.Set("user_id", sddc.UserId)
	d.Set("updated_by_user_id", sddc.UpdatedByUserId)
	d.Set("created", sddc.Created)
	d.Set("version", sddc.Version)
	d.Set("updated_by_user_name", sddc.UpdatedByUserName)
	d.Set("user_name", sddc.UserName)
	d.Set("sddc_state", sddc.SddcState)
	d.Set("org_id", sddc.OrgId)
	d.Set("sddc_type", sddc.SddcType)
	d.Set("provider", sddc.Provider)
	d.Set("account_link_state", sddc.AccountLinkState)
	d.Set("sddc_access_state", sddc.SddcAccessState)
	d.Set("sddc_type", sddc.SddcType)
	return nil
}

func resourceSddcDelete(d *schema.ResourceData, m interface{}) error {
	connector := m.(client.Connector)
	sddcClient := sddcs.NewSddcsClientImpl(connector)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)
	task, err := sddcClient.Delete(orgID, sddcID, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err)
	}

	err = WaitForTask(connector, orgID, task.Id)

	if err != nil {
		return fmt.Errorf("Error while waiting for task %s: %v", task.Id, err)
	}
	d.SetId("")
	return nil
}

func resourceSddcUpdate(d *schema.ResourceData, m interface{}) error {
	connector := m.(client.Connector)
	esxsClient := esxs.NewEsxsClientImpl(connector)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)

	// Add,remove hosts
	if d.HasChange("num_host") {
		oldTmp, newTmp := d.GetChange("num_host")
		oldNum := oldTmp.(int)
		newNum := newTmp.(int)

		action := "add"
		diffNum := newNum - oldNum

		if newNum < oldNum {
			action = "remove"
			diffNum = oldNum - newNum
		}

		esxConfig := model.EsxConfig{
			NumHosts: int64(diffNum),
		}

		task, err := esxsClient.Create(orgID, sddcID, esxConfig, &action)

		if err != nil {
			return fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err)
		}
		err = WaitForTask(connector, orgID, task.Id)
		if err != nil {
			return fmt.Errorf("Error while waiting for task %s: %v", task.Id, err)
		}
	}

	// Update sddc name
	if d.HasChange("sddc_name") {
		sddcClient := sddcs.NewSddcsClientImpl(connector)
		newSDDCName := d.Get("sddc_name").(string)
		sddcPatchRequest := model.SddcPatchRequest{
			Name: &newSDDCName,
		}
		sddc, err := sddcClient.Patch(orgID, sddcID, sddcPatchRequest)

		if err != nil {
			return fmt.Errorf("Error while updating sddc's name %v", err)
		}
		d.Set("sddc_name", sddc.Name)
	}

	return resourceSddcRead(d, m)
}

func expandAccountLinkSddcConfig(l []interface{}) []model.AccountLinkSddcConfig {

	if len(l) == 0 {
		return nil
	}

	var configs []model.AccountLinkSddcConfig

	for _, config := range l {
		c := config.(map[string]interface{})
		fmt.Println("Value of C")
		fmt.Println(c)
		var subnetIds []string
		for _, subnetID := range c["customer_subnet_ids"].([]interface{}) {

			subnetIds = append(subnetIds, subnetID.(string))
			fmt.Println("Inside SDDC creation ")
			fmt.Println(subnetIds)
		}
		var connectedAccId = c["connected_account_id"].(string)
		con := model.AccountLinkSddcConfig{
			CustomerSubnetIds:  subnetIds,
			ConnectedAccountId: &connectedAccId,
		}

		configs = append(configs, con)
	}
	return configs
}
