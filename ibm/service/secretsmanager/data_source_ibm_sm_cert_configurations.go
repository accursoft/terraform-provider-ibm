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

func DataSourceIBMSmCertConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMSmCertConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The secret type. Allowable values are: public_cert, private_cert",
			},
			"config_element": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The configuration element to define or manage. Allowable values are: certificate_authorities, dns_providers, root_certificate_authorities, intermediate_certificate_authorities, certificate_templates",
			},
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
						"certificate_authorities": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The human-readable name to assign to your configuration.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of configuration. Value options differ depending on the 'config_element'. Allowable values are: letsencrypt, letsencrypt-stage, cis, classic_infrastructure, root_certificate_authority, intermediate_certificate_authority",
									},
								},
							},
						},
						"dns_providers": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The human-readable name to assign to your configuration.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of configuration. Value options differ depending on the 'config_element'. Allowable values are: letsencrypt, letsencrypt-stage, cis, classic_infrastructure, root_certificate_authority, intermediate_certificate_authority",
									},
								},
							},
						},
						"root_certificate_authorities": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The human-readable name to assign to your configuration.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of configuration. Value options differ depending on the 'config_element'. Allowable values are: letsencrypt, letsencrypt-stage, cis, classic_infrastructure, root_certificate_authority, intermediate_certificate_authority",
									},
									"config": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Root certificate authority configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max_ttl": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "The maximum time-to-live (TTL) for certificates that are created by this CA.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).",
												},
												"crl_expiry": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "The time until the certificate revocation list (CRL) expires.The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours. In the API response, this value is returned in seconds (integer).**Note:** The CRL is rotated automatically before it expires.",
												},
												"crl_disable": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Disables or enables certificate revocation list (CRL) building. If CRL building is disabled, a signed but zero-length CRL is returned when downloading the CRL. If CRL building is enabled,  it will rebuild the CRL.",
												},
												"crl_distribution_points_encoded": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to encode the certificate revocation list (CRL) distribution points in the certificates that are issued by this certificate authority.",
												},
												"issuing_certificates_urls_encoded": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to encode the URL of the issuing certificate in the certificates that are issued by this certificate authority.",
												},
												"common_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The fully qualified domain name or host domain name for the certificate.",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The status of the certificate authority. The status of a root certificate authority is either `configured` or `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,`signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.",
												},
												"expiration_date": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The date that the certificate expires. The date format follows RFC 3339.",
												},
												"alt_names": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.The alternative names can be host names or email addresses.",
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
												"ttl": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "The time-to-live (TTL) to assign to this CA certificate.The value can be supplied as a string representation of a duration, such as `12h`. The value can't exceed the `max_ttl` that is defined in the associated certificate template. In the API response, this value is returned in seconds (integer).",
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
												"key_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type of private key to generate.",
												},
												"key_bits": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The number of bits to use when generating the private key.Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.",
												},
												"max_path_length": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The maximum path length to encode in the generated certificate. `-1` means no limit.If the signing certificate has a maximum path length set, the path length is set to one less than that of the signing certificate. A limit of `0` means a literal path length of zero.",
												},
												"exclude_cn_from_sans": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Controls whether the common name is excluded from Subject Alternative Names (SANs). If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if the common name is not a hostname or an email address, but is instead a human-readable identifier.",
												},
												"permitted_dns_domains": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ou": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Organizational Unit (OU) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"organization": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Organization (O) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"country": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Country (C) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"locality": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Locality (L) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"province": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Province (ST) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"street_address": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Street Address values in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"postal_code": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Postal Code values in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"serial_number": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.",
												},
												"data": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "The data that is associated with the root certificate authority. The data object contains the following fields:- `certificate`: The root certificate content.- `issuing_ca`: The certificate of the certificate authority that signed and issued this certificate.- `serial_number`: The unique serial number of the root certificate.",
												},
											},
										},
									},
								},
							},
						},
						"intermediate_certificate_authorities": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The human-readable name to assign to your configuration.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of configuration. Value options differ depending on the 'config_element'. Allowable values are: letsencrypt, letsencrypt-stage, cis, classic_infrastructure, root_certificate_authority, intermediate_certificate_authority",
									},
									"config": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Intermediate certificate authority configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max_ttl": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "The maximum time-to-live (TTL) for certificates that are created by this CA.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).",
												},
												"signing_method": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The signing method to use with this certificate authority to generate private certificates.You can choose between internal or externally signed options. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).",
												},
												"issuer": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The certificate authority that signed and issued the certificate.If the certificate is signed internally, the `issuer` field is required and must match the name of a certificate authority that is configured in the Secrets Manager service instance.",
												},
												"crl_expiry": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "The time until the certificate revocation list (CRL) expires.The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours. In the API response, this value is returned in seconds (integer).**Note:** The CRL is rotated automatically before it expires.",
												},
												"crl_disable": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Disables or enables certificate revocation list (CRL) building. If CRL building is disabled, a signed but zero-length CRL is returned when downloading the CRL. If CRL building is enabled,  it will rebuild the CRL.",
												},
												"crl_distribution_points_encoded": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to encode the certificate revocation list (CRL) distribution points in the certificates that are issued by this certificate authority.",
												},
												"issuing_certificates_urls_encoded": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to encode the URL of the issuing certificate in the certificates that are issued by this certificate authority.",
												},
												"common_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The fully qualified domain name or host domain name for the certificate.",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The status of the certificate authority. The status of a root certificate authority is either `configured` or `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,`signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.",
												},
												"expiration_date": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The date that the certificate expires. The date format follows RFC 3339.",
												},
												"alt_names": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.The alternative names can be host names or email addresses.",
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
												"key_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type of private key to generate.",
												},
												"key_bits": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The number of bits to use when generating the private key.Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.",
												},
												"exclude_cn_from_sans": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Controls whether the common name is excluded from Subject Alternative Names (SANs). If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if the common name is not a hostname or an email address, but is instead a human-readable identifier.",
												},
												"ou": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Organizational Unit (OU) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"organization": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Organization (O) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"country": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Country (C) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"locality": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Locality (L) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"province": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Province (ST) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"street_address": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Street Address values in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"postal_code": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Postal Code values in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"serial_number": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.",
												},
												"data": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "The data that is associated with the intermediate certificate authority. The data object contains the following fields:- `csr`: The PEM-encoded certificate signing request.- `private_key`: The private key.- `private_key_type`: The type of private key, for example `rsa`.",
												},
											},
										},
									},
								},
							},
						},
						"certificate_templates": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The human-readable name to assign to your configuration.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of configuration. Value options differ depending on the 'config_element'. Allowable values are: letsencrypt, letsencrypt-stage, cis, classic_infrastructure, root_certificate_authority, intermediate_certificate_authority",
									},
									"config": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Properties that describe a certificate template. You can use a certificate template to control the parameters that are applied to your issued private certificates. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-certificate-templates).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"certificate_authority": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the intermediate certificate authority.",
												},
												"allowed_secret_groups": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Scopes the creation of private certificates to only the secret groups that you specify. This field can be supplied as a comma-delimited list of secret group IDs.",
												},
												"max_ttl": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "The maximum time-to-live (TTL) for certificates that are created by this CA.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).",
												},
												"ttl": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "The time-to-live (TTL) to assign to a private certificate.The value can be supplied as a string representation of a duration, such as `12h`. Hour (`h`) is the largest time suffix. The value can't exceed the `max_ttl` that is defined in the associated certificate template. In the API response, this value is returned in seconds (integer).",
												},
												"allow_localhost": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to allow `localhost` to be included as one of the requested common names.",
												},
												"allowed_domains": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and `allow_subdomains` options.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"allowed_domains_template": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to allow the domains that are supplied in the `allowed_domains` field to contain access control list (ACL) templates.",
												},
												"allow_bare_domains": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to allow clients to request private certificates that match the value of the actual domains on the final certificate.For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a certificate that contains the name `example.com` as one of the DNS values on the final certificate.**Important:** In some scenarios, allowing bare domains can be considered a security risk.",
												},
												"allow_subdomains": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to allow clients to request private certificates with common names (CN) that are subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.**Note:** This field is redundant if you use the `allow_any_name` option.",
												},
												"allow_glob_domains": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are specified in the `allowed_domains` field.If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.",
												},
												"allow_any_name": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to allow clients to request a private certificate that matches any common name.",
												},
												"enforce_hostnames": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to enforce only valid host names for common names, DNS Subject Alternative Names, and the host section of email addresses.",
												},
												"allow_ip_sans": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to allow clients to request a private certificate with IP Subject Alternative Names.",
												},
												"allowed_uri_sans": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The URI Subject Alternative Names to allow for private certificates.Values can contain glob patterns, for example `spiffe://hostname/_*`.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"allowed_other_sans": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private certificates. The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any `other_sans` input.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"server_flag": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether private certificates are flagged for server use.",
												},
												"client_flag": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether private certificates are flagged for client use.",
												},
												"code_signing_flag": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether private certificates are flagged for code signing use.",
												},
												"email_protection_flag": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether private certificates are flagged for email protection use.",
												},
												"key_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type of private key to generate for private certificates and the type of key that is expected for submitted certificate signing requests (CSRs). Allowable values are: `rsa` and `ec`.",
												},
												"key_bits": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The number of bits to use when generating the private key.Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.",
												},
												"key_usage": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The allowed key usage constraint to define for private certificates.You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage).  Omit the `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to an empty list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ext_key_usage": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The allowed extended key usage constraint on private certificates.You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage). Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to an empty list.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ext_key_usage_oids": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A list of extended key usage Object Identifiers (OIDs).",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"use_csr_common_name": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "When used with the `sign_csr` action, this field determines whether to use the common name (CN) from a certificate signing request (CSR) instead of the CN that's included in the JSON data of the certificate.Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include the `use_csr_sans` property.",
												},
												"use_csr_sans": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "When used with the `sign_csr` action, this field determines whether to use the Subject Alternative Names (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the JSON data of the certificate.Does not include the common name in the CSR. To use the common name, include the `use_csr_common_name` property.",
												},
												"ou": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Organizational Unit (OU) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"organization": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Organization (O) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"country": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Country (C) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"locality": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Locality (L) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"province": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Province (ST) values to define in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"street_address": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Street Address values in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"postal_code": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The Postal Code values in the subject field of the resulting certificate.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"serial_number": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.",
												},
												"require_cn": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to require a common name to create a private certificate.By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the `require_cn` option to `false`.",
												},
												"policy_identifiers": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A list of policy Object Identifiers (OIDs).",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"basic_constraints_valid_for_non_ca": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Determines whether to mark the Basic Constraints extension of an issued private certificate as valid for non-CA certificates.",
												},
												"not_before_duration": &schema.Schema{
													Type:        schema.TypeMap,
													Computed:    true,
													Description: "The duration in seconds by which to backdate the `not_before` property of an issued private certificate.The value can be supplied as a string representation of a duration, such as `30s`. In the API response, this value is returned in seconds (integer).",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMSmCertConfigurationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getConfigElementsOptions := &secretsmanagerv1.GetConfigElementsOptions{}

	getConfigElementsOptions.SetSecretType(d.Get("secret_type").(string))
	getConfigElementsOptions.SetConfigElement(d.Get("config_element").(string))

	getConfigElements, response, err := secretsManagerClient.GetConfigElementsWithContext(context, getConfigElementsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetConfigElementsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetConfigElementsWithContext failed %s\n%s", err, response))
	}

	d.SetId(DataSourceIBMSmCertConfigurationsID(d))

	metadata := []map[string]interface{}{}
	if getConfigElements.Metadata != nil {
		modelMap, err := DataSourceIBMSmCertConfigurationsCollectionMetadataToMap(getConfigElements.Metadata)
		if err != nil {
			return diag.FromErr(err)
		}
		metadata = append(metadata, modelMap)
	}
	if err = d.Set("metadata", metadata); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting metadata %s", err))
	}

	resources := []map[string]interface{}{}
	if getConfigElements.Resources != nil {
		for _, modelItem := range getConfigElements.Resources { 
			modelMap, err := DataSourceIBMSmCertConfigurationsGetConfigElementsResourcesItemToMap(modelItem)
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

// DataSourceIBMSmCertConfigurationsID returns a reasonable ID for the list.
func DataSourceIBMSmCertConfigurationsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMSmCertConfigurationsCollectionMetadataToMap(model *secretsmanagerv1.CollectionMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CollectionType != nil {
		modelMap["collection_type"] = *model.CollectionType
	}
	if model.CollectionTotal != nil {
		modelMap["collection_total"] = *model.CollectionTotal
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsGetConfigElementsResourcesItemToMap(model secretsmanagerv1.GetConfigElementsResourcesItemIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv1.GetConfigElementsResourcesItemCertificateAuthoritiesConfig); ok {
		return DataSourceIBMSmCertConfigurationsGetConfigElementsResourcesItemCertificateAuthoritiesConfigToMap(model.(*secretsmanagerv1.GetConfigElementsResourcesItemCertificateAuthoritiesConfig))
	} else if _, ok := model.(*secretsmanagerv1.GetConfigElementsResourcesItemDNSProvidersConfig); ok {
		return DataSourceIBMSmCertConfigurationsGetConfigElementsResourcesItemDNSProvidersConfigToMap(model.(*secretsmanagerv1.GetConfigElementsResourcesItemDNSProvidersConfig))
	} else if _, ok := model.(*secretsmanagerv1.RootCertificateAuthoritiesConfig); ok {
		return DataSourceIBMSmCertConfigurationsRootCertificateAuthoritiesConfigToMap(model.(*secretsmanagerv1.RootCertificateAuthoritiesConfig))
	} else if _, ok := model.(*secretsmanagerv1.IntermediateCertificateAuthoritiesConfig); ok {
		return DataSourceIBMSmCertConfigurationsIntermediateCertificateAuthoritiesConfigToMap(model.(*secretsmanagerv1.IntermediateCertificateAuthoritiesConfig))
	} else if _, ok := model.(*secretsmanagerv1.CertificateTemplatesConfig); ok {
		return DataSourceIBMSmCertConfigurationsCertificateTemplatesConfigToMap(model.(*secretsmanagerv1.CertificateTemplatesConfig))
	} else if _, ok := model.(*secretsmanagerv1.GetConfigElementsResourcesItem); ok {
		modelMap := make(map[string]interface{})
		model := model.(*secretsmanagerv1.GetConfigElementsResourcesItem)
		if model.CertificateAuthorities != nil {
			certificateAuthorities := []map[string]interface{}{}
			for _, certificateAuthoritiesItem := range model.CertificateAuthorities {
				certificateAuthoritiesItemMap, err := DataSourceIBMSmCertConfigurationsConfigElementMetadataToMap(&certificateAuthoritiesItem)
				if err != nil {
					return modelMap, err
				}
				certificateAuthorities = append(certificateAuthorities, certificateAuthoritiesItemMap)
			}
			modelMap["certificate_authorities"] = certificateAuthorities
		}
		if model.DNSProviders != nil {
			dnsProviders := []map[string]interface{}{}
			for _, dnsProvidersItem := range model.DNSProviders {
				dnsProvidersItemMap, err := DataSourceIBMSmCertConfigurationsConfigElementMetadataToMap(&dnsProvidersItem)
				if err != nil {
					return modelMap, err
				}
				dnsProviders = append(dnsProviders, dnsProvidersItemMap)
			}
			modelMap["dns_providers"] = dnsProviders
		}
		if model.RootCertificateAuthorities != nil {
			rootCertificateAuthorities := []map[string]interface{}{}
			for _, rootCertificateAuthoritiesItem := range model.RootCertificateAuthorities {
				rootCertificateAuthoritiesItemMap, err := DataSourceIBMSmCertConfigurationsRootCertificateAuthoritiesConfigItemToMap(&rootCertificateAuthoritiesItem)
				if err != nil {
					return modelMap, err
				}
				rootCertificateAuthorities = append(rootCertificateAuthorities, rootCertificateAuthoritiesItemMap)
			}
			modelMap["root_certificate_authorities"] = rootCertificateAuthorities
		}
		if model.IntermediateCertificateAuthorities != nil {
			intermediateCertificateAuthorities := []map[string]interface{}{}
			for _, intermediateCertificateAuthoritiesItem := range model.IntermediateCertificateAuthorities {
				intermediateCertificateAuthoritiesItemMap, err := DataSourceIBMSmCertConfigurationsIntermediateCertificateAuthoritiesConfigItemToMap(&intermediateCertificateAuthoritiesItem)
				if err != nil {
					return modelMap, err
				}
				intermediateCertificateAuthorities = append(intermediateCertificateAuthorities, intermediateCertificateAuthoritiesItemMap)
			}
			modelMap["intermediate_certificate_authorities"] = intermediateCertificateAuthorities
		}
		if model.CertificateTemplates != nil {
			certificateTemplates := []map[string]interface{}{}
			for _, certificateTemplatesItem := range model.CertificateTemplates {
				certificateTemplatesItemMap, err := DataSourceIBMSmCertConfigurationsCertificateTemplatesConfigItemToMap(&certificateTemplatesItem)
				if err != nil {
					return modelMap, err
				}
				certificateTemplates = append(certificateTemplates, certificateTemplatesItemMap)
			}
			modelMap["certificate_templates"] = certificateTemplates
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv1.GetConfigElementsResourcesItemIntf subtype encountered")
	}
}

func DataSourceIBMSmCertConfigurationsConfigElementMetadataToMap(model *secretsmanagerv1.ConfigElementMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsRootCertificateAuthoritiesConfigItemToMap(model *secretsmanagerv1.RootCertificateAuthoritiesConfigItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Config != nil {
		configMap, err := DataSourceIBMSmCertConfigurationsRootCertificateAuthorityConfigToMap(model.Config)
		if err != nil {
			return modelMap, err
		}
		modelMap["config"] = []map[string]interface{}{configMap}
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsRootCertificateAuthorityConfigToMap(model *secretsmanagerv1.RootCertificateAuthorityConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxTTL != nil {
	}
	if model.CrlExpiry != nil {
	}
	if model.CrlDisable != nil {
		modelMap["crl_disable"] = *model.CrlDisable
	}
	if model.CrlDistributionPointsEncoded != nil {
		modelMap["crl_distribution_points_encoded"] = *model.CrlDistributionPointsEncoded
	}
	if model.IssuingCertificatesUrlsEncoded != nil {
		modelMap["issuing_certificates_urls_encoded"] = *model.IssuingCertificatesUrlsEncoded
	}
	if model.CommonName != nil {
		modelMap["common_name"] = *model.CommonName
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.AltNames != nil {
		modelMap["alt_names"] = *model.AltNames
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
	if model.KeyType != nil {
		modelMap["key_type"] = *model.KeyType
	}
	if model.KeyBits != nil {
		modelMap["key_bits"] = *model.KeyBits
	}
	if model.MaxPathLength != nil {
		modelMap["max_path_length"] = *model.MaxPathLength
	}
	if model.ExcludeCnFromSans != nil {
		modelMap["exclude_cn_from_sans"] = *model.ExcludeCnFromSans
	}
	if model.PermittedDNSDomains != nil {
		modelMap["permitted_dns_domains"] = model.PermittedDNSDomains
	}
	if model.Ou != nil {
		modelMap["ou"] = model.Ou
	}
	if model.Organization != nil {
		modelMap["organization"] = model.Organization
	}
	if model.Country != nil {
		modelMap["country"] = model.Country
	}
	if model.Locality != nil {
		modelMap["locality"] = model.Locality
	}
	if model.Province != nil {
		modelMap["province"] = model.Province
	}
	if model.StreetAddress != nil {
		modelMap["street_address"] = model.StreetAddress
	}
	if model.PostalCode != nil {
		modelMap["postal_code"] = model.PostalCode
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = *model.SerialNumber
	}
	if model.Data != nil {
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsIntermediateCertificateAuthoritiesConfigItemToMap(model *secretsmanagerv1.IntermediateCertificateAuthoritiesConfigItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Config != nil {
		configMap, err := DataSourceIBMSmCertConfigurationsIntermediateCertificateAuthorityConfigToMap(model.Config)
		if err != nil {
			return modelMap, err
		}
		modelMap["config"] = []map[string]interface{}{configMap}
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsIntermediateCertificateAuthorityConfigToMap(model *secretsmanagerv1.IntermediateCertificateAuthorityConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxTTL != nil {
	}
	if model.SigningMethod != nil {
		modelMap["signing_method"] = *model.SigningMethod
	}
	if model.Issuer != nil {
		modelMap["issuer"] = *model.Issuer
	}
	if model.CrlExpiry != nil {
	}
	if model.CrlDisable != nil {
		modelMap["crl_disable"] = *model.CrlDisable
	}
	if model.CrlDistributionPointsEncoded != nil {
		modelMap["crl_distribution_points_encoded"] = *model.CrlDistributionPointsEncoded
	}
	if model.IssuingCertificatesUrlsEncoded != nil {
		modelMap["issuing_certificates_urls_encoded"] = *model.IssuingCertificatesUrlsEncoded
	}
	if model.CommonName != nil {
		modelMap["common_name"] = *model.CommonName
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.AltNames != nil {
		modelMap["alt_names"] = *model.AltNames
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
	if model.KeyType != nil {
		modelMap["key_type"] = *model.KeyType
	}
	if model.KeyBits != nil {
		modelMap["key_bits"] = *model.KeyBits
	}
	if model.ExcludeCnFromSans != nil {
		modelMap["exclude_cn_from_sans"] = *model.ExcludeCnFromSans
	}
	if model.Ou != nil {
		modelMap["ou"] = model.Ou
	}
	if model.Organization != nil {
		modelMap["organization"] = model.Organization
	}
	if model.Country != nil {
		modelMap["country"] = model.Country
	}
	if model.Locality != nil {
		modelMap["locality"] = model.Locality
	}
	if model.Province != nil {
		modelMap["province"] = model.Province
	}
	if model.StreetAddress != nil {
		modelMap["street_address"] = model.StreetAddress
	}
	if model.PostalCode != nil {
		modelMap["postal_code"] = model.PostalCode
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = *model.SerialNumber
	}
	if model.Data != nil {
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsCertificateTemplatesConfigItemToMap(model *secretsmanagerv1.CertificateTemplatesConfigItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Config != nil {
		configMap, err := DataSourceIBMSmCertConfigurationsCertificateTemplateConfigToMap(model.Config)
		if err != nil {
			return modelMap, err
		}
		modelMap["config"] = []map[string]interface{}{configMap}
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsCertificateTemplateConfigToMap(model *secretsmanagerv1.CertificateTemplateConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CertificateAuthority != nil {
		modelMap["certificate_authority"] = *model.CertificateAuthority
	}
	if model.AllowedSecretGroups != nil {
		modelMap["allowed_secret_groups"] = *model.AllowedSecretGroups
	}
	if model.MaxTTL != nil {
	}
	if model.TTL != nil {
	}
	if model.AllowLocalhost != nil {
		modelMap["allow_localhost"] = *model.AllowLocalhost
	}
	if model.AllowedDomains != nil {
		modelMap["allowed_domains"] = model.AllowedDomains
	}
	if model.AllowedDomainsTemplate != nil {
		modelMap["allowed_domains_template"] = *model.AllowedDomainsTemplate
	}
	if model.AllowBareDomains != nil {
		modelMap["allow_bare_domains"] = *model.AllowBareDomains
	}
	if model.AllowSubdomains != nil {
		modelMap["allow_subdomains"] = *model.AllowSubdomains
	}
	if model.AllowGlobDomains != nil {
		modelMap["allow_glob_domains"] = *model.AllowGlobDomains
	}
	if model.AllowAnyName != nil {
		modelMap["allow_any_name"] = *model.AllowAnyName
	}
	if model.EnforceHostnames != nil {
		modelMap["enforce_hostnames"] = *model.EnforceHostnames
	}
	if model.AllowIPSans != nil {
		modelMap["allow_ip_sans"] = *model.AllowIPSans
	}
	if model.AllowedURISans != nil {
		modelMap["allowed_uri_sans"] = model.AllowedURISans
	}
	if model.AllowedOtherSans != nil {
		modelMap["allowed_other_sans"] = model.AllowedOtherSans
	}
	if model.ServerFlag != nil {
		modelMap["server_flag"] = *model.ServerFlag
	}
	if model.ClientFlag != nil {
		modelMap["client_flag"] = *model.ClientFlag
	}
	if model.CodeSigningFlag != nil {
		modelMap["code_signing_flag"] = *model.CodeSigningFlag
	}
	if model.EmailProtectionFlag != nil {
		modelMap["email_protection_flag"] = *model.EmailProtectionFlag
	}
	if model.KeyType != nil {
		modelMap["key_type"] = *model.KeyType
	}
	if model.KeyBits != nil {
		modelMap["key_bits"] = *model.KeyBits
	}
	if model.KeyUsage != nil {
		modelMap["key_usage"] = model.KeyUsage
	}
	if model.ExtKeyUsage != nil {
		modelMap["ext_key_usage"] = model.ExtKeyUsage
	}
	if model.ExtKeyUsageOids != nil {
		modelMap["ext_key_usage_oids"] = model.ExtKeyUsageOids
	}
	if model.UseCsrCommonName != nil {
		modelMap["use_csr_common_name"] = *model.UseCsrCommonName
	}
	if model.UseCsrSans != nil {
		modelMap["use_csr_sans"] = *model.UseCsrSans
	}
	if model.Ou != nil {
		modelMap["ou"] = model.Ou
	}
	if model.Organization != nil {
		modelMap["organization"] = model.Organization
	}
	if model.Country != nil {
		modelMap["country"] = model.Country
	}
	if model.Locality != nil {
		modelMap["locality"] = model.Locality
	}
	if model.Province != nil {
		modelMap["province"] = model.Province
	}
	if model.StreetAddress != nil {
		modelMap["street_address"] = model.StreetAddress
	}
	if model.PostalCode != nil {
		modelMap["postal_code"] = model.PostalCode
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = *model.SerialNumber
	}
	if model.RequireCn != nil {
		modelMap["require_cn"] = *model.RequireCn
	}
	if model.PolicyIdentifiers != nil {
		modelMap["policy_identifiers"] = model.PolicyIdentifiers
	}
	if model.BasicConstraintsValidForNonCa != nil {
		modelMap["basic_constraints_valid_for_non_ca"] = *model.BasicConstraintsValidForNonCa
	}
	if model.NotBeforeDuration != nil {
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsGetConfigElementsResourcesItemCertificateAuthoritiesConfigToMap(model *secretsmanagerv1.GetConfigElementsResourcesItemCertificateAuthoritiesConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CertificateAuthorities != nil {
		certificateAuthorities := []map[string]interface{}{}
		for _, certificateAuthoritiesItem := range model.CertificateAuthorities {
			certificateAuthoritiesItemMap, err := DataSourceIBMSmCertConfigurationsConfigElementMetadataToMap(&certificateAuthoritiesItem)
			if err != nil {
				return modelMap, err
			}
			certificateAuthorities = append(certificateAuthorities, certificateAuthoritiesItemMap)
		}
		modelMap["certificate_authorities"] = certificateAuthorities
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsGetConfigElementsResourcesItemDNSProvidersConfigToMap(model *secretsmanagerv1.GetConfigElementsResourcesItemDNSProvidersConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DNSProviders != nil {
		dnsProviders := []map[string]interface{}{}
		for _, dnsProvidersItem := range model.DNSProviders {
			dnsProvidersItemMap, err := DataSourceIBMSmCertConfigurationsConfigElementMetadataToMap(&dnsProvidersItem)
			if err != nil {
				return modelMap, err
			}
			dnsProviders = append(dnsProviders, dnsProvidersItemMap)
		}
		modelMap["dns_providers"] = dnsProviders
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsRootCertificateAuthoritiesConfigToMap(model *secretsmanagerv1.RootCertificateAuthoritiesConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RootCertificateAuthorities != nil {
		rootCertificateAuthorities := []map[string]interface{}{}
		for _, rootCertificateAuthoritiesItem := range model.RootCertificateAuthorities {
			rootCertificateAuthoritiesItemMap, err := DataSourceIBMSmCertConfigurationsRootCertificateAuthoritiesConfigItemToMap(&rootCertificateAuthoritiesItem)
			if err != nil {
				return modelMap, err
			}
			rootCertificateAuthorities = append(rootCertificateAuthorities, rootCertificateAuthoritiesItemMap)
		}
		modelMap["root_certificate_authorities"] = rootCertificateAuthorities
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsIntermediateCertificateAuthoritiesConfigToMap(model *secretsmanagerv1.IntermediateCertificateAuthoritiesConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IntermediateCertificateAuthorities != nil {
		intermediateCertificateAuthorities := []map[string]interface{}{}
		for _, intermediateCertificateAuthoritiesItem := range model.IntermediateCertificateAuthorities {
			intermediateCertificateAuthoritiesItemMap, err := DataSourceIBMSmCertConfigurationsIntermediateCertificateAuthoritiesConfigItemToMap(&intermediateCertificateAuthoritiesItem)
			if err != nil {
				return modelMap, err
			}
			intermediateCertificateAuthorities = append(intermediateCertificateAuthorities, intermediateCertificateAuthoritiesItemMap)
		}
		modelMap["intermediate_certificate_authorities"] = intermediateCertificateAuthorities
	}
	return modelMap, nil
}

func DataSourceIBMSmCertConfigurationsCertificateTemplatesConfigToMap(model *secretsmanagerv1.CertificateTemplatesConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CertificateTemplates != nil {
		certificateTemplates := []map[string]interface{}{}
		for _, certificateTemplatesItem := range model.CertificateTemplates {
			certificateTemplatesItemMap, err := DataSourceIBMSmCertConfigurationsCertificateTemplatesConfigItemToMap(&certificateTemplatesItem)
			if err != nil {
				return modelMap, err
			}
			certificateTemplates = append(certificateTemplates, certificateTemplatesItemMap)
		}
		modelMap["certificate_templates"] = certificateTemplates
	}
	return modelMap, nil
}
