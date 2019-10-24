provider "vmc" {
  refresh_token = "SPjMe01QUkSx0fulaR7NFBIqQi2MfYQ53VPPNin8jbfXF5qgAg6DmkHdBDzmsiDI"
  //main_token = "12nonEhXDTM2gtXBrVAcm7PMvyhv6t6i6uq9AJksWwoeM4b1Hzh8mtBuVTCgWHbR"
  //staging_refresh_token = "SPjMe01QUkSx0fulaR7NFBIqQi2MfYQ53VPPNin8jbfXF5qgAg6DmkHdBDzmsiDI"


  # for staging environment only
  vmc_url       = "https://stg.skyscraper.vmware.com/vmc/api"
  csp_url       = "https://console-stg.cloud.vmware.com"
}

data "vmc_org" "my_org" {
  id = "05e0a625-3293-41bb-a01f-35e762781c2a"
  //id = "058f47c4-92aa-417f-8747-87f3ed61cb45"



  //stag_id = "05e0a625-3293-41bb-a01f-35e762781c2a"
  //prod_id = "058f47c4-92aa-417f-8747-87f3ed61cb45"
}

data "vmc_connected_accounts" "my_accounts" {
  org_id = "${data.vmc_org.my_org.id}"
}

//data "vmc_customer_subnets" "my_subnets" {
//  org_id               = "${data.vmc_org.my_org.id}"
//  connected_account_id = "${data.vmc_connected_accounts.my_accounts.ids.0}"
//  region               = "us-west-2"
//}

resource "vmc_sddc" "sddc_25" {
  org_id = "${data.vmc_org.my_org.id}"

  storage_capacity    = 0
  sddc_name           = "sumit_sddc_test_real25"
  vpc_cidr            = "10.2.0.0/16"
  num_host            = 3
  provider_type       = "ZEROCLOUD"
  region              = "AP_SOUTHEAST_1"
  vxlan_subnet        = "192.168.1.0/24"
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"

  # sddc_template_id = ""
  deployment_type = "SingleAZ"
  account_link_sddc_config = []
}



resource "vmc_publicips" "IP4" {
  org_id = "${data.vmc_org.my_org.id}"
  sddc_id = "${vmc_sddc.sddc_25.id}"
  private_ip = "10.2.33.45"
  name = "workload"
}


