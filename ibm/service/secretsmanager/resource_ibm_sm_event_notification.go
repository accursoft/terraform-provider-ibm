// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
)

func ResourceIBMSmEventNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext:   ResourceIBMSmEventNotificationCreate,
		ReadContext:     ResourceIBMSmEventNotificationRead,
		DeleteContext:   ResourceIBMSmEventNotificationDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"event_notifications_instance_crn": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Cloud Resource Name (CRN) of the connected Event Notifications instance.",
			},
			"event_notifications_source_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name that is displayed as a source in your Event Notifications instance.",
			},
			"event_notifications_source_description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional description for the source in your Event Notifications instance.",
			},
		},
	}
}

func ResourceIBMSmEventNotificationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createNotificationsRegistrationOptions := &secretsmanagerv1.CreateNotificationsRegistrationOptions{}

	createNotificationsRegistrationOptions.SetEventNotificationsInstanceCRN(d.Get("event_notifications_instance_crn").(string))
	createNotificationsRegistrationOptions.SetEventNotificationsSourceName(d.Get("event_notifications_source_name").(string))
	if _, ok := d.GetOk("event_notifications_source_description"); ok {
		createNotificationsRegistrationOptions.SetEventNotificationsSourceDescription(d.Get("event_notifications_source_description").(string))
	}

	getNotificationsSettings, response, err := secretsManagerClient.CreateNotificationsRegistrationWithContext(context, createNotificationsRegistrationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateNotificationsRegistrationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateNotificationsRegistrationWithContext failed %s\n%s", err, response))
	}

	d.SetId(*getNotificationsSettings.EventNotificationsInstanceCRN)

	return ResourceIBMSmEventNotificationRead(context, d, meta)
}

func ResourceIBMSmEventNotificationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getNotificationsRegistrationOptions := &secretsmanagerv1.GetNotificationsRegistrationOptions{}


	getNotificationsSettings, response, err := secretsManagerClient.GetNotificationsRegistrationWithContext(context, getNotificationsRegistrationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetNotificationsRegistrationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetNotificationsRegistrationWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("event_notifications_source_name", getNotificationsRegistrationOptions.EventNotificationsSourceName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting event_notifications_source_name: %s", err))
	}
	if err = d.Set("event_notifications_source_description", getNotificationsRegistrationOptions.EventNotificationsSourceDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting event_notifications_source_description: %s", err))
	}
	if err = d.Set("event_notifications_instance_crn", getNotificationsSettings.EventNotificationsInstanceCRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting event_notifications_instance_crn: %s", err))
	}

	return nil
}

func ResourceIBMSmEventNotificationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteNotificationsRegistrationOptions := &secretsmanagerv1.DeleteNotificationsRegistrationOptions{}


	response, err := secretsManagerClient.DeleteNotificationsRegistrationWithContext(context, deleteNotificationsRegistrationOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteNotificationsRegistrationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteNotificationsRegistrationWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
