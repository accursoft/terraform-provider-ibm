// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMSmSecretBasic(t *testing.T) {
	var conf secretsmanagerv1.SecretResource
	secretType := "arbitrary"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMSmSecretDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMSmSecretConfigBasic(secretType),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMSmSecretExists("ibm_sm_secret.sm_secret", conf),
					resource.TestCheckResourceAttr("ibm_sm_secret.sm_secret", "secret_type", secretType),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_sm_secret.sm_secret",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMSmSecretConfigBasic(secretType string) string {
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
	`, secretType)
}

func testAccCheckIBMSmSecretExists(n string, obj secretsmanagerv1.SecretResource) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		secretsManagerClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).SecretsManagerV1()
		if err != nil {
			return err
		}

		getSecretOptions := &secretsmanagerv1.GetSecretOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getSecretOptions.SetSecretType(parts[0])
		getSecretOptions.SetID(parts[1])

		secretResourceIntf, _, err := secretsManagerClient.GetSecret(getSecretOptions)
		if err != nil {
			return err
		}

		secretResource := secretResourceIntf.(*secretsmanagerv1.SecretResource)
		obj = *secretResource
		return nil
	}
}

func testAccCheckIBMSmSecretDestroy(s *terraform.State) error {
	secretsManagerClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_sm_secret" {
			continue
		}

		getSecretOptions := &secretsmanagerv1.GetSecretOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getSecretOptions.SetSecretType(parts[0])
		getSecretOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := secretsManagerClient.GetSecret(getSecretOptions)

		if err == nil {
			return fmt.Errorf("sm_secret still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for sm_secret (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
