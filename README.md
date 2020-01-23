# Terraform provider for VMware Cloud on AWS

This is the repository for the Terraform provider for VMware Cloud, which one can use with
Terraform to work with [VMware Cloud on AWS](https://vmc.vmware.com/).

# Using the provider

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Build the Provider

Clone repository to: `$GOPATH/src/github.com/provider/`

```sh
mkdir -p $GOPATH/src/github.com/provider/
cd $GOPATH/src/github.com/provider/
git clone https://github.com/vmware/terraform-provider-vmc.git
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/provider/terraform-provider-vmc
go get
go build -o terraform-provider-vmc
```
## Load the provider
```sh
terraform init
```

## Connect to VMC and create a testing sddc

Update following fields in the [main.tf](main.tf) with your infrastructure settings

* refresh_token
* id
* sddc_name

```
provider "vmc" {
  refresh_token = ""
}

data "vmc_org" "my_org" {
  id = ""
}

data "vmc_connected_accounts" "my_accounts" {
  org_id = data.vmc_org.my_org.id
}

data "vmc_customer_subnets" "my_subnets" {
  org_id               = data.vmc_org.my_org.id
  connected_account_id = data.vmc_connected_accounts.my_accounts.ids[0]
  region               = var.sddc_region
}

resource "vmc_sddc" "sddc_1" {
  org_id = data.vmc_org.my_org.id

  sddc_name           = ""
  vpc_cidr            = var.vpc_cidr
  num_host            = 3
  provider_type       = "AWS"
  region              = data.vmc_customer_subnets.my_subnets.region
  vxlan_subnet        = var.vxlan_subnet
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"

  deployment_type = "SingleAZ"

  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.ids[0]
  }
  timeouts {
    create = "300m"
    update = "300m"
    delete = "180m"
  }
}

resource "vmc_publicips" "IP1" {
  org_id     = data.vmc_org.my_org.id
  sddc_id    = vmc_sddc.sddc_1.id
  private_ip = var.private_ip
  name       = "vm1"
}
```

## Try a dry run

```sh
terraform plan
```

Check if the terraform plan looks good

## Execute the plan

```sh
terraform apply
```

Verified the sddc is created

## Add/Remove hosts

Update the "num_host" field in [main.tf](main.tf) to expected number.   
Review and execute the plan

```sh
terraform plan
terraform apply
```

Verified the hosts are added/removed successfully.

## To delete the sddc

```sh
terraform destroy
```

# Testing the Provider

## Set required environment variable

```sh
$ export REFRESH_TOKEN=xxx
$ export ORG_ID=xxxx
$ export TEST_SDDC_ID=xxx
$ make testacc
```

# License

Copyright 2019 VMware, Inc.

The Terraform provider for VMware Cloud on AWS is available under [MPL2.0 license](https://github.com/vmware/terraform-provider-vmc/blob/master/LICENSE).
