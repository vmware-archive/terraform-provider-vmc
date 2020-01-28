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

# Using the provider

The following instructions assume that the required parameters have updated in [variables.tf](https://github.com/vmware/terraform-provider-vmc/blob/master/examples/variables.tf) with your infrastructure settings. 
In order to provide the required parameters you need to specify -var flag to the terraform commands.

Load the provider
---------------------

```sh
terraform init
```

Try a dry run
---------------------

```sh
terraform plan
```

OR

```sh
terraform plan -var="api_token=xxxx" -var="org_id=xxxx"
```

Check if the terraform plan looks good.

Execute the plan
---------------------

```sh
   terraform apply
```

OR

```sh
   terraform apply -var="api_token=xxxx" -var="org_id=xxxx"
```

Verify the SDDC has been created successfully.

Check the terraform state
-------------------------
```sh
terraform show
```

Add / Remove hosts
---------------------

Update the "num_host" field in [main.tf](https://github.com/vmware/terraform-provider-vmc/blob/master/examples/main.tf).


```sh
terraform apply
```
OR

You can update num_host from command line :

```sh
terraform apply -var="api_token=xxxx" -var="org_id=xxxx" -var="num_hosts=2"
```

Verify the hosts are added/removed successfully.

Delete the SDDC
-----------------

```sh
terraform destroy
```

OR

```sh
terraform destroy -var="api_token=xxxx" -var="org_id=xxxx"
```
