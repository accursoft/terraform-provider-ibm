// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
)

func TestAccIBMSmCertConfigurationBasic(t *testing.T) {
	var conf secretsmanagerv1.ConfigElementDef
	secretType := "public_cert"
	configElement := "certificate_authorities"
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	typeVar := "letsencrypt"
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	typeVarUpdate := "certificate_template"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMSmCertConfigurationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMSmCertConfigurationConfigBasic(secretType, configElement, name, typeVar),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMSmCertConfigurationExists("ibm_sm_cert_configuration.sm_cert_configuration", conf),
					resource.TestCheckResourceAttr("ibm_sm_cert_configuration.sm_cert_configuration", "secret_type", secretType),
					resource.TestCheckResourceAttr("ibm_sm_cert_configuration.sm_cert_configuration", "config_element", configElement),
					resource.TestCheckResourceAttr("ibm_sm_cert_configuration.sm_cert_configuration", "name", name),
					resource.TestCheckResourceAttr("ibm_sm_cert_configuration.sm_cert_configuration", "type", typeVar),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMSmCertConfigurationConfigBasic(secretType, configElement, nameUpdate, typeVarUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_sm_cert_configuration.sm_cert_configuration", "secret_type", secretType),
					resource.TestCheckResourceAttr("ibm_sm_cert_configuration.sm_cert_configuration", "config_element", configElement),
					resource.TestCheckResourceAttr("ibm_sm_cert_configuration.sm_cert_configuration", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_sm_cert_configuration.sm_cert_configuration", "type", typeVarUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_sm_cert_configuration.sm_cert_configuration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMSmCertConfigurationConfigBasic(secretType string, configElement string, name string, typeVar string) string {
	return fmt.Sprintf(`

		resource "ibm_sm_cert_configuration" "sm_cert_configuration" {
			secret_type = "%s"
			config_element = "%s"
			name = "%s"
			type = "%s"
			config {
				private_key = "private_key"
			}
		}
	`, secretType, configElement, name, typeVar)
}

func testAccCheckIBMSmCertConfigurationExists(n string, obj secretsmanagerv1.ConfigElementDef) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		secretsManagerClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).SecretsManagerV1()
		if err != nil {
			return err
		}

		getConfigElementOptions := &secretsmanagerv1.GetConfigElementOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getConfigElementOptions.SetSecretType(parts[0])
		getConfigElementOptions.SetConfigElement(parts[1])
		getConfigElementOptions.SetConfigName(parts[2])

		configElementDef, _, err := secretsManagerClient.GetConfigElement(getConfigElementOptions)
		if err != nil {
			return err
		}

		obj = *configElementDef
		return nil
	}
}

func testAccCheckIBMSmCertConfigurationDestroy(s *terraform.State) error {
	secretsManagerClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_sm_cert_configuration" {
			continue
		}

		getConfigElementOptions := &secretsmanagerv1.GetConfigElementOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getConfigElementOptions.SetSecretType(parts[0])
		getConfigElementOptions.SetConfigElement(parts[1])
		getConfigElementOptions.SetConfigName(parts[2])

		// Try to find the key
		_, response, err := secretsManagerClient.GetConfigElement(getConfigElementOptions)

		if err == nil {
			return fmt.Errorf("sm_cert_configuration still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for sm_cert_configuration (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
