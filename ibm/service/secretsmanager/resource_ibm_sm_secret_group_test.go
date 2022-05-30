// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMSmSecretGroupBasic(t *testing.T) {
	var conf secretsmanagerv1.SecretGroupResource

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMSmSecretGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMSmSecretGroupConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMSmSecretGroupExists("ibm_sm_secret_group.sm_secret_group", conf),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_sm_secret_group.sm_secret_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMSmSecretGroupConfigBasic() string {
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
	`, )
}

func testAccCheckIBMSmSecretGroupExists(n string, obj secretsmanagerv1.SecretGroupResource) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		secretsManagerClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).SecretsManagerV1()
		if err != nil {
			return err
		}

		getSecretGroupOptions := &secretsmanagerv1.GetSecretGroupOptions{}

		getSecretGroupOptions.SetID(rs.Primary.ID)

		secretGroupResource, _, err := secretsManagerClient.GetSecretGroup(getSecretGroupOptions)
		if err != nil {
			return err
		}

		obj = *secretGroupResource
		return nil
	}
}

func testAccCheckIBMSmSecretGroupDestroy(s *terraform.State) error {
	secretsManagerClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_sm_secret_group" {
			continue
		}

		getSecretGroupOptions := &secretsmanagerv1.GetSecretGroupOptions{}

		getSecretGroupOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := secretsManagerClient.GetSecretGroup(getSecretGroupOptions)

		if err == nil {
			return fmt.Errorf("sm_secret_group still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for sm_secret_group (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
