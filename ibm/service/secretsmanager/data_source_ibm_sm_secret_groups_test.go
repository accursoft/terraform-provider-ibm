// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMSmSecretGroupsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMSmSecretGroupsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_sm_secret_groups.sm_secret_groups", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_sm_secret_groups.sm_secret_groups", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_sm_secret_groups.sm_secret_groups", "metadata.#"),
					resource.TestCheckResourceAttrSet("data.ibm_sm_secret_groups.sm_secret_groups", "resources.#"),
				),
			},
		},
	})
}

func testAccCheckIBMSmSecretGroupsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		resource "ibm_sm_secret_group" "sm_secret_group" {
			metadata {
				collection_type = "application/vnd.ibm.secrets-manager.config+json"
				collection_total = 1
			}
			resources {
				id = "bc656587-8fda-4d05-9ad8-b1de1ec7e712"
				name = "my-secret-group"
				description = "Extended description for this group."
				creation_date = 2018-04-12T23:20:50.520Z
				last_update_date = 2018-05-12T23:20:50.520Z
				type = "application/vnd.ibm.secrets-manager.secret.group+json"
			}
		}

		data "ibm_sm_secret_groups" "sm_secret_groups" {
			id = "id"
		}
	`)
}

