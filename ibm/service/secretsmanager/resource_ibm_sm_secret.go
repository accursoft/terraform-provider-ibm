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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
)

func ResourceIBMSmSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIBMSmSecretCreate,
		ReadContext:   ResourceIBMSmSecretRead,
		UpdateContext: ResourceIBMSmSecretUpdate,
		DeleteContext: ResourceIBMSmSecretDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"secret_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_sm_secret", "secret_type"),
				Description:  "The secret type. Allowable values are: arbitrary, iam_credentials, imported_cert, public_cert, username_password, kv.",
			},
			"secret_resource": &schema.Schema{
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The v4 UUID that uniquely identifies the secret.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A human-readable alias to assign to your secret.To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.",
						},
						"secret_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The v4 UUID that uniquely identifies the secret group to assign to this secret.If you omit this parameter, your secret is assigned to the `default` secret group.",
						},
						"labels": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Labels that you can use to filter for secrets in your instance.Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"state": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,  Suspended = 2, Deactivated = 3, and Destroyed = 5 values.",
						},
						"state_description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "A text representation of the secret state.",
						},
						"secret_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The secret type.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.",
						},
						"creation_date": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The date the secret was created. The date format follows RFC 3339.",
						},
						"created_by": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique identifier for the entity that created the secret.",
						},
						"last_update_date": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Updates when the actual secret is modified. The date format follows RFC 3339.",
						},
						"versions_total": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The number of versions that are associated with a secret.",
						},
						"versions": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "An array that contains metadata for each secret version. For more information on the metadata properties, see [Get secret version metadata](#get-secret-version-metadata).",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"expiration_date": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The date the secret material expires. The date format follows RFC 3339.You can set an expiration date on supported secret types at their creation. If you create a secret without specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the following secret types:- `arbitrary`- `username_password`.",
						},
						"payload": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The new secret data to assign to the secret.",
						},
						"secret_data": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "The data that is associated with the secret version.The data object contains the field `payload`.",
						},
						"username": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The username to assign to this secret.",
						},
						"password": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The password to assign to this secret.",
						},
						"next_rotation_date": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The date that the secret is scheduled for automatic rotation.The service automatically creates a new version of the secret on its next rotation date. This field exists only for secrets that can be auto-rotated and have an existing rotation policy.",
						},
						"ttl": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The time-to-live (TTL) or lease duration to assign to generated credentials.For `iam_credentials` secrets, the TTL defines for how long each generated API key remains valid. The value can be either an integer that specifies the number of seconds, or the string representation of a duration, such as `120m` or `24h`.Minimum duration is 1 minute. Maximum is 90 days.",
						},
						"access_groups": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The access groups that define the capabilities of the service ID and API key that are generated for an`iam_credentials` secret. If you prefer to use an existing service ID that is already assigned the access policies that you require, you can omit this parameter and use the `service_id` field instead.**Tip:** To list the access groups that are available in an account, you can use the [IAM Access Groups API](https://cloud.ibm.com/apidocs/iam-access-groups#list-access-groups). To find the ID of an access group in the console, go to **Manage > Access (IAM) > Access groups**. Select the access group to inspect, and click **Details** to view its ID.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"api_key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The API key that is generated for this secret.After the secret reaches the end of its lease (see the `ttl` field), the API key is deleted automatically. If you want to continue to use the same API key for future read operations, see the `reuse_api_key` field.",
						},
						"api_key_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The ID of the API key that is generated for this secret.",
						},
						"service_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "(IAM credentials) The service ID under which the API key is created. To have Secrets Manager generate a new service ID, omit this option and include 'access_groups'.",
						},
						"service_id_is_static": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether an `iam_credentials` secret was created with a static service ID.If `true`, the service ID for the secret was provided by the user at secret creation. If `false`, the service ID was generated by Secrets Manager.",
						},
						"reuse_api_key": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "(IAM credentials) Reuse the service ID and API key for future read operations.",
						},
						"certificate": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The contents of your certificate.",
						},
						"private_key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "(Optional) The private key to associate with the certificate.",
						},
						"intermediate": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "(Optional) The intermediate certificate to associate with the root certificate.",
						},
						"serial_number": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique serial number that was assigned to the certificate by the issuing certificate authority.",
						},
						"algorithm": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign the certificate.",
						},
						"key_algorithm": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The identifier for the cryptographic algorithm that was used to generate the public and private keys that are associated with the certificate.",
						},
						"issuer": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The distinguished name that identifies the entity that signed and issued the certificate.",
						},
						"validity": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"not_before": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The date and time that the certificate validity period begins.",
									},
									"not_after": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The date and time that the certificate validity period ends.",
									},
								},
							},
						},
						"common_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The fully qualified domain name or host domain name that is defined for the certificate.",
						},
						"intermediate_included": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the certificate was imported with an associated intermediate certificate.",
						},
						"private_key_included": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the certificate was imported with an associated private key.",
						},
						"alt_names": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The alternative names that are defined for the certificate.For public certificates, this value is provided as an array of strings. For private certificates, this value is provided as a comma-delimited list (string). In the API response, this value is returned as an array of strings for all the types of certificate secrets.",
						},
						"bundle_certs": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Determines whether your issued certificate is bundled with intermediate certificates.Set to `false` for the certificate file to contain only the issued certificate.",
						},
						"ca": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the certificate authority configuration.",
						},
						"dns": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the DNS provider configuration.",
						},
						"rotation": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_rotate": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Determines whether Secrets Manager rotates your certificate automatically.For public certificates, if `auto_rotate` is set to `true` the service reorders your certificate 31 days before it expires. For private certificates, the certificate is rotated according to the time interval specified in the `interval` and `unit` fields.To access the previous version of the certificate, you can use the[Get a version of a secret](#get-secret-version) method.",
									},
									"rotate_keys": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Determines whether Secrets Manager rotates the private key for your certificate automatically.If set to `true`, the service generates and stores a new private key for your rotated certificate.**Note:** Use this field only for public certificates. It is ignored for private certificates.",
									},
									"interval": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Used together with the `unit` field to specify the rotation interval. The minimum interval is one day, and the maximum interval is 3 years (1095 days). Required in case `auto_rotate` is set to `true`.**Note:** Use this field only for private certificates. It is ignored for public certificates.",
									},
									"unit": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The time unit of the rotation interval.**Note:** Use this field only for private certificates. It is ignored for public certificates.",
									},
								},
							},
						},
						"issuance_info": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Issuance information that is associated with your certificate.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ordered_on": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The date the certificate was ordered. The date format follows RFC 3339.",
									},
									"error_code": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "A code that identifies an issuance error.This field, along with `error_message`, is returned when Secrets Manager successfully processes your request, but a certificate is unable to be issued by the certificate authority.",
									},
									"error_message": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "A human-readable message that provides details about the issuance error.",
									},
									"bundle_certs": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Indicates whether the issued certificate is bundled with intermediate certificates.",
									},
									"state": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,  Suspended = 2, Deactivated = 3, and Destroyed = 5 values.",
									},
									"state_description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "A text representation of the secret state.",
									},
									"auto_rotated": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Indicates whether the issued certificate is configured with an automatic rotation policy.",
									},
									"ca": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The name that was assigned to the certificate authority configuration.",
									},
									"dns": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The name that was assigned to the DNS provider configuration.",
									},
								},
							},
						},
						"certificate_template": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the certificate template.",
						},
						"certificate_authority": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The intermediate certificate authority that signed this certificate.",
						},
						"ip_sans": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.",
						},
						"uri_sans": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.",
						},
						"other_sans": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"format": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "pem",
							Description: "The format of the returned data.",
						},
						"private_key_format": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "der",
							Description: "The format of the generated private key.",
						},
						"exclude_cn_from_sans": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Controls whether the common name is excluded from Subject Alternative Names (SANs). If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if the common name is not a hostname or an email address, but is instead a human-readable identifier.",
						},
						"revocation_time": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The timestamp of the certificate revocation.",
						},
						"revocation_time_rfc3339": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The date and time that the certificate was revoked. The date format follows RFC 3339.",
						},
					},
				},
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
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Elem:        &schema.Schema{Type: schema.TypeString},
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"not_before": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The date and time that the certificate validity period begins.",
						},
						"not_after": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_rotate": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether Secrets Manager rotates your certificate automatically.For public certificates, if `auto_rotate` is set to `true` the service reorders your certificate 31 days before it expires. For private certificates, the certificate is rotated according to the time interval specified in the `interval` and `unit` fields.To access the previous version of the certificate, you can use the[Get a version of a secret](#get-secret-version) method.",
						},
						"rotate_keys": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether Secrets Manager rotates the private key for your certificate automatically.If set to `true`, the service generates and stores a new private key for your rotated certificate.**Note:** Use this field only for public certificates. It is ignored for private certificates.",
						},
						"interval": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Used together with the `unit` field to specify the rotation interval. The minimum interval is one day, and the maximum interval is 3 years (1095 days). Required in case `auto_rotate` is set to `true`.**Note:** Use this field only for private certificates. It is ignored for public certificates.",
						},
						"unit": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "The date the certificate was ordered. The date format follows RFC 3339.",
						},
						"error_code": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "A code that identifies an issuance error.This field, along with `error_message`, is returned when Secrets Manager successfully processes your request, but a certificate is unable to be issued by the certificate authority.",
						},
						"error_message": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "A human-readable message that provides details about the issuance error.",
						},
						"bundle_certs": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the issued certificate is bundled with intermediate certificates.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,  Suspended = 2, Deactivated = 3, and Destroyed = 5 values.",
						},
						"state_description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "A text representation of the secret state.",
						},
						"auto_rotated": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the issued certificate is configured with an automatic rotation policy.",
						},
						"ca": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name that was assigned to the certificate authority configuration.",
						},
						"dns": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
				Elem:        &schema.Schema{Type: schema.TypeString},
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
			"secret_resource_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The v4 UUID that uniquely identifies the secret.",
			},
		},
	}
}

func ResourceIBMSmSecretValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "secret_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "arbitrary, iam_credentials, imported_cert, kv, private_cert, public_cert, username_password",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_sm_secret", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIBMSmSecretCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createSecretOptions := &secretsmanagerv1.CreateSecretOptions{}

	createSecretOptions.SetSecretType(d.Get("secret_type").(string))
	secretResourceModel, err := ResourceIBMSmSecretMapToSecretResource(d.Get("secret_resource.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createSecretOptions.SetSecretResource(secretResourceModel)

	secretResourceIntf, response, err := secretsManagerClient.CreateSecretWithContext(context, createSecretOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateSecretWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateSecretWithContext failed %s\n%s", err, response))
	}

	secretResource := secretResourceIntf.(*secretsmanagerv1.SecretResource)
	d.SetId(fmt.Sprintf("%s/%s", *createSecretOptions.SecretType, *secretResource.ID))

	return ResourceIBMSmSecretRead(context, d, meta)
}

func ResourceIBMSmSecretRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getSecretOptions := &secretsmanagerv1.GetSecretOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getSecretOptions.SetSecretType(parts[0])
	getSecretOptions.SetID(parts[1])

	secretResourceIntf, response, err := secretsManagerClient.GetSecretWithContext(context, getSecretOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSecretWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetSecretWithContext failed %s\n%s", err, response))
	}

	secretResource := secretResourceIntf.(*secretsmanagerv1.SecretResource)
	// TODO: handle argument of type SecretResource
	if err = d.Set("secret_type", secretResource.SecretType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_type: %s", err))
	}
	if err = d.Set("name", secretResource.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("description", secretResource.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if err = d.Set("secret_group_id", secretResource.SecretGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_group_id: %s", err))
	}
	if secretResource.Labels != nil {
		if err = d.Set("labels", secretResource.Labels); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting labels: %s", err))
		}
	}
	if err = d.Set("state", flex.IntValue(secretResource.State)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}
	if err = d.Set("state_description", secretResource.StateDescription); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state_description: %s", err))
	}
	if err = d.Set("crn", secretResource.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("creation_date", flex.DateTimeToString(secretResource.CreationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting creation_date: %s", err))
	}
	if err = d.Set("created_by", secretResource.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("last_update_date", flex.DateTimeToString(secretResource.LastUpdateDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_update_date: %s", err))
	}
	if err = d.Set("versions_total", flex.IntValue(secretResource.VersionsTotal)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting versions_total: %s", err))
	}
	versions := []map[string]interface{}{}
	if secretResource.Versions != nil {
		for _, versionsItem := range secretResource.Versions {
			// TODO: handle Versions of type TypeList -- list of non-primitive, not model items
		}
	}
	if err = d.Set("versions", versions); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting versions: %s", err))
	}
	if err = d.Set("expiration_date", flex.DateTimeToString(secretResource.ExpirationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting expiration_date: %s", err))
	}
	if err = d.Set("payload", secretResource.Payload); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting payload: %s", err))
	}
	if err = d.Set("secret_data", secretResource.SecretData); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_data: %s", err))
	}
	if err = d.Set("username", secretResource.Username); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting username: %s", err))
	}
	if err = d.Set("password", secretResource.Password); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting password: %s", err))
	}
	if err = d.Set("next_rotation_date", flex.DateTimeToString(secretResource.NextRotationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting next_rotation_date: %s", err))
	}
	if err = d.Set("ttl", secretResource.TTL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ttl: %s", err))
	}
	if secretResource.AccessGroups != nil {
		if err = d.Set("access_groups", secretResource.AccessGroups); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting access_groups: %s", err))
		}
	}
	if err = d.Set("api_key", secretResource.APIKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting api_key: %s", err))
	}
	if err = d.Set("api_key_id", secretResource.APIKeyID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting api_key_id: %s", err))
	}
	if err = d.Set("service_id", secretResource.ServiceID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_id: %s", err))
	}
	if err = d.Set("service_id_is_static", secretResource.ServiceIDIsStatic); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_id_is_static: %s", err))
	}
	if err = d.Set("reuse_api_key", secretResource.ReuseAPIKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting reuse_api_key: %s", err))
	}
	if err = d.Set("certificate", secretResource.Certificate); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting certificate: %s", err))
	}
	if err = d.Set("private_key", secretResource.PrivateKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting private_key: %s", err))
	}
	if err = d.Set("intermediate", secretResource.Intermediate); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting intermediate: %s", err))
	}
	if err = d.Set("serial_number", secretResource.SerialNumber); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting serial_number: %s", err))
	}
	if err = d.Set("algorithm", secretResource.Algorithm); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting algorithm: %s", err))
	}
	if err = d.Set("key_algorithm", secretResource.KeyAlgorithm); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key_algorithm: %s", err))
	}
	if err = d.Set("issuer", secretResource.Issuer); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting issuer: %s", err))
	}
	if secretResource.Validity != nil {
		validityMap, err := ResourceIBMSmSecretCertificateValidityToMap(secretResource.Validity)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("validity", []map[string]interface{}{validityMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting validity: %s", err))
		}
	}
	if err = d.Set("common_name", secretResource.CommonName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting common_name: %s", err))
	}
	if err = d.Set("intermediate_included", secretResource.IntermediateIncluded); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting intermediate_included: %s", err))
	}
	if err = d.Set("private_key_included", secretResource.PrivateKeyIncluded); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting private_key_included: %s", err))
	}
	if err = d.Set("alt_names", secretResource.AltNames); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting alt_names: %s", err))
	}
	if err = d.Set("bundle_certs", secretResource.BundleCerts); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting bundle_certs: %s", err))
	}
	if err = d.Set("ca", secretResource.Ca); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ca: %s", err))
	}
	if err = d.Set("dns", secretResource.DNS); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting dns: %s", err))
	}
	if secretResource.Rotation != nil {
		rotationMap, err := ResourceIBMSmSecretRotationToMap(secretResource.Rotation)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("rotation", []map[string]interface{}{rotationMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting rotation: %s", err))
		}
	}
	if secretResource.IssuanceInfo != nil {
		issuanceInfoMap, err := ResourceIBMSmSecretIssuanceInfoToMap(secretResource.IssuanceInfo)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("issuance_info", []map[string]interface{}{issuanceInfoMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting issuance_info: %s", err))
		}
	}
	if err = d.Set("certificate_template", secretResource.CertificateTemplate); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting certificate_template: %s", err))
	}
	if err = d.Set("certificate_authority", secretResource.CertificateAuthority); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting certificate_authority: %s", err))
	}
	if err = d.Set("ip_sans", secretResource.IPSans); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ip_sans: %s", err))
	}
	if err = d.Set("uri_sans", secretResource.URISans); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting uri_sans: %s", err))
	}
	if secretResource.OtherSans != nil {
		if err = d.Set("other_sans", secretResource.OtherSans); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting other_sans: %s", err))
		}
	}
	if err = d.Set("format", secretResource.Format); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting format: %s", err))
	}
	if err = d.Set("private_key_format", secretResource.PrivateKeyFormat); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting private_key_format: %s", err))
	}
	if err = d.Set("exclude_cn_from_sans", secretResource.ExcludeCnFromSans); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting exclude_cn_from_sans: %s", err))
	}
	if err = d.Set("revocation_time", flex.IntValue(secretResource.RevocationTime)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting revocation_time: %s", err))
	}
	if err = d.Set("revocation_time_rfc3339", flex.DateTimeToString(secretResource.RevocationTimeRfc3339)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting revocation_time_rfc3339: %s", err))
	}
	if err = d.Set("secret_resource_id", secretResource.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_resource_id: %s", err))
	}

	return nil
}

func ResourceIBMSmSecretUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateSecretMetadataOptions := &secretsmanagerv1.UpdateSecretMetadataOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateSecretMetadataOptions.SetSecretType(parts[0])
	updateSecretMetadataOptions.SetID(parts[1])

	hasChange := false

	if d.HasChange("secret_type") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "secret_type"))
	}
	if d.HasChange("secret_resource") {
		secretResource := d.Get("secret_resource.0").(map[string]interface{})
		updateSecretMetadataOptions.SetSecretResourcePatch(secretResource)
		hasChange = true
	}

	if hasChange {
		_, response, err := secretsManagerClient.UpdateSecretMetadataWithContext(context, updateSecretMetadataOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateSecretMetadataWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateSecretMetadataWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIBMSmSecretRead(context, d, meta)
}

func ResourceIBMSmSecretDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteSecretOptions := &secretsmanagerv1.DeleteSecretOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteSecretOptions.SetSecretType(parts[0])
	deleteSecretOptions.SetID(parts[1])

	response, err := secretsManagerClient.DeleteSecretWithContext(context, deleteSecretOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteSecretWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteSecretWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func ResourceIBMSmSecretMapToSecretResource(modelMap map[string]interface{}) (secretsmanagerv1.SecretResourceIntf, error) {
	model := &secretsmanagerv1.SecretResource{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["secret_group_id"] != nil && modelMap["secret_group_id"].(string) != "" {
		model.SecretGroupID = core.StringPtr(modelMap["secret_group_id"].(string))
	}
	if modelMap["labels"] != nil {
		labels := []string{}
		for _, labelsItem := range modelMap["labels"].([]interface{}) {
			labels = append(labels, labelsItem.(string))
		}
		model.Labels = labels
	}
	if modelMap["state"] != nil {
		model.State = core.Int64Ptr(int64(modelMap["state"].(int)))
	}
	if modelMap["state_description"] != nil && modelMap["state_description"].(string) != "" {
		model.StateDescription = core.StringPtr(modelMap["state_description"].(string))
	}
	if modelMap["secret_type"] != nil && modelMap["secret_type"].(string) != "" {
		model.SecretType = core.StringPtr(modelMap["secret_type"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["creation_date"] != nil {

	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["last_update_date"] != nil {

	}
	if modelMap["versions_total"] != nil {
		model.VersionsTotal = core.Int64Ptr(int64(modelMap["versions_total"].(int)))
	}
	if modelMap["versions"] != nil {
		versions := []map[string]interface{}{}
		for _, versionsItem := range modelMap["versions"].([]interface{}) {
			versions = append(versions, versionsItem.(map[string]interface{}))
		}
		model.Versions = versions
	}
	if modelMap["expiration_date"] != nil {

	}
	if modelMap["payload"] != nil && modelMap["payload"].(string) != "" {
		model.Payload = core.StringPtr(modelMap["payload"].(string))
	}
	if modelMap["secret_data"] != nil {

	}
	if modelMap["username"] != nil && modelMap["username"].(string) != "" {
		model.Username = core.StringPtr(modelMap["username"].(string))
	}
	if modelMap["password"] != nil && modelMap["password"].(string) != "" {
		model.Password = core.StringPtr(modelMap["password"].(string))
	}
	if modelMap["next_rotation_date"] != nil {

	}
	if modelMap["ttl"] != nil {

	}
	if modelMap["access_groups"] != nil {
		accessGroups := []string{}
		for _, accessGroupsItem := range modelMap["access_groups"].([]interface{}) {
			accessGroups = append(accessGroups, accessGroupsItem.(string))
		}
		model.AccessGroups = accessGroups
	}
	if modelMap["api_key"] != nil && modelMap["api_key"].(string) != "" {
		model.APIKey = core.StringPtr(modelMap["api_key"].(string))
	}
	if modelMap["api_key_id"] != nil && modelMap["api_key_id"].(string) != "" {
		model.APIKeyID = core.StringPtr(modelMap["api_key_id"].(string))
	}
	if modelMap["service_id"] != nil && modelMap["service_id"].(string) != "" {
		model.ServiceID = core.StringPtr(modelMap["service_id"].(string))
	}
	if modelMap["service_id_is_static"] != nil {
		model.ServiceIDIsStatic = core.BoolPtr(modelMap["service_id_is_static"].(bool))
	}
	if modelMap["reuse_api_key"] != nil {
		model.ReuseAPIKey = core.BoolPtr(modelMap["reuse_api_key"].(bool))
	}
	if modelMap["certificate"] != nil && modelMap["certificate"].(string) != "" {
		model.Certificate = core.StringPtr(modelMap["certificate"].(string))
	}
	if modelMap["private_key"] != nil && modelMap["private_key"].(string) != "" {
		model.PrivateKey = core.StringPtr(modelMap["private_key"].(string))
	}
	if modelMap["intermediate"] != nil && modelMap["intermediate"].(string) != "" {
		model.Intermediate = core.StringPtr(modelMap["intermediate"].(string))
	}
	if modelMap["serial_number"] != nil && modelMap["serial_number"].(string) != "" {
		model.SerialNumber = core.StringPtr(modelMap["serial_number"].(string))
	}
	if modelMap["algorithm"] != nil && modelMap["algorithm"].(string) != "" {
		model.Algorithm = core.StringPtr(modelMap["algorithm"].(string))
	}
	if modelMap["key_algorithm"] != nil && modelMap["key_algorithm"].(string) != "" {
		model.KeyAlgorithm = core.StringPtr(modelMap["key_algorithm"].(string))
	}
	if modelMap["issuer"] != nil && modelMap["issuer"].(string) != "" {
		model.Issuer = core.StringPtr(modelMap["issuer"].(string))
	}
	if modelMap["validity"] != nil && len(modelMap["validity"].([]interface{})) > 0 {
		ValidityModel, err := ResourceIBMSmSecretMapToCertificateValidity(modelMap["validity"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Validity = ValidityModel
	}
	if modelMap["common_name"] != nil && modelMap["common_name"].(string) != "" {
		model.CommonName = core.StringPtr(modelMap["common_name"].(string))
	}
	if modelMap["intermediate_included"] != nil {
		model.IntermediateIncluded = core.BoolPtr(modelMap["intermediate_included"].(bool))
	}
	if modelMap["private_key_included"] != nil {
		model.PrivateKeyIncluded = core.BoolPtr(modelMap["private_key_included"].(bool))
	}
	if modelMap["alt_names"] != nil {

	}
	if modelMap["bundle_certs"] != nil {
		model.BundleCerts = core.BoolPtr(modelMap["bundle_certs"].(bool))
	}
	if modelMap["ca"] != nil && modelMap["ca"].(string) != "" {
		model.Ca = core.StringPtr(modelMap["ca"].(string))
	}
	if modelMap["dns"] != nil && modelMap["dns"].(string) != "" {
		model.DNS = core.StringPtr(modelMap["dns"].(string))
	}
	if modelMap["rotation"] != nil && len(modelMap["rotation"].([]interface{})) > 0 {
		RotationModel, err := ResourceIBMSmSecretMapToRotation(modelMap["rotation"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Rotation = RotationModel
	}
	if modelMap["issuance_info"] != nil && len(modelMap["issuance_info"].([]interface{})) > 0 {
		IssuanceInfoModel, err := ResourceIBMSmSecretMapToIssuanceInfo(modelMap["issuance_info"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.IssuanceInfo = IssuanceInfoModel
	}
	if modelMap["certificate_template"] != nil && modelMap["certificate_template"].(string) != "" {
		model.CertificateTemplate = core.StringPtr(modelMap["certificate_template"].(string))
	}
	if modelMap["certificate_authority"] != nil && modelMap["certificate_authority"].(string) != "" {
		model.CertificateAuthority = core.StringPtr(modelMap["certificate_authority"].(string))
	}
	if modelMap["ip_sans"] != nil && modelMap["ip_sans"].(string) != "" {
		model.IPSans = core.StringPtr(modelMap["ip_sans"].(string))
	}
	if modelMap["uri_sans"] != nil && modelMap["uri_sans"].(string) != "" {
		model.URISans = core.StringPtr(modelMap["uri_sans"].(string))
	}
	if modelMap["other_sans"] != nil {
		otherSans := []string{}
		for _, otherSansItem := range modelMap["other_sans"].([]interface{}) {
			otherSans = append(otherSans, otherSansItem.(string))
		}
		model.OtherSans = otherSans
	}
	if modelMap["format"] != nil && modelMap["format"].(string) != "" {
		model.Format = core.StringPtr(modelMap["format"].(string))
	}
	if modelMap["private_key_format"] != nil && modelMap["private_key_format"].(string) != "" {
		model.PrivateKeyFormat = core.StringPtr(modelMap["private_key_format"].(string))
	}
	if modelMap["exclude_cn_from_sans"] != nil {
		model.ExcludeCnFromSans = core.BoolPtr(modelMap["exclude_cn_from_sans"].(bool))
	}
	if modelMap["revocation_time"] != nil {
		model.RevocationTime = core.Int64Ptr(int64(modelMap["revocation_time"].(int)))
	}
	if modelMap["revocation_time_rfc3339"] != nil {

	}
	return model, nil
}

func ResourceIBMSmSecretMapToCertificateValidity(modelMap map[string]interface{}) (*secretsmanagerv1.CertificateValidity, error) {
	model := &secretsmanagerv1.CertificateValidity{}
	if modelMap["not_before"] != nil {

	}
	if modelMap["not_after"] != nil {

	}
	return model, nil
}

func ResourceIBMSmSecretMapToRotation(modelMap map[string]interface{}) (*secretsmanagerv1.Rotation, error) {
	model := &secretsmanagerv1.Rotation{}
	if modelMap["auto_rotate"] != nil {
		model.AutoRotate = core.BoolPtr(modelMap["auto_rotate"].(bool))
	}
	if modelMap["rotate_keys"] != nil {
		model.RotateKeys = core.BoolPtr(modelMap["rotate_keys"].(bool))
	}
	if modelMap["interval"] != nil {
		model.Interval = core.Int64Ptr(int64(modelMap["interval"].(int)))
	}
	if modelMap["unit"] != nil && modelMap["unit"].(string) != "" {
		model.Unit = core.StringPtr(modelMap["unit"].(string))
	}
	return model, nil
}

func ResourceIBMSmSecretMapToIssuanceInfo(modelMap map[string]interface{}) (*secretsmanagerv1.IssuanceInfo, error) {
	model := &secretsmanagerv1.IssuanceInfo{}
	if modelMap["ordered_on"] != nil {

	}
	if modelMap["error_code"] != nil && modelMap["error_code"].(string) != "" {
		model.ErrorCode = core.StringPtr(modelMap["error_code"].(string))
	}
	if modelMap["error_message"] != nil && modelMap["error_message"].(string) != "" {
		model.ErrorMessage = core.StringPtr(modelMap["error_message"].(string))
	}
	if modelMap["bundle_certs"] != nil {
		model.BundleCerts = core.BoolPtr(modelMap["bundle_certs"].(bool))
	}
	if modelMap["state"] != nil {
		model.State = core.Int64Ptr(int64(modelMap["state"].(int)))
	}
	if modelMap["state_description"] != nil && modelMap["state_description"].(string) != "" {
		model.StateDescription = core.StringPtr(modelMap["state_description"].(string))
	}
	if modelMap["auto_rotated"] != nil {
		model.AutoRotated = core.BoolPtr(modelMap["auto_rotated"].(bool))
	}
	if modelMap["ca"] != nil && modelMap["ca"].(string) != "" {
		model.Ca = core.StringPtr(modelMap["ca"].(string))
	}
	if modelMap["dns"] != nil && modelMap["dns"].(string) != "" {
		model.DNS = core.StringPtr(modelMap["dns"].(string))
	}
	return model, nil
}

func ResourceIBMSmSecretMapToArbitrarySecretResource(modelMap map[string]interface{}) (*secretsmanagerv1.ArbitrarySecretResource, error) {
	model := &secretsmanagerv1.ArbitrarySecretResource{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["secret_group_id"] != nil && modelMap["secret_group_id"].(string) != "" {
		model.SecretGroupID = core.StringPtr(modelMap["secret_group_id"].(string))
	}
	if modelMap["labels"] != nil {
		labels := []string{}
		for _, labelsItem := range modelMap["labels"].([]interface{}) {
			labels = append(labels, labelsItem.(string))
		}
		model.Labels = labels
	}
	if modelMap["state"] != nil {
		model.State = core.Int64Ptr(int64(modelMap["state"].(int)))
	}
	if modelMap["state_description"] != nil && modelMap["state_description"].(string) != "" {
		model.StateDescription = core.StringPtr(modelMap["state_description"].(string))
	}
	if modelMap["secret_type"] != nil && modelMap["secret_type"].(string) != "" {
		model.SecretType = core.StringPtr(modelMap["secret_type"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["creation_date"] != nil {

	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["last_update_date"] != nil {

	}
	if modelMap["versions_total"] != nil {
		model.VersionsTotal = core.Int64Ptr(int64(modelMap["versions_total"].(int)))
	}
	if modelMap["versions"] != nil {
		versions := []map[string]interface{}{}
		for _, versionsItem := range modelMap["versions"].([]interface{}) {
			versions = append(versions, versionsItem.(map[string]interface{}))
		}
		model.Versions = versions
	}
	if modelMap["expiration_date"] != nil {

	}
	if modelMap["payload"] != nil && modelMap["payload"].(string) != "" {
		model.Payload = core.StringPtr(modelMap["payload"].(string))
	}
	if modelMap["secret_data"] != nil {

	}
	return model, nil
}

func ResourceIBMSmSecretMapToUsernamePasswordSecretResource(modelMap map[string]interface{}) (*secretsmanagerv1.UsernamePasswordSecretResource, error) {
	model := &secretsmanagerv1.UsernamePasswordSecretResource{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["secret_group_id"] != nil && modelMap["secret_group_id"].(string) != "" {
		model.SecretGroupID = core.StringPtr(modelMap["secret_group_id"].(string))
	}
	if modelMap["labels"] != nil {
		labels := []string{}
		for _, labelsItem := range modelMap["labels"].([]interface{}) {
			labels = append(labels, labelsItem.(string))
		}
		model.Labels = labels
	}
	if modelMap["state"] != nil {
		model.State = core.Int64Ptr(int64(modelMap["state"].(int)))
	}
	if modelMap["state_description"] != nil && modelMap["state_description"].(string) != "" {
		model.StateDescription = core.StringPtr(modelMap["state_description"].(string))
	}
	if modelMap["secret_type"] != nil && modelMap["secret_type"].(string) != "" {
		model.SecretType = core.StringPtr(modelMap["secret_type"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["creation_date"] != nil {

	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["last_update_date"] != nil {

	}
	if modelMap["versions_total"] != nil {
		model.VersionsTotal = core.Int64Ptr(int64(modelMap["versions_total"].(int)))
	}
	if modelMap["versions"] != nil {
		versions := []map[string]interface{}{}
		for _, versionsItem := range modelMap["versions"].([]interface{}) {
			versions = append(versions, versionsItem.(map[string]interface{}))
		}
		model.Versions = versions
	}
	if modelMap["username"] != nil && modelMap["username"].(string) != "" {
		model.Username = core.StringPtr(modelMap["username"].(string))
	}
	if modelMap["password"] != nil && modelMap["password"].(string) != "" {
		model.Password = core.StringPtr(modelMap["password"].(string))
	}
	if modelMap["secret_data"] != nil {

	}
	if modelMap["expiration_date"] != nil {

	}
	if modelMap["next_rotation_date"] != nil {

	}
	return model, nil
}

func ResourceIBMSmSecretMapToIamCredentialsSecretResource(modelMap map[string]interface{}) (*secretsmanagerv1.IamCredentialsSecretResource, error) {
	model := &secretsmanagerv1.IamCredentialsSecretResource{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["secret_group_id"] != nil && modelMap["secret_group_id"].(string) != "" {
		model.SecretGroupID = core.StringPtr(modelMap["secret_group_id"].(string))
	}
	if modelMap["labels"] != nil {
		labels := []string{}
		for _, labelsItem := range modelMap["labels"].([]interface{}) {
			labels = append(labels, labelsItem.(string))
		}
		model.Labels = labels
	}
	if modelMap["state"] != nil {
		model.State = core.Int64Ptr(int64(modelMap["state"].(int)))
	}
	if modelMap["state_description"] != nil && modelMap["state_description"].(string) != "" {
		model.StateDescription = core.StringPtr(modelMap["state_description"].(string))
	}
	if modelMap["secret_type"] != nil && modelMap["secret_type"].(string) != "" {
		model.SecretType = core.StringPtr(modelMap["secret_type"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["creation_date"] != nil {

	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["last_update_date"] != nil {

	}
	if modelMap["versions_total"] != nil {
		model.VersionsTotal = core.Int64Ptr(int64(modelMap["versions_total"].(int)))
	}
	if modelMap["versions"] != nil {
		versions := []map[string]interface{}{}
		for _, versionsItem := range modelMap["versions"].([]interface{}) {
			versions = append(versions, versionsItem.(map[string]interface{}))
		}
		model.Versions = versions
	}
	if modelMap["ttl"] != nil {

	}
	if modelMap["access_groups"] != nil {
		accessGroups := []string{}
		for _, accessGroupsItem := range modelMap["access_groups"].([]interface{}) {
			accessGroups = append(accessGroups, accessGroupsItem.(string))
		}
		model.AccessGroups = accessGroups
	}
	if modelMap["api_key"] != nil && modelMap["api_key"].(string) != "" {
		model.APIKey = core.StringPtr(modelMap["api_key"].(string))
	}
	if modelMap["api_key_id"] != nil && modelMap["api_key_id"].(string) != "" {
		model.APIKeyID = core.StringPtr(modelMap["api_key_id"].(string))
	}
	if modelMap["service_id"] != nil && modelMap["service_id"].(string) != "" {
		model.ServiceID = core.StringPtr(modelMap["service_id"].(string))
	}
	if modelMap["service_id_is_static"] != nil {
		model.ServiceIDIsStatic = core.BoolPtr(modelMap["service_id_is_static"].(bool))
	}
	if modelMap["reuse_api_key"] != nil {
		model.ReuseAPIKey = core.BoolPtr(modelMap["reuse_api_key"].(bool))
	}
	return model, nil
}

func ResourceIBMSmSecretMapToCertificateSecretResource(modelMap map[string]interface{}) (*secretsmanagerv1.CertificateSecretResource, error) {
	model := &secretsmanagerv1.CertificateSecretResource{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["secret_group_id"] != nil && modelMap["secret_group_id"].(string) != "" {
		model.SecretGroupID = core.StringPtr(modelMap["secret_group_id"].(string))
	}
	if modelMap["labels"] != nil {
		labels := []string{}
		for _, labelsItem := range modelMap["labels"].([]interface{}) {
			labels = append(labels, labelsItem.(string))
		}
		model.Labels = labels
	}
	if modelMap["state"] != nil {
		model.State = core.Int64Ptr(int64(modelMap["state"].(int)))
	}
	if modelMap["state_description"] != nil && modelMap["state_description"].(string) != "" {
		model.StateDescription = core.StringPtr(modelMap["state_description"].(string))
	}
	if modelMap["secret_type"] != nil && modelMap["secret_type"].(string) != "" {
		model.SecretType = core.StringPtr(modelMap["secret_type"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["creation_date"] != nil {

	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["last_update_date"] != nil {

	}
	if modelMap["versions_total"] != nil {
		model.VersionsTotal = core.Int64Ptr(int64(modelMap["versions_total"].(int)))
	}
	if modelMap["versions"] != nil {
		versions := []map[string]interface{}{}
		for _, versionsItem := range modelMap["versions"].([]interface{}) {
			versions = append(versions, versionsItem.(map[string]interface{}))
		}
		model.Versions = versions
	}
	if modelMap["certificate"] != nil && modelMap["certificate"].(string) != "" {
		model.Certificate = core.StringPtr(modelMap["certificate"].(string))
	}
	if modelMap["private_key"] != nil && modelMap["private_key"].(string) != "" {
		model.PrivateKey = core.StringPtr(modelMap["private_key"].(string))
	}
	if modelMap["intermediate"] != nil && modelMap["intermediate"].(string) != "" {
		model.Intermediate = core.StringPtr(modelMap["intermediate"].(string))
	}
	if modelMap["secret_data"] != nil {

	}
	if modelMap["serial_number"] != nil && modelMap["serial_number"].(string) != "" {
		model.SerialNumber = core.StringPtr(modelMap["serial_number"].(string))
	}
	if modelMap["algorithm"] != nil && modelMap["algorithm"].(string) != "" {
		model.Algorithm = core.StringPtr(modelMap["algorithm"].(string))
	}
	if modelMap["key_algorithm"] != nil && modelMap["key_algorithm"].(string) != "" {
		model.KeyAlgorithm = core.StringPtr(modelMap["key_algorithm"].(string))
	}
	if modelMap["issuer"] != nil && modelMap["issuer"].(string) != "" {
		model.Issuer = core.StringPtr(modelMap["issuer"].(string))
	}
	if modelMap["validity"] != nil && len(modelMap["validity"].([]interface{})) > 0 {
		ValidityModel, err := ResourceIBMSmSecretMapToCertificateValidity(modelMap["validity"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Validity = ValidityModel
	}
	if modelMap["common_name"] != nil && modelMap["common_name"].(string) != "" {
		model.CommonName = core.StringPtr(modelMap["common_name"].(string))
	}
	if modelMap["intermediate_included"] != nil {
		model.IntermediateIncluded = core.BoolPtr(modelMap["intermediate_included"].(bool))
	}
	if modelMap["private_key_included"] != nil {
		model.PrivateKeyIncluded = core.BoolPtr(modelMap["private_key_included"].(bool))
	}
	if modelMap["alt_names"] != nil {

	}
	if modelMap["expiration_date"] != nil {

	}
	return model, nil
}

func ResourceIBMSmSecretMapToPublicCertificateSecretResource(modelMap map[string]interface{}) (*secretsmanagerv1.PublicCertificateSecretResource, error) {
	model := &secretsmanagerv1.PublicCertificateSecretResource{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["secret_group_id"] != nil && modelMap["secret_group_id"].(string) != "" {
		model.SecretGroupID = core.StringPtr(modelMap["secret_group_id"].(string))
	}
	if modelMap["labels"] != nil {
		labels := []string{}
		for _, labelsItem := range modelMap["labels"].([]interface{}) {
			labels = append(labels, labelsItem.(string))
		}
		model.Labels = labels
	}
	if modelMap["state"] != nil {
		model.State = core.Int64Ptr(int64(modelMap["state"].(int)))
	}
	if modelMap["state_description"] != nil && modelMap["state_description"].(string) != "" {
		model.StateDescription = core.StringPtr(modelMap["state_description"].(string))
	}
	if modelMap["secret_type"] != nil && modelMap["secret_type"].(string) != "" {
		model.SecretType = core.StringPtr(modelMap["secret_type"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["creation_date"] != nil {

	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["last_update_date"] != nil {

	}
	if modelMap["versions_total"] != nil {
		model.VersionsTotal = core.Int64Ptr(int64(modelMap["versions_total"].(int)))
	}
	if modelMap["versions"] != nil {
		versions := []map[string]interface{}{}
		for _, versionsItem := range modelMap["versions"].([]interface{}) {
			versions = append(versions, versionsItem.(map[string]interface{}))
		}
		model.Versions = versions
	}
	if modelMap["issuer"] != nil && modelMap["issuer"].(string) != "" {
		model.Issuer = core.StringPtr(modelMap["issuer"].(string))
	}
	if modelMap["bundle_certs"] != nil {
		model.BundleCerts = core.BoolPtr(modelMap["bundle_certs"].(bool))
	}
	if modelMap["ca"] != nil && modelMap["ca"].(string) != "" {
		model.Ca = core.StringPtr(modelMap["ca"].(string))
	}
	if modelMap["dns"] != nil && modelMap["dns"].(string) != "" {
		model.DNS = core.StringPtr(modelMap["dns"].(string))
	}
	if modelMap["algorithm"] != nil && modelMap["algorithm"].(string) != "" {
		model.Algorithm = core.StringPtr(modelMap["algorithm"].(string))
	}
	if modelMap["key_algorithm"] != nil && modelMap["key_algorithm"].(string) != "" {
		model.KeyAlgorithm = core.StringPtr(modelMap["key_algorithm"].(string))
	}
	if modelMap["alt_names"] != nil {

	}
	if modelMap["common_name"] != nil && modelMap["common_name"].(string) != "" {
		model.CommonName = core.StringPtr(modelMap["common_name"].(string))
	}
	if modelMap["private_key_included"] != nil {
		model.PrivateKeyIncluded = core.BoolPtr(modelMap["private_key_included"].(bool))
	}
	if modelMap["intermediate_included"] != nil {
		model.IntermediateIncluded = core.BoolPtr(modelMap["intermediate_included"].(bool))
	}
	if modelMap["rotation"] != nil && len(modelMap["rotation"].([]interface{})) > 0 {
		RotationModel, err := ResourceIBMSmSecretMapToRotation(modelMap["rotation"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Rotation = RotationModel
	}
	if modelMap["issuance_info"] != nil && len(modelMap["issuance_info"].([]interface{})) > 0 {
		IssuanceInfoModel, err := ResourceIBMSmSecretMapToIssuanceInfo(modelMap["issuance_info"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.IssuanceInfo = IssuanceInfoModel
	}
	if modelMap["validity"] != nil && len(modelMap["validity"].([]interface{})) > 0 {
		ValidityModel, err := ResourceIBMSmSecretMapToCertificateValidity(modelMap["validity"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Validity = ValidityModel
	}
	if modelMap["serial_number"] != nil && modelMap["serial_number"].(string) != "" {
		model.SerialNumber = core.StringPtr(modelMap["serial_number"].(string))
	}
	if modelMap["secret_data"] != nil {

	}
	return model, nil
}

func ResourceIBMSmSecretMapToPrivateCertificateSecretResource(modelMap map[string]interface{}) (*secretsmanagerv1.PrivateCertificateSecretResource, error) {
	model := &secretsmanagerv1.PrivateCertificateSecretResource{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["secret_group_id"] != nil && modelMap["secret_group_id"].(string) != "" {
		model.SecretGroupID = core.StringPtr(modelMap["secret_group_id"].(string))
	}
	if modelMap["labels"] != nil {
		labels := []string{}
		for _, labelsItem := range modelMap["labels"].([]interface{}) {
			labels = append(labels, labelsItem.(string))
		}
		model.Labels = labels
	}
	if modelMap["state"] != nil {
		model.State = core.Int64Ptr(int64(modelMap["state"].(int)))
	}
	if modelMap["state_description"] != nil && modelMap["state_description"].(string) != "" {
		model.StateDescription = core.StringPtr(modelMap["state_description"].(string))
	}
	if modelMap["secret_type"] != nil && modelMap["secret_type"].(string) != "" {
		model.SecretType = core.StringPtr(modelMap["secret_type"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["creation_date"] != nil {

	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["last_update_date"] != nil {

	}
	if modelMap["versions_total"] != nil {
		model.VersionsTotal = core.Int64Ptr(int64(modelMap["versions_total"].(int)))
	}
	if modelMap["versions"] != nil {
		versions := []map[string]interface{}{}
		for _, versionsItem := range modelMap["versions"].([]interface{}) {
			versions = append(versions, versionsItem.(map[string]interface{}))
		}
		model.Versions = versions
	}
	model.CertificateTemplate = core.StringPtr(modelMap["certificate_template"].(string))
	if modelMap["certificate_authority"] != nil && modelMap["certificate_authority"].(string) != "" {
		model.CertificateAuthority = core.StringPtr(modelMap["certificate_authority"].(string))
	}
	model.CommonName = core.StringPtr(modelMap["common_name"].(string))
	if modelMap["alt_names"] != nil {

	}
	if modelMap["ip_sans"] != nil && modelMap["ip_sans"].(string) != "" {
		model.IPSans = core.StringPtr(modelMap["ip_sans"].(string))
	}
	if modelMap["uri_sans"] != nil && modelMap["uri_sans"].(string) != "" {
		model.URISans = core.StringPtr(modelMap["uri_sans"].(string))
	}
	if modelMap["other_sans"] != nil {
		otherSans := []string{}
		for _, otherSansItem := range modelMap["other_sans"].([]interface{}) {
			otherSans = append(otherSans, otherSansItem.(string))
		}
		model.OtherSans = otherSans
	}
	if modelMap["ttl"] != nil {

	}
	if modelMap["format"] != nil && modelMap["format"].(string) != "" {
		model.Format = core.StringPtr(modelMap["format"].(string))
	}
	if modelMap["private_key_format"] != nil && modelMap["private_key_format"].(string) != "" {
		model.PrivateKeyFormat = core.StringPtr(modelMap["private_key_format"].(string))
	}
	if modelMap["exclude_cn_from_sans"] != nil {
		model.ExcludeCnFromSans = core.BoolPtr(modelMap["exclude_cn_from_sans"].(bool))
	}
	if modelMap["rotation"] != nil && len(modelMap["rotation"].([]interface{})) > 0 {
		RotationModel, err := ResourceIBMSmSecretMapToRotation(modelMap["rotation"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Rotation = RotationModel
	}
	if modelMap["algorithm"] != nil && modelMap["algorithm"].(string) != "" {
		model.Algorithm = core.StringPtr(modelMap["algorithm"].(string))
	}
	if modelMap["key_algorithm"] != nil && modelMap["key_algorithm"].(string) != "" {
		model.KeyAlgorithm = core.StringPtr(modelMap["key_algorithm"].(string))
	}
	if modelMap["issuer"] != nil && modelMap["issuer"].(string) != "" {
		model.Issuer = core.StringPtr(modelMap["issuer"].(string))
	}
	if modelMap["validity"] != nil && len(modelMap["validity"].([]interface{})) > 0 {
		ValidityModel, err := ResourceIBMSmSecretMapToCertificateValidity(modelMap["validity"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Validity = ValidityModel
	}
	if modelMap["serial_number"] != nil && modelMap["serial_number"].(string) != "" {
		model.SerialNumber = core.StringPtr(modelMap["serial_number"].(string))
	}
	if modelMap["revocation_time"] != nil {
		model.RevocationTime = core.Int64Ptr(int64(modelMap["revocation_time"].(int)))
	}
	if modelMap["revocation_time_rfc3339"] != nil {

	}
	if modelMap["secret_data"] != nil {

	}
	return model, nil
}

func ResourceIBMSmSecretMapToKvSecretResource(modelMap map[string]interface{}) (*secretsmanagerv1.KvSecretResource, error) {
	model := &secretsmanagerv1.KvSecretResource{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["secret_group_id"] != nil && modelMap["secret_group_id"].(string) != "" {
		model.SecretGroupID = core.StringPtr(modelMap["secret_group_id"].(string))
	}
	if modelMap["labels"] != nil {
		labels := []string{}
		for _, labelsItem := range modelMap["labels"].([]interface{}) {
			labels = append(labels, labelsItem.(string))
		}
		model.Labels = labels
	}
	if modelMap["state"] != nil {
		model.State = core.Int64Ptr(int64(modelMap["state"].(int)))
	}
	if modelMap["state_description"] != nil && modelMap["state_description"].(string) != "" {
		model.StateDescription = core.StringPtr(modelMap["state_description"].(string))
	}
	if modelMap["secret_type"] != nil && modelMap["secret_type"].(string) != "" {
		model.SecretType = core.StringPtr(modelMap["secret_type"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["creation_date"] != nil {

	}
	if modelMap["created_by"] != nil && modelMap["created_by"].(string) != "" {
		model.CreatedBy = core.StringPtr(modelMap["created_by"].(string))
	}
	if modelMap["last_update_date"] != nil {

	}
	if modelMap["versions_total"] != nil {
		model.VersionsTotal = core.Int64Ptr(int64(modelMap["versions_total"].(int)))
	}
	if modelMap["versions"] != nil {
		versions := []map[string]interface{}{}
		for _, versionsItem := range modelMap["versions"].([]interface{}) {
			versions = append(versions, versionsItem.(map[string]interface{}))
		}
		model.Versions = versions
	}
	if modelMap["expiration_date"] != nil {

	}
	if modelMap["payload"] != nil {

	}
	if modelMap["secret_data"] != nil {

	}
	return model, nil
}

func ResourceIBMSmSecretSecretResourceToMap(model secretsmanagerv1.SecretResourceIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv1.ArbitrarySecretResource); ok {
		return ResourceIBMSmSecretArbitrarySecretResourceToMap(model.(*secretsmanagerv1.ArbitrarySecretResource))
	} else if _, ok := model.(*secretsmanagerv1.UsernamePasswordSecretResource); ok {
		return ResourceIBMSmSecretUsernamePasswordSecretResourceToMap(model.(*secretsmanagerv1.UsernamePasswordSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.IamCredentialsSecretResource); ok {
		return ResourceIBMSmSecretIamCredentialsSecretResourceToMap(model.(*secretsmanagerv1.IamCredentialsSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.CertificateSecretResource); ok {
		return ResourceIBMSmSecretCertificateSecretResourceToMap(model.(*secretsmanagerv1.CertificateSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.PublicCertificateSecretResource); ok {
		return ResourceIBMSmSecretPublicCertificateSecretResourceToMap(model.(*secretsmanagerv1.PublicCertificateSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.PrivateCertificateSecretResource); ok {
		return ResourceIBMSmSecretPrivateCertificateSecretResourceToMap(model.(*secretsmanagerv1.PrivateCertificateSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.KvSecretResource); ok {
		return ResourceIBMSmSecretKvSecretResourceToMap(model.(*secretsmanagerv1.KvSecretResource))
	} else if _, ok := model.(*secretsmanagerv1.SecretResource); ok {
		modelMap := make(map[string]interface{})
		model := model.(*secretsmanagerv1.SecretResource)
		if model.ID != nil {
			modelMap["id"] = model.ID
		}
		if model.Name != nil {
			modelMap["name"] = model.Name
		}
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.SecretGroupID != nil {
			modelMap["secret_group_id"] = model.SecretGroupID
		}
		if model.Labels != nil {
			modelMap["labels"] = model.Labels
		}
		if model.State != nil {
			modelMap["state"] = flex.IntValue(model.State)
		}
		if model.StateDescription != nil {
			modelMap["state_description"] = model.StateDescription
		}
		if model.SecretType != nil {
			modelMap["secret_type"] = model.SecretType
		}
		if model.CRN != nil {
			modelMap["crn"] = model.CRN
		}
		if model.CreationDate != nil {
			modelMap["creation_date"] = model.CreationDate.String()
		}
		if model.CreatedBy != nil {
			modelMap["created_by"] = model.CreatedBy
		}
		if model.LastUpdateDate != nil {
			modelMap["last_update_date"] = model.LastUpdateDate.String()
		}
		if model.VersionsTotal != nil {
			modelMap["versions_total"] = flex.IntValue(model.VersionsTotal)
		}
		if model.Versions != nil {
			modelMap["versions"] = model.Versions
		}
		if model.ExpirationDate != nil {
			modelMap["expiration_date"] = model.ExpirationDate.String()
		}
		if model.Payload != nil {
			modelMap["payload"] = model.Payload
		}
		if model.SecretData != nil {
			modelMap["secret_data"] = model.SecretData
		}
		if model.Username != nil {
			modelMap["username"] = model.Username
		}
		if model.Password != nil {
			modelMap["password"] = model.Password
		}
		if model.NextRotationDate != nil {
			modelMap["next_rotation_date"] = model.NextRotationDate.String()
		}
		if model.TTL != nil {
			modelMap["ttl"] = model.TTL
		}
		if model.AccessGroups != nil {
			modelMap["access_groups"] = model.AccessGroups
		}
		if model.APIKey != nil {
			modelMap["api_key"] = model.APIKey
		}
		if model.APIKeyID != nil {
			modelMap["api_key_id"] = model.APIKeyID
		}
		if model.ServiceID != nil {
			modelMap["service_id"] = model.ServiceID
		}
		if model.ServiceIDIsStatic != nil {
			modelMap["service_id_is_static"] = model.ServiceIDIsStatic
		}
		if model.ReuseAPIKey != nil {
			modelMap["reuse_api_key"] = model.ReuseAPIKey
		}
		if model.Certificate != nil {
			modelMap["certificate"] = model.Certificate
		}
		if model.PrivateKey != nil {
			modelMap["private_key"] = model.PrivateKey
		}
		if model.Intermediate != nil {
			modelMap["intermediate"] = model.Intermediate
		}
		if model.SerialNumber != nil {
			modelMap["serial_number"] = model.SerialNumber
		}
		if model.Algorithm != nil {
			modelMap["algorithm"] = model.Algorithm
		}
		if model.KeyAlgorithm != nil {
			modelMap["key_algorithm"] = model.KeyAlgorithm
		}
		if model.Issuer != nil {
			modelMap["issuer"] = model.Issuer
		}
		if model.Validity != nil {
			validityMap, err := ResourceIBMSmSecretCertificateValidityToMap(model.Validity)
			if err != nil {
				return modelMap, err
			}
			modelMap["validity"] = []map[string]interface{}{validityMap}
		}
		if model.CommonName != nil {
			modelMap["common_name"] = model.CommonName
		}
		if model.IntermediateIncluded != nil {
			modelMap["intermediate_included"] = model.IntermediateIncluded
		}
		if model.PrivateKeyIncluded != nil {
			modelMap["private_key_included"] = model.PrivateKeyIncluded
		}
		if model.AltNames != nil {
			modelMap["alt_names"] = model.AltNames
		}
		if model.BundleCerts != nil {
			modelMap["bundle_certs"] = model.BundleCerts
		}
		if model.Ca != nil {
			modelMap["ca"] = model.Ca
		}
		if model.DNS != nil {
			modelMap["dns"] = model.DNS
		}
		if model.Rotation != nil {
			rotationMap, err := ResourceIBMSmSecretRotationToMap(model.Rotation)
			if err != nil {
				return modelMap, err
			}
			modelMap["rotation"] = []map[string]interface{}{rotationMap}
		}
		if model.IssuanceInfo != nil {
			issuanceInfoMap, err := ResourceIBMSmSecretIssuanceInfoToMap(model.IssuanceInfo)
			if err != nil {
				return modelMap, err
			}
			modelMap["issuance_info"] = []map[string]interface{}{issuanceInfoMap}
		}
		if model.CertificateTemplate != nil {
			modelMap["certificate_template"] = model.CertificateTemplate
		}
		if model.CertificateAuthority != nil {
			modelMap["certificate_authority"] = model.CertificateAuthority
		}
		if model.IPSans != nil {
			modelMap["ip_sans"] = model.IPSans
		}
		if model.URISans != nil {
			modelMap["uri_sans"] = model.URISans
		}
		if model.OtherSans != nil {
			modelMap["other_sans"] = model.OtherSans
		}
		if model.Format != nil {
			modelMap["format"] = model.Format
		}
		if model.PrivateKeyFormat != nil {
			modelMap["private_key_format"] = model.PrivateKeyFormat
		}
		if model.ExcludeCnFromSans != nil {
			modelMap["exclude_cn_from_sans"] = model.ExcludeCnFromSans
		}
		if model.RevocationTime != nil {
			modelMap["revocation_time"] = flex.IntValue(model.RevocationTime)
		}
		if model.RevocationTimeRfc3339 != nil {
			modelMap["revocation_time_rfc3339"] = model.RevocationTimeRfc3339.String()
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv1.SecretResourceIntf subtype encountered")
	}
}

func ResourceIBMSmSecretCertificateValidityToMap(model *secretsmanagerv1.CertificateValidity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NotBefore != nil {
		modelMap["not_before"] = model.NotBefore.String()
	}
	if model.NotAfter != nil {
		modelMap["not_after"] = model.NotAfter.String()
	}
	return modelMap, nil
}

func ResourceIBMSmSecretRotationToMap(model *secretsmanagerv1.Rotation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotate != nil {
		modelMap["auto_rotate"] = model.AutoRotate
	}
	if model.RotateKeys != nil {
		modelMap["rotate_keys"] = model.RotateKeys
	}
	if model.Interval != nil {
		modelMap["interval"] = flex.IntValue(model.Interval)
	}
	if model.Unit != nil {
		modelMap["unit"] = model.Unit
	}
	return modelMap, nil
}

func ResourceIBMSmSecretIssuanceInfoToMap(model *secretsmanagerv1.IssuanceInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.OrderedOn != nil {
		modelMap["ordered_on"] = model.OrderedOn.String()
	}
	if model.ErrorCode != nil {
		modelMap["error_code"] = model.ErrorCode
	}
	if model.ErrorMessage != nil {
		modelMap["error_message"] = model.ErrorMessage
	}
	if model.BundleCerts != nil {
		modelMap["bundle_certs"] = model.BundleCerts
	}
	if model.State != nil {
		modelMap["state"] = flex.IntValue(model.State)
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = model.StateDescription
	}
	if model.AutoRotated != nil {
		modelMap["auto_rotated"] = model.AutoRotated
	}
	if model.Ca != nil {
		modelMap["ca"] = model.Ca
	}
	if model.DNS != nil {
		modelMap["dns"] = model.DNS
	}
	return modelMap, nil
}

func ResourceIBMSmSecretArbitrarySecretResourceToMap(model *secretsmanagerv1.ArbitrarySecretResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = flex.IntValue(model.State)
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = flex.IntValue(model.VersionsTotal)
	}
	if model.Versions != nil {
		modelMap["versions"] = model.Versions
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.Payload != nil {
		modelMap["payload"] = model.Payload
	}
	if model.SecretData != nil {
		modelMap["secret_data"] = model.SecretData
	}
	return modelMap, nil
}

func ResourceIBMSmSecretUsernamePasswordSecretResourceToMap(model *secretsmanagerv1.UsernamePasswordSecretResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = flex.IntValue(model.State)
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = flex.IntValue(model.VersionsTotal)
	}
	if model.Versions != nil {
		modelMap["versions"] = model.Versions
	}
	if model.Username != nil {
		modelMap["username"] = model.Username
	}
	if model.Password != nil {
		modelMap["password"] = model.Password
	}
	if model.SecretData != nil {
		modelMap["secret_data"] = model.SecretData
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.NextRotationDate != nil {
		modelMap["next_rotation_date"] = model.NextRotationDate.String()
	}
	return modelMap, nil
}

func ResourceIBMSmSecretIamCredentialsSecretResourceToMap(model *secretsmanagerv1.IamCredentialsSecretResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = flex.IntValue(model.State)
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = flex.IntValue(model.VersionsTotal)
	}
	if model.Versions != nil {
		modelMap["versions"] = model.Versions
	}
	if model.TTL != nil {
		modelMap["ttl"] = model.TTL
	}
	if model.AccessGroups != nil {
		modelMap["access_groups"] = model.AccessGroups
	}
	if model.APIKey != nil {
		modelMap["api_key"] = model.APIKey
	}
	if model.APIKeyID != nil {
		modelMap["api_key_id"] = model.APIKeyID
	}
	if model.ServiceID != nil {
		modelMap["service_id"] = model.ServiceID
	}
	if model.ServiceIDIsStatic != nil {
		modelMap["service_id_is_static"] = model.ServiceIDIsStatic
	}
	if model.ReuseAPIKey != nil {
		modelMap["reuse_api_key"] = model.ReuseAPIKey
	}
	return modelMap, nil
}

func ResourceIBMSmSecretCertificateSecretResourceToMap(model *secretsmanagerv1.CertificateSecretResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = flex.IntValue(model.State)
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = flex.IntValue(model.VersionsTotal)
	}
	if model.Versions != nil {
		modelMap["versions"] = model.Versions
	}
	if model.Certificate != nil {
		modelMap["certificate"] = model.Certificate
	}
	if model.PrivateKey != nil {
		modelMap["private_key"] = model.PrivateKey
	}
	if model.Intermediate != nil {
		modelMap["intermediate"] = model.Intermediate
	}
	if model.SecretData != nil {
		modelMap["secret_data"] = model.SecretData
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = model.SerialNumber
	}
	if model.Algorithm != nil {
		modelMap["algorithm"] = model.Algorithm
	}
	if model.KeyAlgorithm != nil {
		modelMap["key_algorithm"] = model.KeyAlgorithm
	}
	if model.Issuer != nil {
		modelMap["issuer"] = model.Issuer
	}
	if model.Validity != nil {
		validityMap, err := ResourceIBMSmSecretCertificateValidityToMap(model.Validity)
		if err != nil {
			return modelMap, err
		}
		modelMap["validity"] = []map[string]interface{}{validityMap}
	}
	if model.CommonName != nil {
		modelMap["common_name"] = model.CommonName
	}
	if model.IntermediateIncluded != nil {
		modelMap["intermediate_included"] = model.IntermediateIncluded
	}
	if model.PrivateKeyIncluded != nil {
		modelMap["private_key_included"] = model.PrivateKeyIncluded
	}
	if model.AltNames != nil {
		modelMap["alt_names"] = model.AltNames
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	return modelMap, nil
}

func ResourceIBMSmSecretPublicCertificateSecretResourceToMap(model *secretsmanagerv1.PublicCertificateSecretResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = flex.IntValue(model.State)
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = flex.IntValue(model.VersionsTotal)
	}
	if model.Versions != nil {
		modelMap["versions"] = model.Versions
	}
	if model.Issuer != nil {
		modelMap["issuer"] = model.Issuer
	}
	if model.BundleCerts != nil {
		modelMap["bundle_certs"] = model.BundleCerts
	}
	if model.Ca != nil {
		modelMap["ca"] = model.Ca
	}
	if model.DNS != nil {
		modelMap["dns"] = model.DNS
	}
	if model.Algorithm != nil {
		modelMap["algorithm"] = model.Algorithm
	}
	if model.KeyAlgorithm != nil {
		modelMap["key_algorithm"] = model.KeyAlgorithm
	}
	if model.AltNames != nil {
		modelMap["alt_names"] = model.AltNames
	}
	if model.CommonName != nil {
		modelMap["common_name"] = model.CommonName
	}
	if model.PrivateKeyIncluded != nil {
		modelMap["private_key_included"] = model.PrivateKeyIncluded
	}
	if model.IntermediateIncluded != nil {
		modelMap["intermediate_included"] = model.IntermediateIncluded
	}
	if model.Rotation != nil {
		rotationMap, err := ResourceIBMSmSecretRotationToMap(model.Rotation)
		if err != nil {
			return modelMap, err
		}
		modelMap["rotation"] = []map[string]interface{}{rotationMap}
	}
	if model.IssuanceInfo != nil {
		issuanceInfoMap, err := ResourceIBMSmSecretIssuanceInfoToMap(model.IssuanceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["issuance_info"] = []map[string]interface{}{issuanceInfoMap}
	}
	if model.Validity != nil {
		validityMap, err := ResourceIBMSmSecretCertificateValidityToMap(model.Validity)
		if err != nil {
			return modelMap, err
		}
		modelMap["validity"] = []map[string]interface{}{validityMap}
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = model.SerialNumber
	}
	if model.SecretData != nil {
		modelMap["secret_data"] = model.SecretData
	}
	return modelMap, nil
}

func ResourceIBMSmSecretPrivateCertificateSecretResourceToMap(model *secretsmanagerv1.PrivateCertificateSecretResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = flex.IntValue(model.State)
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = flex.IntValue(model.VersionsTotal)
	}
	if model.Versions != nil {
		modelMap["versions"] = model.Versions
	}
	modelMap["certificate_template"] = model.CertificateTemplate
	if model.CertificateAuthority != nil {
		modelMap["certificate_authority"] = model.CertificateAuthority
	}
	modelMap["common_name"] = model.CommonName
	if model.AltNames != nil {
		modelMap["alt_names"] = model.AltNames
	}
	if model.IPSans != nil {
		modelMap["ip_sans"] = model.IPSans
	}
	if model.URISans != nil {
		modelMap["uri_sans"] = model.URISans
	}
	if model.OtherSans != nil {
		modelMap["other_sans"] = model.OtherSans
	}
	if model.TTL != nil {
		modelMap["ttl"] = model.TTL
	}
	if model.Format != nil {
		modelMap["format"] = model.Format
	}
	if model.PrivateKeyFormat != nil {
		modelMap["private_key_format"] = model.PrivateKeyFormat
	}
	if model.ExcludeCnFromSans != nil {
		modelMap["exclude_cn_from_sans"] = model.ExcludeCnFromSans
	}
	if model.Rotation != nil {
		rotationMap, err := ResourceIBMSmSecretRotationToMap(model.Rotation)
		if err != nil {
			return modelMap, err
		}
		modelMap["rotation"] = []map[string]interface{}{rotationMap}
	}
	if model.Algorithm != nil {
		modelMap["algorithm"] = model.Algorithm
	}
	if model.KeyAlgorithm != nil {
		modelMap["key_algorithm"] = model.KeyAlgorithm
	}
	if model.Issuer != nil {
		modelMap["issuer"] = model.Issuer
	}
	if model.Validity != nil {
		validityMap, err := ResourceIBMSmSecretCertificateValidityToMap(model.Validity)
		if err != nil {
			return modelMap, err
		}
		modelMap["validity"] = []map[string]interface{}{validityMap}
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = model.SerialNumber
	}
	if model.RevocationTime != nil {
		modelMap["revocation_time"] = flex.IntValue(model.RevocationTime)
	}
	if model.RevocationTimeRfc3339 != nil {
		modelMap["revocation_time_rfc3339"] = model.RevocationTimeRfc3339.String()
	}
	if model.SecretData != nil {
		modelMap["secret_data"] = model.SecretData
	}
	return modelMap, nil
}

func ResourceIBMSmSecretKvSecretResourceToMap(model *secretsmanagerv1.KvSecretResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = model.SecretGroupID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.State != nil {
		modelMap["state"] = flex.IntValue(model.State)
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = model.StateDescription
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = model.SecretType
	}
	if model.CRN != nil {
		modelMap["crn"] = model.CRN
	}
	if model.CreationDate != nil {
		modelMap["creation_date"] = model.CreationDate.String()
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = model.CreatedBy
	}
	if model.LastUpdateDate != nil {
		modelMap["last_update_date"] = model.LastUpdateDate.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = flex.IntValue(model.VersionsTotal)
	}
	if model.Versions != nil {
		modelMap["versions"] = model.Versions
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.Payload != nil {
		modelMap["payload"] = model.Payload
	}
	if model.SecretData != nil {
		modelMap["secret_data"] = model.SecretData
	}
	return modelMap, nil
}
