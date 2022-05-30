// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMSmSecretsDataSourceBasic(t *testing.T) {
	secretResourceSecretType := "arbitrary"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMSmSecretsDataSourceConfigBasic(secretResourceSecretType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_sm_secrets.sm_secrets", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_sm_secrets.sm_secrets", "metadata.#"),
				),
			},
		},
	})
}

func testAccCheckIBMSmSecretsDataSourceConfigBasic(secretResourceSecretType string) string {
	return fmt.Sprintf(`
		resource "ibm_sm_secret" "sm_secret" {
			secret_type = "%s"
			secret_resource {
				id = "id"
				name = "name"
				description = "description"
				secret_group_id = "secret_group_id"
				labels = [ "labels" ]
				state = 1
				state_description = "Active"
				secret_type = "arbitrary"
				crn = "crn:v1:bluemix:public:secrets-manager:<region>:a/<account-id>:<service-instance>:secret:<secret-id>"
				creation_date = 2018-04-12T23:20:50.520Z
				created_by = "created_by"
				last_update_date = 2018-04-12T23:20:50.520Z
				versions_total = 1
				versions = [ { "key": null } ]
				expiration_date = 2030-04-01T09:30:00.000Z
				payload = "payload"
			}
		}

		data "ibm_sm_secrets" "sm_secrets" {
			depends_on = [
				ibm_sm_secret.sm_secret
			]
		}
	`, secretResourceSecretType)
}
