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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
)

func ResourceIBMSmSecretGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext:   ResourceIBMSmSecretGroupCreate,
		ReadContext:     ResourceIBMSmSecretGroupRead,
		UpdateContext:   ResourceIBMSmSecretGroupUpdate,
		DeleteContext:   ResourceIBMSmSecretGroupDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"metadata": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The metadata that describes the resource array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collection_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of resources in the resource array.",
						},
						"collection_total": &schema.Schema{
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The number of elements in the resource array.",
						},
					},
				},
			},
			"resources": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The properties in JSON format to define, such as the name and description. For more information, see the docs: https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-cli-plugin-secrets-manager-cli#secrets-manager-cli-secret-group-create-command",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The v4 UUID that uniquely identifies the secret group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of policy. Allowable values are: rotation",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An extended description of your secret group.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
						},
						"creation_date": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The date the secret group was created. The date format follows RFC 3339.",
						},
						"last_update_date": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Updates when the metadata of the secret group is modified. The date format follows RFC 3339.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The MIME type that represents the secret group.",
						},
					},
				},
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
	}
}

func ResourceIBMSmSecretGroupCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createSecretGroupOptions := &secretsmanagerv1.CreateSecretGroupOptions{}

	metadataModel, err := ResourceIBMSmSecretGroupMapToCollectionMetadata(d.Get("metadata.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createSecretGroupOptions.SetMetadata(metadataModel)
	var resources []secretsmanagerv1.SecretGroupResource
	for _, e := range d.Get("resources").([]interface{}) {
		value := e.(map[string]interface{})
		resourcesItem, err := ResourceIBMSmSecretGroupMapToSecretGroupResource(value)
		if err != nil {
			return diag.FromErr(err)
		}
		resources = append(resources, *resourcesItem)
	}
	createSecretGroupOptions.SetResources(resources)

	secretGroupDef, response, err := secretsManagerClient.CreateSecretGroupWithContext(context, createSecretGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateSecretGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateSecretGroupWithContext failed %s\n%s", err, response))
	}

	d.SetId(*secretGroupDef.ID)

	return ResourceIBMSmSecretGroupRead(context, d, meta)
}

func ResourceIBMSmSecretGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getSecretGroupOptions := &secretsmanagerv1.GetSecretGroupOptions{}

	getSecretGroupOptions.SetID(d.Id())

	secretGroupDef, response, err := secretsManagerClient.GetSecretGroupWithContext(context, getSecretGroupOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSecretGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSecretGroupWithContext failed %s\n%s", err, response))
	}

	// TODO: handle argument of type CollectionMetadata
	// TODO: handle argument of type []interface{}
	if err = d.Set("name", secretGroupDef.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("description", secretGroupDef.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if err = d.Set("creation_date", flex.DateTimeToString(secretGroupDef.CreationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting creation_date: %s", err))
	}
	if err = d.Set("last_update_date", flex.DateTimeToString(secretGroupDef.LastUpdateDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_update_date: %s", err))
	}
	if err = d.Set("type", secretGroupDef.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	return nil
}

func ResourceIBMSmSecretGroupUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateSecretGroupMetadataOptions := &secretsmanagerv1.UpdateSecretGroupMetadataOptions{}

	updateSecretGroupMetadataOptions.SetID(d.Id())

	hasChange := false

	if d.HasChange("metadata") || d.HasChange("resources") {
		metadata, err := ResourceIBMSmSecretGroupMapToCollectionMetadata(d.Get("metadata.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateSecretGroupMetadataOptions.SetMetadata(metadata)
		// TODO: handle Resources of type TypeList -- not primitive, not model
		hasChange = true
	}

	if hasChange {
		_, response, err := secretsManagerClient.UpdateSecretGroupMetadataWithContext(context, updateSecretGroupMetadataOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateSecretGroupMetadataWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateSecretGroupMetadataWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIBMSmSecretGroupRead(context, d, meta)
}

func ResourceIBMSmSecretGroupDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteSecretGroupOptions := &secretsmanagerv1.DeleteSecretGroupOptions{}

	deleteSecretGroupOptions.SetID(d.Id())

	response, err := secretsManagerClient.DeleteSecretGroupWithContext(context, deleteSecretGroupOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteSecretGroupWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteSecretGroupWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func ResourceIBMSmSecretGroupMapToCollectionMetadata(modelMap map[string]interface{}) (*secretsmanagerv1.CollectionMetadata, error) {
	model := &secretsmanagerv1.CollectionMetadata{}
	model.CollectionType = core.StringPtr(modelMap["collection_type"].(string))
	model.CollectionTotal = core.Int64Ptr(int64(modelMap["collection_total"].(int)))
	return model, nil
}

func ResourceIBMSmSecretGroupMapToSecretGroupResource(modelMap map[string]interface{}) (*secretsmanagerv1.SecretGroupResource, error) {
	model := &secretsmanagerv1.SecretGroupResource{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["creation_date"] != nil {
	
	}
	if modelMap["last_update_date"] != nil {
	
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func ResourceIBMSmSecretGroupCollectionMetadataToMap(model *secretsmanagerv1.CollectionMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["collection_type"] = model.CollectionType
	modelMap["collection_total"] = flex.IntValue(model.CollectionTotal)
	return modelMap, nil
}

func ResourceIBMSmSecretGroupSecretGroupResourceToMap(model *secretsmanagerv1.SecretGroupResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}
