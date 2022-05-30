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

func DataSourceIBMSmSecrets() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMSmSecretsRead,

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
				Description: "A collection of resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The v4 UUID that uniquely identifies the secret.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A human-readable alias to assign to your secret.To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.",
						},
						"secret_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The v4 UUID that uniquely identifies the secret group to assign to this secret.If you omit this parameter, your secret is assigned to the `default` secret group.",
						},
						"labels": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Labels that you can use to filter for secrets in your instance.Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"state": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,  Suspended = 2, Deactivated = 3, and Destroyed = 5 values.",
						},
						"state_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A text representation of the secret state.",
						},
						"secret_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret type.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.",
						},
						"creation_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date the secret was created. The date format follows RFC 3339.",
						},
						"created_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the entity that created the secret.",
						},
						"last_update_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updates when the actual secret is modified. The date format follows RFC 3339.",
						},
						"versions_total": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of versions that are associated with a secret.",
						},
						"versions": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An array that contains metadata for each secret version. For more information on the metadata properties, see [Get secret version metadata](#get-secret-version-metadata).",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"expiration_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date the secret material expires. The date format follows RFC 3339.You can set an expiration date on supported secret types at their creation. If you create a secret without specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the following secret types:- `arbitrary`- `username_password`.",
						},
						"payload": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The new secret data to assign to the secret.",
						},
						"secret_data": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The data that is associated with the secret version.The data object contains the field `payload`.",
						},
						"username": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username to assign to this secret.",
						},
						"password": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The password to assign to this secret.",
						},
						"next_rotation_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date that the secret is scheduled for automatic rotation.The service automatically creates a new version of the secret on its next rotation date. This field exists only for secrets that can be auto-rotated and have an existing rotation policy.",
						},
						"ttl": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The time-to-live (TTL) or lease duration to assign to generated credentials.For `iam_credentials` secrets, the TTL defines for how long each generated API key remains valid. The value can be either an integer that specifies the number of seconds, or the string representation of a duration, such as `120m` or `24h`.Minimum duration is 1 minute. Maximum is 90 days.",
						},
						"access_groups": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The access groups that define the capabilities of the service ID and API key that are generated for an`iam_credentials` secret. If you prefer to use an existing service ID that is already assigned the access policies that you require, you can omit this parameter and use the `service_id` field instead.**Tip:** To list the access groups that are available in an account, you can use the [IAM Access Groups API](https://cloud.ibm.com/apidocs/iam-access-groups#list-access-groups). To find the ID of an access group in the console, go to **Manage > Access (IAM) > Access groups**. Select the access group to inspect, and click **Details** to view its ID.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"api_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API key that is generated for this secret.After the secret reaches the end of its lease (see the `ttl` field), the API key is deleted automatically. If you want to continue to use the same API key for future read operations, see the `reuse_api_key` field.",
						},
						"api_key_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the API key that is generated for this secret.",
						},
						"service_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "(IAM credentials) The service ID under which the API key is created. To have Secrets Manager generate a new service ID, omit this option and include 'access_groups'.",
						},
						"service_id_is_static": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether an `iam_credentials` secret was created with a static service ID.If `true`, the service ID for the secret was provided by the user at secret creation. If `false`, the service ID was generated by Secrets Manager.",
						},
						"reuse_api_key": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "(IAM credentials) Reuse the service ID and API key for future read operations.",
						},
						"certificate": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The contents of your certificate.",
						},
						"private_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "(Optional) The private key to associate with the certificate.",
						},
						"intermediate": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "(Optional) The intermediate certificate to associate with the root certificate.",
						},
						"serial_number": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique serial number that was assigned to the certificate by the issuing certificate authority.",
						},
						"algorithm": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign the certificate.",
						},
						"key_algorithm": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identifier for the cryptographic algorithm that was used to generate the public and private keys that are associated with the certificate.",
						},
						"issuer": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The distinguished name that identifies the entity that signed and issued the certificate.",
						},
						"validity": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"not_before": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that the certificate validity period begins.",
									},
									"not_after": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that the certificate validity period ends.",
									},
								},
							},
						},
						"common_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The fully qualified domain name or host domain name that is defined for the certificate.",
						},
						"intermediate_included": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the certificate was imported with an associated intermediate certificate.",
						},
						"private_key_included": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the certificate was imported with an associated private key.",
						},
						"alt_names": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The alternative names that are defined for the certificate.For public certificates, this value is provided as an array of strings. For private certificates, this value is provided as a comma-delimited list (string). In the API response, this value is returned as an array of strings for all the types of certificate secrets.",
						},
						"bundle_certs": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determines whether your issued certificate is bundled with intermediate certificates.Set to `false` for the certificate file to contain only the issued certificate.",
						},
						"ca": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the certificate authority configuration.",
						},
						"dns": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the DNS provider configuration.",
						},
						"rotation": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_rotate": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Determines whether Secrets Manager rotates your certificate automatically.For public certificates, if `auto_rotate` is set to `true` the service reorders your certificate 31 days before it expires. For private certificates, the certificate is rotated according to the time interval specified in the `interval` and `unit` fields.To access the previous version of the certificate, you can use the[Get a version of a secret](#get-secret-version) method.",
									},
									"rotate_keys": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Determines whether Secrets Manager rotates the private key for your certificate automatically.If set to `true`, the service generates and stores a new private key for your rotated certificate.**Note:** Use this field only for public certificates. It is ignored for private certificates.",
									},
									"interval": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Used together with the `unit` field to specify the rotation interval. The minimum interval is one day, and the maximum interval is 3 years (1095 days). Required in case `auto_rotate` is set to `true`.**Note:** Use this field only for private certificates. It is ignored for public certificates.",
									},
									"unit": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The time unit of the rotation interval.**Note:** Use this field only for private certificates. It is ignored for public certificates.",
									},
								},
							},
						},
						"issuance_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Issuance information that is associated with your certificate.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ordered_on": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date the certificate was ordered. The date format follows RFC 3339.",
									},
									"error_code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A code that identifies an issuance error.This field, along with `error_message`, is returned when Secrets Manager successfully processes your request, but a certificate is unable to be issued by the certificate authority.",
									},
									"error_message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A human-readable message that provides details about the issuance error.",
									},
									"bundle_certs": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the issued certificate is bundled with intermediate certificates.",
									},
									"state": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,  Suspended = 2, Deactivated = 3, and Destroyed = 5 values.",
									},
									"state_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A text representation of the secret state.",
									},
									"auto_rotated": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the issued certificate is configured with an automatic rotation policy.",
									},
									"ca": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name that was assigned to the certificate authority configuration.",
									},
									"dns": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name that was assigned to the DNS provider configuration.",
									},
								},
							},
						},
						"certificate_template": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the certificate template.",
						},
						"certificate_authority": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The intermediate certificate authority that signed this certificate.",
						},
						"ip_sans": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.",
						},
						"uri_sans": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.",
						},
						"other_sans": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"format": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The format of the returned data.",
						},
						"private_key_format": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The format of the generated private key.",
						},
						"exclude_cn_from_sans": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Controls whether the common name is excluded from Subject Alternative Names (SANs). If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if the common name is not a hostname or an email address, but is instead a human-readable identifier.",
						},
						"revocation_time": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timestamp of the certificate revocation.",
						},
						"revocation_time_rfc3339": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the certificate was revoked. The date format follows RFC 3339.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMSmSecretsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	listAllSecretsOptions := &secretsmanagerv1.ListAllSecretsOptions{}

	listSecrets, response, err := secretsManagerClient.ListAllSecretsWithContext(context, listAllSecretsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListAllSecretsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListAllSecretsWithContext failed %s\n%s", err, response))
	}

	d.SetId(DataSourceIBMSmSecretsID(d))

	metadata := []map[string]interface{}{}
	if listSecrets.Metadata != nil {
		modelMap, err := DataSourceIBMSmSecretsCollectionMetadataToMap(listSecrets.Metadata)
		if err != nil {
			return diag.FromErr(err)
		}
		metadata = append(metadata, modelMap)
	}
	if err = d.Set("metadata", metadata); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting metadata %s", err))
	}

	resources := []map[string]interface{}{}
	if listSecrets.Resources != nil {
		for _, modelItem := range listSecrets.Resources { 
			modelMap, err := DataSourceIBMSmSecretsSecretResourceToMap(modelItem)
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

// DataSourceIBMSmSecretsID returns a reasonable ID for the list.
func DataSourceIBMSmSecretsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMSmSecretsCollectionMetadataToMap(model *secretsmanagerv1.CollectionMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CollectionType != nil {
		modelMap["collection_type"] = *model.CollectionType
	}
	if model.CollectionTotal != nil {
		modelMap["collection_total"] = *model.CollectionTotal
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretsSecretResourceToMap(model secretsmanagerv1.SecretResourceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv1.ArbitrarySecretResource); ok {
		return DataSourceIBMSmSecretsArbitrarySecretResourceToMap(model.(*secretsmanagerv1.ArbitrarySecretResource))
	} else if _, ok := model.(*secretsmanagerv1.UsernamePasswordSecretResource); ok {
		return DataSourceIBMSmSecretsUsernamePasswordSecretResourceToMap(model.(*secretsmanagerv1.UsernamePasswordSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.IamCredentialsSecretResource); ok {
		return DataSourceIBMSmSecretsIamCredentialsSecretResourceToMap(model.(*secretsmanagerv1.IamCredentialsSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.CertificateSecretResource); ok {
		return DataSourceIBMSmSecretsCertificateSecretResourceToMap(model.(*secretsmanagerv1.CertificateSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.PublicCertificateSecretResource); ok {
		return DataSourceIBMSmSecretsPublicCertificateSecretResourceToMap(model.(*secretsmanagerv1.PublicCertificateSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.PrivateCertificateSecretResource); ok {
		return DataSourceIBMSmSecretsPrivateCertificateSecretResourceToMap(model.(*secretsmanagerv1.PrivateCertificateSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.KvSecretResource); ok {
		return DataSourceIBMSmSecretsKvSecretResourceToMap(model.(*secretsmanagerv1.KvSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.SecretResource); ok {
		modelMap := make(map[string]interface{})
		model := model.(*secretsmanagerv1.SecretResource)
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.Description != nil {
			modelMap["description"] = *model.Description
		}
		if model.SecretGroupID != nil {
			modelMap["secret_group_id"] = *model.SecretGroupID
		}
		if model.Labels != nil {
			modelMap["labels"] = model.Labels
		}
		if model.State != nil {
			modelMap["state"] = *model.State
		}
		if model.StateDescription != nil {
			modelMap["state_description"] = *model.StateDescription
		}
		if model.SecretType != nil {
			modelMap["secret_type"] = *model.SecretType
		}
		if model.CRN != nil {
			modelMap["crn"] = *model.CRN
		}
		if model.CreationDate != nil {
			modelMap["creation_date"] = model.CreationDate.String()
		}
		if model.CreatedBy != nil {
			modelMap["created_by"] = *model.CreatedBy
		}
		if model.LastUpdateDate != nil {
			modelMap["last_update_date"] = model.LastUpdateDate.String()
		}
		if model.VersionsTotal != nil {
			modelMap["versions_total"] = *model.VersionsTotal
		}
		if model.Versions != nil {
		}
		if model.ExpirationDate != nil {
			modelMap["expiration_date"] = model.ExpirationDate.String()
		}
		if model.Payload != nil {
			modelMap["payload"] = *model.Payload
		}
		if model.SecretData != nil {
		}
		if model.Username != nil {
			modelMap["username"] = *model.Username
		}
		if model.Password != nil {
			modelMap["password"] = *model.Password
		}
		if model.NextRotationDate != nil {
			modelMap["next_rotation_date"] = model.NextRotationDate.String()
		}
		if model.TTL != nil {
		}
		if model.AccessGroups != nil {
			modelMap["access_groups"] = model.AccessGroups
		}
		if model.APIKey != nil {
			modelMap["api_key"] = *model.APIKey
		}
		if model.APIKeyID != nil {
			modelMap["api_key_id"] = *model.APIKeyID
		}
		if model.ServiceID != nil {
			modelMap["service_id"] = *model.ServiceID
		}
		if model.ServiceIDIsStatic != nil {
			modelMap["service_id_is_static"] = *model.ServiceIDIsStatic
		}
		if model.ReuseAPIKey != nil {
			modelMap["reuse_api_key"] = *model.ReuseAPIKey
		}
		if model.Certificate != nil {
			modelMap["certificate"] = *model.Certificate
		}
		if model.PrivateKey != nil {
			modelMap["private_key"] = *model.PrivateKey
		}
		if model.Intermediate != nil {
			modelMap["intermediate"] = *model.Intermediate
		}
		if model.SerialNumber != nil {
			modelMap["serial_number"] = *model.SerialNumber
		}
		if model.Algorithm != nil {
			modelMap["algorithm"] = *model.Algorithm
		}
		if model.KeyAlgorithm != nil {
			modelMap["key_algorithm"] = *model.KeyAlgorithm
		}
		if model.Issuer != nil {
			modelMap["issuer"] = *model.Issuer
		}
		if model.Validity != nil {
			validityMap, err := DataSourceIBMSmSecretsCertificateValidityToMap(model.Validity)
			if err != nil {
				return modelMap, err
			}
			modelMap["validity"] = []map[string]interface{}{validityMap}
		}
		if model.CommonName != nil {
			modelMap["common_name"] = *model.CommonName
		}
		if model.IntermediateIncluded != nil {
			modelMap["intermediate_included"] = *model.IntermediateIncluded
		}
		if model.PrivateKeyIncluded != nil {
			modelMap["private_key_included"] = *model.PrivateKeyIncluded
		}
		if model.AltNames != nil {
		}
		if model.BundleCerts != nil {
			modelMap["bundle_certs"] = *model.BundleCerts
		}
		if model.Ca != nil {
			modelMap["ca"] = *model.Ca
		}
		if model.DNS != nil {
			modelMap["dns"] = *model.DNS
		}
		if model.Rotation != nil {
			rotationMap, err := DataSourceIBMSmSecretsRotationToMap(model.Rotation)
			if err != nil {
				return modelMap, err
			}
			modelMap["rotation"] = []map[string]interface{}{rotationMap}
		}
		if model.IssuanceInfo != nil {
			issuanceInfoMap, err := DataSourceIBMSmSecretsIssuanceInfoToMap(model.IssuanceInfo)
			if err != nil {
				return modelMap, err
			}
			modelMap["issuance_info"] = []map[string]interface{}{issuanceInfoMap}
		}
		if model.CertificateTemplate != nil {
			modelMap["certificate_template"] = *model.CertificateTemplate
		}
		if model.CertificateAuthority != nil {
			modelMap["certificate_authority"] = *model.CertificateAuthority
		}
		if model.IPSans != nil {
			modelMap["ip_sans"] = *model.IPSans
		}
		if model.URISans != nil {
			modelMap["uri_sans"] = *model.URISans
		}
		if model.OtherSans != nil {
			modelMap["other_sans"] = model.OtherSans
		}
		if model.Format != nil {
			modelMap["format"] = *model.Format
		}
		if model.PrivateKeyFormat != nil {
			modelMap["private_key_format"] = *model.PrivateKeyFormat
		}
		if model.ExcludeCnFromSans != nil {
			modelMap["exclude_cn_from_sans"] = *model.ExcludeCnFromSans
		}
		if model.RevocationTime != nil {
			modelMap["revocation_time"] = *model.RevocationTime
		}
		if model.RevocationTimeRfc3339 != nil {
			modelMap["revocation_time_rfc3339"] = model.RevocationTimeRfc3339.String()
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv1.SecretResourceIntf subtype encountered")
	}
}

func DataSourceIBMSmSecretsCertificateValidityToMap(model *secretsmanagerv1.CertificateValidity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NotBefore != nil {
		modelMap["not_before"] = model.NotBefore.String()
	}
	if model.NotAfter != nil {
		modelMap["not_after"] = model.NotAfter.String()
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretsRotationToMap(model *secretsmanagerv1.Rotation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotate != nil {
		modelMap["auto_rotate"] = *model.AutoRotate
	}
	if model.RotateKeys != nil {
		modelMap["rotate_keys"] = *model.RotateKeys
	}
	if model.Interval != nil {
		modelMap["interval"] = *model.Interval
	}
	if model.Unit != nil {
		modelMap["unit"] = *model.Unit
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretsIssuanceInfoToMap(model *secretsmanagerv1.IssuanceInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.OrderedOn != nil {
		modelMap["ordered_on"] = model.OrderedOn.String()
	}
	if model.ErrorCode != nil {
		modelMap["error_code"] = *model.ErrorCode
	}
	if model.ErrorMessage != nil {
		modelMap["error_message"] = *model.ErrorMessage
	}
	if model.BundleCerts != nil {
		modelMap["bundle_certs"] = *model.BundleCerts
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.AutoRotated != nil {
		modelMap["auto_rotated"] = *model.AutoRotated
	}
	if model.Ca != nil {
		modelMap["ca"] = *model.Ca
	}
	if model.DNS != nil {
		modelMap["dns"] = *model.DNS
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretsArbitrarySecretResourceToMap(model *secretsmanagerv1.ArbitrarySecretResource) (map[string]interface{}, error) {
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
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.Versions != nil {
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.Payload != nil {
		modelMap["payload"] = *model.Payload
	}
	if model.SecretData != nil {
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretsUsernamePasswordSecretResourceToMap(model *secretsmanagerv1.UsernamePasswordSecretResource) (map[string]interface{}, error) {
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
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.Versions != nil {
	}
	if model.Username != nil {
		modelMap["username"] = *model.Username
	}
	if model.Password != nil {
		modelMap["password"] = *model.Password
	}
	if model.SecretData != nil {
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.NextRotationDate != nil {
		modelMap["next_rotation_date"] = model.NextRotationDate.String()
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretsIamCredentialsSecretResourceToMap(model *secretsmanagerv1.IamCredentialsSecretResource) (map[string]interface{}, error) {
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
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.Versions != nil {
	}
	if model.TTL != nil {
	}
	if model.AccessGroups != nil {
		modelMap["access_groups"] = model.AccessGroups
	}
	if model.APIKey != nil {
		modelMap["api_key"] = *model.APIKey
	}
	if model.APIKeyID != nil {
		modelMap["api_key_id"] = *model.APIKeyID
	}
	if model.ServiceID != nil {
		modelMap["service_id"] = *model.ServiceID
	}
	if model.ServiceIDIsStatic != nil {
		modelMap["service_id_is_static"] = *model.ServiceIDIsStatic
	}
	if model.ReuseAPIKey != nil {
		modelMap["reuse_api_key"] = *model.ReuseAPIKey
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretsCertificateSecretResourceToMap(model *secretsmanagerv1.CertificateSecretResource) (map[string]interface{}, error) {
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
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.Versions != nil {
	}
	if model.Certificate != nil {
		modelMap["certificate"] = *model.Certificate
	}
	if model.PrivateKey != nil {
		modelMap["private_key"] = *model.PrivateKey
	}
	if model.Intermediate != nil {
		modelMap["intermediate"] = *model.Intermediate
	}
	if model.SecretData != nil {
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = *model.SerialNumber
	}
	if model.Algorithm != nil {
		modelMap["algorithm"] = *model.Algorithm
	}
	if model.KeyAlgorithm != nil {
		modelMap["key_algorithm"] = *model.KeyAlgorithm
	}
	if model.Issuer != nil {
		modelMap["issuer"] = *model.Issuer
	}
	if model.Validity != nil {
		validityMap, err := DataSourceIBMSmSecretsCertificateValidityToMap(model.Validity)
		if err != nil {
			return modelMap, err
		}
		modelMap["validity"] = []map[string]interface{}{validityMap}
	}
	if model.CommonName != nil {
		modelMap["common_name"] = *model.CommonName
	}
	if model.IntermediateIncluded != nil {
		modelMap["intermediate_included"] = *model.IntermediateIncluded
	}
	if model.PrivateKeyIncluded != nil {
		modelMap["private_key_included"] = *model.PrivateKeyIncluded
	}
	if model.AltNames != nil {
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretsPublicCertificateSecretResourceToMap(model *secretsmanagerv1.PublicCertificateSecretResource) (map[string]interface{}, error) {
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
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.Versions != nil {
	}
	if model.Issuer != nil {
		modelMap["issuer"] = *model.Issuer
	}
	if model.BundleCerts != nil {
		modelMap["bundle_certs"] = *model.BundleCerts
	}
	if model.Ca != nil {
		modelMap["ca"] = *model.Ca
	}
	if model.DNS != nil {
		modelMap["dns"] = *model.DNS
	}
	if model.Algorithm != nil {
		modelMap["algorithm"] = *model.Algorithm
	}
	if model.KeyAlgorithm != nil {
		modelMap["key_algorithm"] = *model.KeyAlgorithm
	}
	if model.AltNames != nil {
	}
	if model.CommonName != nil {
		modelMap["common_name"] = *model.CommonName
	}
	if model.PrivateKeyIncluded != nil {
		modelMap["private_key_included"] = *model.PrivateKeyIncluded
	}
	if model.IntermediateIncluded != nil {
		modelMap["intermediate_included"] = *model.IntermediateIncluded
	}
	if model.Rotation != nil {
		rotationMap, err := DataSourceIBMSmSecretsRotationToMap(model.Rotation)
		if err != nil {
			return modelMap, err
		}
		modelMap["rotation"] = []map[string]interface{}{rotationMap}
	}
	if model.IssuanceInfo != nil {
		issuanceInfoMap, err := DataSourceIBMSmSecretsIssuanceInfoToMap(model.IssuanceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["issuance_info"] = []map[string]interface{}{issuanceInfoMap}
	}
	if model.Validity != nil {
		validityMap, err := DataSourceIBMSmSecretsCertificateValidityToMap(model.Validity)
		if err != nil {
			return modelMap, err
		}
		modelMap["validity"] = []map[string]interface{}{validityMap}
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = *model.SerialNumber
	}
	if model.SecretData != nil {
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretsPrivateCertificateSecretResourceToMap(model *secretsmanagerv1.PrivateCertificateSecretResource) (map[string]interface{}, error) {
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
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.Versions != nil {
	}
	if model.CertificateTemplate != nil {
		modelMap["certificate_template"] = *model.CertificateTemplate
	}
	if model.CertificateAuthority != nil {
		modelMap["certificate_authority"] = *model.CertificateAuthority
	}
	if model.CommonName != nil {
		modelMap["common_name"] = *model.CommonName
	}
	if model.AltNames != nil {
	}
	if model.IPSans != nil {
		modelMap["ip_sans"] = *model.IPSans
	}
	if model.URISans != nil {
		modelMap["uri_sans"] = *model.URISans
	}
	if model.OtherSans != nil {
		modelMap["other_sans"] = model.OtherSans
	}
	if model.TTL != nil {
	}
	if model.Format != nil {
		modelMap["format"] = *model.Format
	}
	if model.PrivateKeyFormat != nil {
		modelMap["private_key_format"] = *model.PrivateKeyFormat
	}
	if model.ExcludeCnFromSans != nil {
		modelMap["exclude_cn_from_sans"] = *model.ExcludeCnFromSans
	}
	if model.Rotation != nil {
		rotationMap, err := DataSourceIBMSmSecretsRotationToMap(model.Rotation)
		if err != nil {
			return modelMap, err
		}
		modelMap["rotation"] = []map[string]interface{}{rotationMap}
	}
	if model.Algorithm != nil {
		modelMap["algorithm"] = *model.Algorithm
	}
	if model.KeyAlgorithm != nil {
		modelMap["key_algorithm"] = *model.KeyAlgorithm
	}
	if model.Issuer != nil {
		modelMap["issuer"] = *model.Issuer
	}
	if model.Validity != nil {
		validityMap, err := DataSourceIBMSmSecretsCertificateValidityToMap(model.Validity)
		if err != nil {
			return modelMap, err
		}
		modelMap["validity"] = []map[string]interface{}{validityMap}
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = *model.SerialNumber
	}
	if model.RevocationTime != nil {
		modelMap["revocation_time"] = *model.RevocationTime
	}
	if model.RevocationTimeRfc3339 != nil {
		modelMap["revocation_time_rfc3339"] = model.RevocationTimeRfc3339.String()
	}
	if model.SecretData != nil {
	}
	return modelMap, nil
}

func DataSourceIBMSmSecretsKvSecretResourceToMap(model *secretsmanagerv1.KvSecretResource) (map[string]interface{}, error) {
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
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.Versions != nil {
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.Payload != nil {
	}
	if model.SecretData != nil {
	}
	return modelMap, nil
}
