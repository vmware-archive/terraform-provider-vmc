# Terraform provider for VMware Cloud on AWS

This is the repository for the Terraform provider for VMware Cloud, which one can use with
Terraform to work with [VMware Cloud on AWS](https://vmc.vmware.com/).

**This repository is Work in Progress and has been made public only for review purposes.**

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Build the Provider

The instructions outlined below to build the provider are specific to Mac OS or Linux OS only.

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

## Configure the provider

The instructions and examples to configure the provider can be found in [examples/README.md](https://github.com/vmware/terraform-provider-vmc/blob/master/examples/README.md)

## Load the provider

```sh
terraform init
```

## Try a dry run

```sh
terraform plan
```

Check if the terraform plan looks good.

## Execute the plan

```sh
   terraform apply
```

Verify the SDDC has been created successfully.

## Check the terraform state 
```sh
terraform show
```

## Add/Remove hosts

Update the "num_host" field in [main.tf](https://github.com/vmware/terraform-provider-vmc/blob/master/examples/main.tf).
Execute the plan

```sh
terraform apply
```

Verify the hosts are added/removed successfully.

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

The Terraform provider for VMware Cloud on AWS is available under [MPL2.0 license](https://github.com/vmware/terraform-provider-vmc/blob/master/LICENSE.txt).
