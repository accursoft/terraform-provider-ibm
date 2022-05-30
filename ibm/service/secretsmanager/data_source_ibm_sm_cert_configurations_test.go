// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMSmCertConfigurationsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMSmCertConfigurationsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_sm_cert_configurations.sm_cert_configurations", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_sm_cert_configurations.sm_cert_configurations", "secret_type"),
					resource.TestCheckResourceAttrSet("data.ibm_sm_cert_configurations.sm_cert_configurations", "config_element"),
					resource.TestCheckResourceAttrSet("data.ibm_sm_cert_configurations.sm_cert_configurations", "metadata.#"),
					resource.TestCheckResourceAttrSet("data.ibm_sm_cert_configurations.sm_cert_configurations", "resources.#"),
				),
			},
		},
	})
}

func testAccCheckIBMSmCertConfigurationsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_sm_cert_configurations" "sm_cert_configurations" {
			secret_type = "public_cert"
			config_element = "certificate_authorities"
		}
	`)
}
