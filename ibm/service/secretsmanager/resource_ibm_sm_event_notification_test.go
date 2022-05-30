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

func TestAccIBMSmEventNotificationBasic(t *testing.T) {
	var conf secretsmanagerv1.NotificationsSettings
	eventNotificationsInstanceCRN := fmt.Sprintf("tf_event_notifications_instance_crn_%d", acctest.RandIntRange(10, 100))
	eventNotificationsSourceName := fmt.Sprintf("tf_event_notifications_source_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMSmEventNotificationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMSmEventNotificationConfigBasic(eventNotificationsInstanceCRN, eventNotificationsSourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMSmEventNotificationExists("ibm_sm_event_notification.sm_event_notification", conf),
					resource.TestCheckResourceAttr("ibm_sm_event_notification.sm_event_notification", "event_notifications_instance_crn", eventNotificationsInstanceCRN),
					resource.TestCheckResourceAttr("ibm_sm_event_notification.sm_event_notification", "event_notifications_source_name", eventNotificationsSourceName),
				),
			},
		},
	})
}

func TestAccIBMSmEventNotificationAllArgs(t *testing.T) {
	var conf secretsmanagerv1.NotificationsSettings
	eventNotificationsInstanceCRN := fmt.Sprintf("tf_event_notifications_instance_crn_%d", acctest.RandIntRange(10, 100))
	eventNotificationsSourceName := fmt.Sprintf("tf_event_notifications_source_name_%d", acctest.RandIntRange(10, 100))
	eventNotificationsSourceDescription := fmt.Sprintf("tf_event_notifications_source_description_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMSmEventNotificationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMSmEventNotificationConfig(eventNotificationsInstanceCRN, eventNotificationsSourceName, eventNotificationsSourceDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMSmEventNotificationExists("ibm_sm_event_notification.sm_event_notification", conf),
					resource.TestCheckResourceAttr("ibm_sm_event_notification.sm_event_notification", "event_notifications_instance_crn", eventNotificationsInstanceCRN),
					resource.TestCheckResourceAttr("ibm_sm_event_notification.sm_event_notification", "event_notifications_source_name", eventNotificationsSourceName),
					resource.TestCheckResourceAttr("ibm_sm_event_notification.sm_event_notification", "event_notifications_source_description", eventNotificationsSourceDescription),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_sm_event_notification.sm_event_notification",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMSmEventNotificationConfigBasic(eventNotificationsInstanceCRN string, eventNotificationsSourceName string) string {
	return fmt.Sprintf(`

		resource "ibm_sm_event_notification" "sm_event_notification" {
			event_notifications_instance_crn = "%s"
			event_notifications_source_name = "%s"
		}
	`, eventNotificationsInstanceCRN, eventNotificationsSourceName)
}

func testAccCheckIBMSmEventNotificationConfig(eventNotificationsInstanceCRN string, eventNotificationsSourceName string, eventNotificationsSourceDescription string) string {
	return fmt.Sprintf(`

		resource "ibm_sm_event_notification" "sm_event_notification" {
			event_notifications_instance_crn = "%s"
			event_notifications_source_name = "%s"
			event_notifications_source_description = "%s"
		}
	`, eventNotificationsInstanceCRN, eventNotificationsSourceName, eventNotificationsSourceDescription)
}

func testAccCheckIBMSmEventNotificationExists(n string, obj secretsmanagerv1.NotificationsSettings) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		secretsManagerClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).SecretsManagerV1()
		if err != nil {
			return err
		}

		getNotificationsRegistrationOptions := &secretsmanagerv1.GetNotificationsRegistrationOptions{}

		getNotificationsRegistrationOptions.(rs.Primary.ID)

		notificationsSettings, _, err := secretsManagerClient.GetNotificationsRegistration(getNotificationsRegistrationOptions)
		if err != nil {
			return err
		}

		obj = *notificationsSettings
		return nil
	}
}

func testAccCheckIBMSmEventNotificationDestroy(s *terraform.State) error {
	secretsManagerClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_sm_event_notification" {
			continue
		}

		getNotificationsRegistrationOptions := &secretsmanagerv1.GetNotificationsRegistrationOptions{}

		getNotificationsRegistrationOptions.(rs.Primary.ID)

		// Try to find the key
		_, response, err := secretsManagerClient.GetNotificationsRegistration(getNotificationsRegistrationOptions)

		if err == nil {
			return fmt.Errorf("sm_event_notification still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for sm_event_notification (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
