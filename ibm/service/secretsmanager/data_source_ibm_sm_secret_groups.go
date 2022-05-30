// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
)

func DataSourceIBMSmSecretGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMSmSecretGroupsRead,

		Schema: map[string]*schema.Schema{
			"metadata": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The metadata that describes the resource array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collection_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resources in the resource array.",
						},
						"collection_total": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of elements in the resource array.",
						},
					},
				},
			},
			"resources": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The properties in JSON format to define, such as the name and description. For more information, see the docs: https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-cli-plugin-secrets-manager-cli#secrets-manager-cli-secret-group-create-command",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The v4 UUID that uniquely identifies the secret group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of policy. Allowable values are: rotation",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An extended description of your secret group.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
						},
						"creation_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date the secret group was created. The date format follows RFC 3339.",
						},
						"last_update_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updates when the metadata of the secret group is modified. The date format follows RFC 3339.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MIME type that represents the secret group.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMSmSecretGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	listSecretGroupsOptions := &secretsmanagerv1.ListSecretGroupsOptions{}

	secretGroupDef, response, err := secretsManagerClient.ListSecretGroupsWithContext(context, listSecretGroupsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListSecretGroupsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListSecretGroupsWithContext failed %s\n%s", err, response))
	}

	d.SetId(DataSourceIBMSmSecretGroupsID(d))

	metadata := []map[string]interface{}{}
	if secretGroupDef.Metadata != nil {
		modelMap, err := DataSourceIBMSmSecretGroupsCollectionMetadataToMap(secretGroupDef.Metadata)
		if err != nil {
			return diag.FromErr(err)
		}
		metadata = append(metadata, modelMap)
	}
	if err = d.Set("metadata", metadata); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting metadata %s", err))
	}

	resources := []map[string]interface{}{}
	if secretGroupDef.Resources != nil {
		for _, modelItem := range secretGroupDef.Resources {
			modelMap, err := DataSourceIBMSmSecretGroupsSecretGroupResourceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			resources = append(resources, modelMap)
		}
	}
	if err = d.Set("resources", resources); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resources %s", err))
	}

	return nil
}

// DataSourceIBMSmSecretGroupsID returns a reasonable ID for the list.
func DataSourceIBMSmSecretGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMSmSecretGroupsCollectionMetadataToMap(model *secretsmanagerv1.CollectionMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CollectionType != nil {
		modelMap["collection_type"] = *model.CollectionType
	}
	if model.CollectionTotal != nil {
		modelMap["collection_total"] = *model.CollectionTotal
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretGroupsSecretGroupResourceToMap(model *secretsmanagerv1.SecretGroupResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}
