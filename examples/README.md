# Provision an SDDC Example

This is an example that supports SDDC management actions like creating an SDDC, 
updating or deleting an existing SDDC.

To run the example:

* Generate an API token using [VMware Cloud on AWS console] (https://vmc.vmware.com/console/)

* Update the variables required parameters api_token and org_id in [variables.tf](https://github.com/vmware/terraform-provider-vmc/blob/master/examples/variables.tf) with your infrastructure settings. 
 
* Alternately you can provide the required parameters through command line :
  
  ```sh 
  terraform apply \
  -var="api_token=xxxx" \
  -var="org_id=xxxx"
  ```

