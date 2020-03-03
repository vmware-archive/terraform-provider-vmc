package vmc

import (
	"fmt"
	"log"
)

func GetSDDC(connector client.Connector, orgID string, sddcID string) (model.Sddc, error) {
	sddcClient := orgs.NewDefaultSddcsClient(connector)
	sddc, err := sddcClient.Get(orgID, sddcID)
	return sddc, err
}

func DeleteSDDC(d *schema.ResourceData, connector client.Connector, orgID string, sddcID string) error {
	sddcClient := orgs.NewDefaultSddcsClient(connector)
	task, err := sddcClient.Delete(orgID, sddcID, nil, nil, nil)
	if err != nil {
		if err.Error() == errors.NewInvalidRequest().Error() {
			log.Printf("Can't Delete : SDDC with ID %s not found or already deleted %v", sddcID, err)
			return nil
		}
		return fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err)
	}
	tasksClient := orgs.NewDefaultTasksClient(connector)
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("Expected instance to be deleted but was in state %s", *task.Status))
		}
		d.SetId("")
		return resource.NonRetryableError(nil)
	})
}
