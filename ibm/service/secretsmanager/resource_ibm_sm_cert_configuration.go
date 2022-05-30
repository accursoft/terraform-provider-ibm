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

func ResourceIBMSmCertConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext:   ResourceIBMSmCertConfigurationCreate,
		ReadContext:     ResourceIBMSmCertConfigurationRead,
		UpdateContext:   ResourceIBMSmCertConfigurationUpdate,
		DeleteContext:   ResourceIBMSmCertConfigurationDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				ValidateFunc: validate.InvokeValidator("ibm_sm_cert_configuration", "secret_type"),
				Description: "The secret type. Allowable values are: public_cert, private_cert",
			},
			"config_element": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				ValidateFunc: validate.InvokeValidator("ibm_sm_cert_configuration", "config_element"),
				Description: "The configuration element to define or manage. Allowable values are: certificate_authorities, dns_providers, root_certificate_authorities, intermediate_certificate_authorities, certificate_templates",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_sm_cert_configuration", "name"),
				Description: "The human-readable name to assign to your configuration.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_sm_cert_configuration", "type"),
				Description: "The type of configuration. Value options differ depending on the 'config_element'. Allowable values are: letsencrypt, letsencrypt-stage, cis, classic_infrastructure, root_certificate_authority, intermediate_certificate_authority",
			},
			"config": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The configuration to define for the specified secret type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_key": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The private key that is associated with your Automatic Certificate Management Environment (ACME) account.If you have a working ACME client or account for Let's Encrypt, you can use the existing private key to enable communications with Secrets Manager. If you don't have an account yet, you can create one. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#create-acme-account).",
						},
						"cis_crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Cloud Resource Name (CRN) that is associated with the CIS instance.",
						},
						"cis_apikey": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An IBM Cloud API key that can to list domains in your CIS instance.To grant Secrets Manager the ability to view the CIS instance and all of its domains, the API key must be assigned the Reader service role on Internet Services (`internet-svcs`).If you need to manage specific domains, you can assign the Manager role. For production environments, it is recommended that you assign the Reader access role, and then use the[IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific domains. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).",
						},
						"classic_infrastructure_username": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The username that is associated with your classic infrastructure account.In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information, see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).",
						},
						"classic_infrastructure_password": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Your classic infrastructure API key.For information about viewing and accessing your classic infrastructure API key, see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).",
						},
						"max_ttl": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The maximum time-to-live (TTL) for certificates that are created by this CA.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).",
						},
						"crl_expiry": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The time until the certificate revocation list (CRL) expires.The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours. In the API response, this value is returned in seconds (integer).**Note:** The CRL is rotated automatically before it expires.",
						},
						"crl_disable": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Disables or enables certificate revocation list (CRL) building. If CRL building is disabled, a signed but zero-length CRL is returned when downloading the CRL. If CRL building is enabled,  it will rebuild the CRL.",
						},
						"crl_distribution_points_encoded": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether to encode the certificate revocation list (CRL) distribution points in the certificates that are issued by this certificate authority.",
						},
						"issuing_certificates_urls_encoded": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether to encode the URL of the issuing certificate in the certificates that are issued by this certificate authority.",
						},
						"common_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The fully qualified domain name or host domain name for the certificate.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The status of the certificate authority. The status of a root certificate authority is either `configured` or `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,`signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.",
						},
						"expiration_date": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The date that the certificate expires. The date format follows RFC 3339.",
						},
						"alt_names": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.The alternative names can be host names or email addresses.",
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
						"ttl": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The time-to-live (TTL) to assign to this CA certificate.The value can be supplied as a string representation of a duration, such as `12h`. The value can't exceed the `max_ttl` that is defined in the associated certificate template. In the API response, this value is returned in seconds (integer).",
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
						"key_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of private key to generate.",
						},
						"key_bits": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of bits to use when generating the private key.Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.",
						},
						"max_path_length": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum path length to encode in the generated certificate. `-1` means no limit.If the signing certificate has a maximum path length set, the path length is set to one less than that of the signing certificate. A limit of `0` means a literal path length of zero.",
						},
						"exclude_cn_from_sans": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Controls whether the common name is excluded from Subject Alternative Names (SANs). If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if the common name is not a hostname or an email address, but is instead a human-readable identifier.",
						},
						"permitted_dns_domains": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ou": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Organizational Unit (OU) values to define in the subject field of the resulting certificate.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"organization": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Organization (O) values to define in the subject field of the resulting certificate.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"country": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Country (C) values to define in the subject field of the resulting certificate.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"locality": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Locality (L) values to define in the subject field of the resulting certificate.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"province": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Province (ST) values to define in the subject field of the resulting certificate.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"street_address": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Street Address values in the subject field of the resulting certificate.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"postal_code": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The Postal Code values in the subject field of the resulting certificate.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"serial_number": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.",
						},
						"data": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "The data that is associated with the root certificate authority. The data object contains the following fields:- `certificate`: The root certificate content.- `issuing_ca`: The certificate of the certificate authority that signed and issued this certificate.- `serial_number`: The unique serial number of the root certificate.",
						},
						"signing_method": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The signing method to use with this certificate authority to generate private certificates.You can choose between internal or externally signed options. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).",
						},
						"issuer": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The certificate authority that signed and issued the certificate.If the certificate is signed internally, the `issuer` field is required and must match the name of a certificate authority that is configured in the Secrets Manager service instance.",
						},
						"certificate_authority": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the intermediate certificate authority.",
						},
						"allowed_secret_groups": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Scopes the creation of private certificates to only the secret groups that you specify. This field can be supplied as a comma-delimited list of secret group IDs.",
						},
						"allow_localhost": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Determines whether to allow `localhost` to be included as one of the requested common names.",
						},
						"allowed_domains": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and `allow_subdomains` options.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"allowed_domains_template": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether to allow the domains that are supplied in the `allowed_domains` field to contain access control list (ACL) templates.",
						},
						"allow_bare_domains": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether to allow clients to request private certificates that match the value of the actual domains on the final certificate.For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a certificate that contains the name `example.com` as one of the DNS values on the final certificate.**Important:** In some scenarios, allowing bare domains can be considered a security risk.",
						},
						"allow_subdomains": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether to allow clients to request private certificates with common names (CN) that are subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.**Note:** This field is redundant if you use the `allow_any_name` option.",
						},
						"allow_glob_domains": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are specified in the `allowed_domains` field.If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.",
						},
						"allow_any_name": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether to allow clients to request a private certificate that matches any common name.",
						},
						"enforce_hostnames": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Determines whether to enforce only valid host names for common names, DNS Subject Alternative Names, and the host section of email addresses.",
						},
						"allow_ip_sans": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Determines whether to allow clients to request a private certificate with IP Subject Alternative Names.",
						},
						"allowed_uri_sans": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The URI Subject Alternative Names to allow for private certificates.Values can contain glob patterns, for example `spiffe://hostname/_*`.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"allowed_other_sans": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private certificates. The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any `other_sans` input.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"server_flag": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Determines whether private certificates are flagged for server use.",
						},
						"client_flag": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Determines whether private certificates are flagged for client use.",
						},
						"code_signing_flag": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether private certificates are flagged for code signing use.",
						},
						"email_protection_flag": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether private certificates are flagged for email protection use.",
						},
						"key_usage": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The allowed key usage constraint to define for private certificates.You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage).  Omit the `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to an empty list.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ext_key_usage": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The allowed extended key usage constraint on private certificates.You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage). Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to an empty list.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ext_key_usage_oids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of extended key usage Object Identifiers (OIDs).",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"use_csr_common_name": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "When used with the `sign_csr` action, this field determines whether to use the common name (CN) from a certificate signing request (CSR) instead of the CN that's included in the JSON data of the certificate.Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include the `use_csr_sans` property.",
						},
						"use_csr_sans": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "When used with the `sign_csr` action, this field determines whether to use the Subject Alternative Names (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the JSON data of the certificate.Does not include the common name in the CSR. To use the common name, include the `use_csr_common_name` property.",
						},
						"require_cn": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Determines whether to require a common name to create a private certificate.By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the `require_cn` option to `false`.",
						},
						"policy_identifiers": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of policy Object Identifiers (OIDs).",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"basic_constraints_valid_for_non_ca": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Determines whether to mark the Basic Constraints extension of an issued private certificate as valid for non-CA certificates.",
						},
						"not_before_duration": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The duration in seconds by which to backdate the `not_before` property of an issued private certificate.The value can be supplied as a string representation of a duration, such as `30s`. In the API response, this value is returned in seconds (integer).",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMSmCertConfigurationValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "secret_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "private_cert, public_cert",
		},
		validate.ValidateSchema{
			Identifier:                 "config_element",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "certificate_authorities, certificate_templates, dns_providers, intermediate_certificate_authorities, root_certificate_authorities",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.StringLenBetween,
			Type:                       validate.TypeString,
			Required:                   true,
			MinValueLength:             2,
			MaxValueLength:             256,
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "certificate_template, cis, classic_infrastructure, intermediate_certificate_authority, letsencrypt, letsencrypt-stage, root_certificate_authority",
			MinValueLength:             2,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_sm_cert_configuration", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIBMSmCertConfigurationCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createConfigElementOptions := &secretsmanagerv1.CreateConfigElementOptions{}

	createConfigElementOptions.SetSecretType(d.Get("secret_type").(string))
	createConfigElementOptions.SetConfigElement(d.Get("config_element").(string))
	createConfigElementOptions.SetName(d.Get("name").(string))
	createConfigElementOptions.SetType(d.Get("type").(string))
	configModel, err := ResourceIBMSmCertConfigurationMapToConfigElementDefConfig(d.Get("config.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createConfigElementOptions.SetConfig(configModel)

	configElementDef, response, err := secretsManagerClient.CreateConfigElementWithContext(context, createConfigElementOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigElementWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateConfigElementWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createConfigElementOptions.SecretType, *createConfigElementOptions.ConfigElement, *configElementDef.Name))

	return ResourceIBMSmCertConfigurationRead(context, d, meta)
}

func ResourceIBMSmCertConfigurationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getConfigElementOptions := &secretsmanagerv1.GetConfigElementOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getConfigElementOptions.SetSecretType(parts[0])
	getConfigElementOptions.SetConfigElement(parts[1])
	getConfigElementOptions.SetConfigName(parts[2])

	configElementDef, response, err := secretsManagerClient.GetConfigElementWithContext(context, getConfigElementOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetConfigElementWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetConfigElementWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("secret_type", getConfigElementOptions.SecretType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_type: %s", err))
	}
	if err = d.Set("config_element", getConfigElementOptions.ConfigElement); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting config_element: %s", err))
	}
	if err = d.Set("name", configElementDef.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("type", configElementDef.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	configMap, err := ResourceIBMSmCertConfigurationConfigElementDefConfigToMap(configElementDef.Config)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("config", []map[string]interface{}{configMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting config: %s", err))
	}

	return nil
}

func ResourceIBMSmCertConfigurationUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateConfigElementOptions := &secretsmanagerv1.UpdateConfigElementOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateConfigElementOptions.SetSecretType(parts[0])
	updateConfigElementOptions.SetConfigElement(parts[1])
	updateConfigElementOptions.SetConfigName(parts[2])

	hasChange := false

	if d.HasChange("secret_type") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation." +
				" The resource must be re-created to update this property.", "secret_type"))
	}
	if d.HasChange("config_element") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation." +
				" The resource must be re-created to update this property.", "config_element"))
	}
	if d.HasChange("name") || d.HasChange("type") || d.HasChange("config") {
		updateConfigElementOptions.SetName(d.Get("name").(string))
		updateConfigElementOptions.SetType(d.Get("type").(string))
		config, err := ResourceIBMSmCertConfigurationMapToConfigElementDefConfig(d.Get("config.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateConfigElementOptions.SetConfig(config)
		hasChange = true
	}

	if hasChange {
		_, response, err := secretsManagerClient.UpdateConfigElementWithContext(context, updateConfigElementOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateConfigElementWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateConfigElementWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIBMSmCertConfigurationRead(context, d, meta)
}

func ResourceIBMSmCertConfigurationDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteConfigElementOptions := &secretsmanagerv1.DeleteConfigElementOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteConfigElementOptions.SetSecretType(parts[0])
	deleteConfigElementOptions.SetConfigElement(parts[1])
	deleteConfigElementOptions.SetConfigName(parts[2])

	response, err := secretsManagerClient.DeleteConfigElementWithContext(context, deleteConfigElementOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteConfigElementWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteConfigElementWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func ResourceIBMSmCertConfigurationMapToConfigElementDefConfig(modelMap map[string]interface{}) (secretsmanagerv1.ConfigElementDefConfigIntf, error) {
	model := &secretsmanagerv1.ConfigElementDefConfig{}
	if modelMap["private_key"] != nil && modelMap["private_key"].(string) != "" {
		model.PrivateKey = core.StringPtr(modelMap["private_key"].(string))
	}
	if modelMap["cis_crn"] != nil && modelMap["cis_crn"].(string) != "" {
		model.CisCRN = core.StringPtr(modelMap["cis_crn"].(string))
	}
	if modelMap["cis_apikey"] != nil && modelMap["cis_apikey"].(string) != "" {
		model.CisApikey = core.StringPtr(modelMap["cis_apikey"].(string))
	}
	if modelMap["classic_infrastructure_username"] != nil && modelMap["classic_infrastructure_username"].(string) != "" {
		model.ClassicInfrastructureUsername = core.StringPtr(modelMap["classic_infrastructure_username"].(string))
	}
	if modelMap["classic_infrastructure_password"] != nil && modelMap["classic_infrastructure_password"].(string) != "" {
		model.ClassicInfrastructurePassword = core.StringPtr(modelMap["classic_infrastructure_password"].(string))
	}
	if modelMap["max_ttl"] != nil {
	
	}
	if modelMap["crl_expiry"] != nil {
	
	}
	if modelMap["crl_disable"] != nil {
		model.CrlDisable = core.BoolPtr(modelMap["crl_disable"].(bool))
	}
	if modelMap["crl_distribution_points_encoded"] != nil {
		model.CrlDistributionPointsEncoded = core.BoolPtr(modelMap["crl_distribution_points_encoded"].(bool))
	}
	if modelMap["issuing_certificates_urls_encoded"] != nil {
		model.IssuingCertificatesUrlsEncoded = core.BoolPtr(modelMap["issuing_certificates_urls_encoded"].(bool))
	}
	if modelMap["common_name"] != nil && modelMap["common_name"].(string) != "" {
		model.CommonName = core.StringPtr(modelMap["common_name"].(string))
	}
	if modelMap["status"] != nil && modelMap["status"].(string) != "" {
		model.Status = core.StringPtr(modelMap["status"].(string))
	}
	if modelMap["expiration_date"] != nil {
	
	}
	if modelMap["alt_names"] != nil && modelMap["alt_names"].(string) != "" {
		model.AltNames = core.StringPtr(modelMap["alt_names"].(string))
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
	if modelMap["key_type"] != nil && modelMap["key_type"].(string) != "" {
		model.KeyType = core.StringPtr(modelMap["key_type"].(string))
	}
	if modelMap["key_bits"] != nil {
		model.KeyBits = core.Int64Ptr(int64(modelMap["key_bits"].(int)))
	}
	if modelMap["max_path_length"] != nil {
		model.MaxPathLength = core.Int64Ptr(int64(modelMap["max_path_length"].(int)))
	}
	if modelMap["exclude_cn_from_sans"] != nil {
		model.ExcludeCnFromSans = core.BoolPtr(modelMap["exclude_cn_from_sans"].(bool))
	}
	if modelMap["permitted_dns_domains"] != nil {
		permittedDNSDomains := []string{}
		for _, permittedDNSDomainsItem := range modelMap["permitted_dns_domains"].([]interface{}) {
			permittedDNSDomains = append(permittedDNSDomains, permittedDNSDomainsItem.(string))
		}
		model.PermittedDNSDomains = permittedDNSDomains
	}
	if modelMap["ou"] != nil {
		ou := []string{}
		for _, ouItem := range modelMap["ou"].([]interface{}) {
			ou = append(ou, ouItem.(string))
		}
		model.Ou = ou
	}
	if modelMap["organization"] != nil {
		organization := []string{}
		for _, organizationItem := range modelMap["organization"].([]interface{}) {
			organization = append(organization, organizationItem.(string))
		}
		model.Organization = organization
	}
	if modelMap["country"] != nil {
		country := []string{}
		for _, countryItem := range modelMap["country"].([]interface{}) {
			country = append(country, countryItem.(string))
		}
		model.Country = country
	}
	if modelMap["locality"] != nil {
		locality := []string{}
		for _, localityItem := range modelMap["locality"].([]interface{}) {
			locality = append(locality, localityItem.(string))
		}
		model.Locality = locality
	}
	if modelMap["province"] != nil {
		province := []string{}
		for _, provinceItem := range modelMap["province"].([]interface{}) {
			province = append(province, provinceItem.(string))
		}
		model.Province = province
	}
	if modelMap["street_address"] != nil {
		streetAddress := []string{}
		for _, streetAddressItem := range modelMap["street_address"].([]interface{}) {
			streetAddress = append(streetAddress, streetAddressItem.(string))
		}
		model.StreetAddress = streetAddress
	}
	if modelMap["postal_code"] != nil {
		postalCode := []string{}
		for _, postalCodeItem := range modelMap["postal_code"].([]interface{}) {
			postalCode = append(postalCode, postalCodeItem.(string))
		}
		model.PostalCode = postalCode
	}
	if modelMap["serial_number"] != nil && modelMap["serial_number"].(string) != "" {
		model.SerialNumber = core.StringPtr(modelMap["serial_number"].(string))
	}
	if modelMap["data"] != nil {
	
	}
	if modelMap["signing_method"] != nil && modelMap["signing_method"].(string) != "" {
		model.SigningMethod = core.StringPtr(modelMap["signing_method"].(string))
	}
	if modelMap["issuer"] != nil && modelMap["issuer"].(string) != "" {
		model.Issuer = core.StringPtr(modelMap["issuer"].(string))
	}
	if modelMap["certificate_authority"] != nil && modelMap["certificate_authority"].(string) != "" {
		model.CertificateAuthority = core.StringPtr(modelMap["certificate_authority"].(string))
	}
	if modelMap["allowed_secret_groups"] != nil && modelMap["allowed_secret_groups"].(string) != "" {
		model.AllowedSecretGroups = core.StringPtr(modelMap["allowed_secret_groups"].(string))
	}
	if modelMap["allow_localhost"] != nil {
		model.AllowLocalhost = core.BoolPtr(modelMap["allow_localhost"].(bool))
	}
	if modelMap["allowed_domains"] != nil {
		allowedDomains := []string{}
		for _, allowedDomainsItem := range modelMap["allowed_domains"].([]interface{}) {
			allowedDomains = append(allowedDomains, allowedDomainsItem.(string))
		}
		model.AllowedDomains = allowedDomains
	}
	if modelMap["allowed_domains_template"] != nil {
		model.AllowedDomainsTemplate = core.BoolPtr(modelMap["allowed_domains_template"].(bool))
	}
	if modelMap["allow_bare_domains"] != nil {
		model.AllowBareDomains = core.BoolPtr(modelMap["allow_bare_domains"].(bool))
	}
	if modelMap["allow_subdomains"] != nil {
		model.AllowSubdomains = core.BoolPtr(modelMap["allow_subdomains"].(bool))
	}
	if modelMap["allow_glob_domains"] != nil {
		model.AllowGlobDomains = core.BoolPtr(modelMap["allow_glob_domains"].(bool))
	}
	if modelMap["allow_any_name"] != nil {
		model.AllowAnyName = core.BoolPtr(modelMap["allow_any_name"].(bool))
	}
	if modelMap["enforce_hostnames"] != nil {
		model.EnforceHostnames = core.BoolPtr(modelMap["enforce_hostnames"].(bool))
	}
	if modelMap["allow_ip_sans"] != nil {
		model.AllowIPSans = core.BoolPtr(modelMap["allow_ip_sans"].(bool))
	}
	if modelMap["allowed_uri_sans"] != nil {
		allowedURISans := []string{}
		for _, allowedURISansItem := range modelMap["allowed_uri_sans"].([]interface{}) {
			allowedURISans = append(allowedURISans, allowedURISansItem.(string))
		}
		model.AllowedURISans = allowedURISans
	}
	if modelMap["allowed_other_sans"] != nil {
		allowedOtherSans := []string{}
		for _, allowedOtherSansItem := range modelMap["allowed_other_sans"].([]interface{}) {
			allowedOtherSans = append(allowedOtherSans, allowedOtherSansItem.(string))
		}
		model.AllowedOtherSans = allowedOtherSans
	}
	if modelMap["server_flag"] != nil {
		model.ServerFlag = core.BoolPtr(modelMap["server_flag"].(bool))
	}
	if modelMap["client_flag"] != nil {
		model.ClientFlag = core.BoolPtr(modelMap["client_flag"].(bool))
	}
	if modelMap["code_signing_flag"] != nil {
		model.CodeSigningFlag = core.BoolPtr(modelMap["code_signing_flag"].(bool))
	}
	if modelMap["email_protection_flag"] != nil {
		model.EmailProtectionFlag = core.BoolPtr(modelMap["email_protection_flag"].(bool))
	}
	if modelMap["key_usage"] != nil {
		keyUsage := []string{}
		for _, keyUsageItem := range modelMap["key_usage"].([]interface{}) {
			keyUsage = append(keyUsage, keyUsageItem.(string))
		}
		model.KeyUsage = keyUsage
	}
	if modelMap["ext_key_usage"] != nil {
		extKeyUsage := []string{}
		for _, extKeyUsageItem := range modelMap["ext_key_usage"].([]interface{}) {
			extKeyUsage = append(extKeyUsage, extKeyUsageItem.(string))
		}
		model.ExtKeyUsage = extKeyUsage
	}
	if modelMap["ext_key_usage_oids"] != nil {
		extKeyUsageOids := []string{}
		for _, extKeyUsageOidsItem := range modelMap["ext_key_usage_oids"].([]interface{}) {
			extKeyUsageOids = append(extKeyUsageOids, extKeyUsageOidsItem.(string))
		}
		model.ExtKeyUsageOids = extKeyUsageOids
	}
	if modelMap["use_csr_common_name"] != nil {
		model.UseCsrCommonName = core.BoolPtr(modelMap["use_csr_common_name"].(bool))
	}
	if modelMap["use_csr_sans"] != nil {
		model.UseCsrSans = core.BoolPtr(modelMap["use_csr_sans"].(bool))
	}
	if modelMap["require_cn"] != nil {
		model.RequireCn = core.BoolPtr(modelMap["require_cn"].(bool))
	}
	if modelMap["policy_identifiers"] != nil {
		policyIdentifiers := []string{}
		for _, policyIdentifiersItem := range modelMap["policy_identifiers"].([]interface{}) {
			policyIdentifiers = append(policyIdentifiers, policyIdentifiersItem.(string))
		}
		model.PolicyIdentifiers = policyIdentifiers
	}
	if modelMap["basic_constraints_valid_for_non_ca"] != nil {
		model.BasicConstraintsValidForNonCa = core.BoolPtr(modelMap["basic_constraints_valid_for_non_ca"].(bool))
	}
	if modelMap["not_before_duration"] != nil {
	
	}
	return model, nil
}

func ResourceIBMSmCertConfigurationMapToConfigElementDefConfigLetsEncryptConfig(modelMap map[string]interface{}) (*secretsmanagerv1.ConfigElementDefConfigLetsEncryptConfig, error) {
	model := &secretsmanagerv1.ConfigElementDefConfigLetsEncryptConfig{}
	model.PrivateKey = core.StringPtr(modelMap["private_key"].(string))
	return model, nil
}

func ResourceIBMSmCertConfigurationMapToConfigElementDefConfigCloudInternetServicesConfig(modelMap map[string]interface{}) (*secretsmanagerv1.ConfigElementDefConfigCloudInternetServicesConfig, error) {
	model := &secretsmanagerv1.ConfigElementDefConfigCloudInternetServicesConfig{}
	model.CisCRN = core.StringPtr(modelMap["cis_crn"].(string))
	if modelMap["cis_apikey"] != nil && modelMap["cis_apikey"].(string) != "" {
		model.CisApikey = core.StringPtr(modelMap["cis_apikey"].(string))
	}
	return model, nil
}

func ResourceIBMSmCertConfigurationMapToConfigElementDefConfigClassicInfrastructureConfig(modelMap map[string]interface{}) (*secretsmanagerv1.ConfigElementDefConfigClassicInfrastructureConfig, error) {
	model := &secretsmanagerv1.ConfigElementDefConfigClassicInfrastructureConfig{}
	model.ClassicInfrastructureUsername = core.StringPtr(modelMap["classic_infrastructure_username"].(string))
	model.ClassicInfrastructurePassword = core.StringPtr(modelMap["classic_infrastructure_password"].(string))
	return model, nil
}

func ResourceIBMSmCertConfigurationMapToRootCertificateAuthorityConfig(modelMap map[string]interface{}) (*secretsmanagerv1.RootCertificateAuthorityConfig, error) {
	model := &secretsmanagerv1.RootCertificateAuthorityConfig{}

	if modelMap["crl_expiry"] != nil {
	
	}
	if modelMap["crl_disable"] != nil {
		model.CrlDisable = core.BoolPtr(modelMap["crl_disable"].(bool))
	}
	if modelMap["crl_distribution_points_encoded"] != nil {
		model.CrlDistributionPointsEncoded = core.BoolPtr(modelMap["crl_distribution_points_encoded"].(bool))
	}
	if modelMap["issuing_certificates_urls_encoded"] != nil {
		model.IssuingCertificatesUrlsEncoded = core.BoolPtr(modelMap["issuing_certificates_urls_encoded"].(bool))
	}
	model.CommonName = core.StringPtr(modelMap["common_name"].(string))
	if modelMap["status"] != nil && modelMap["status"].(string) != "" {
		model.Status = core.StringPtr(modelMap["status"].(string))
	}
	if modelMap["expiration_date"] != nil {
	
	}
	if modelMap["alt_names"] != nil && modelMap["alt_names"].(string) != "" {
		model.AltNames = core.StringPtr(modelMap["alt_names"].(string))
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
	if modelMap["key_type"] != nil && modelMap["key_type"].(string) != "" {
		model.KeyType = core.StringPtr(modelMap["key_type"].(string))
	}
	if modelMap["key_bits"] != nil {
		model.KeyBits = core.Int64Ptr(int64(modelMap["key_bits"].(int)))
	}
	if modelMap["max_path_length"] != nil {
		model.MaxPathLength = core.Int64Ptr(int64(modelMap["max_path_length"].(int)))
	}
	if modelMap["exclude_cn_from_sans"] != nil {
		model.ExcludeCnFromSans = core.BoolPtr(modelMap["exclude_cn_from_sans"].(bool))
	}
	if modelMap["permitted_dns_domains"] != nil {
		permittedDNSDomains := []string{}
		for _, permittedDNSDomainsItem := range modelMap["permitted_dns_domains"].([]interface{}) {
			permittedDNSDomains = append(permittedDNSDomains, permittedDNSDomainsItem.(string))
		}
		model.PermittedDNSDomains = permittedDNSDomains
	}
	if modelMap["ou"] != nil {
		ou := []string{}
		for _, ouItem := range modelMap["ou"].([]interface{}) {
			ou = append(ou, ouItem.(string))
		}
		model.Ou = ou
	}
	if modelMap["organization"] != nil {
		organization := []string{}
		for _, organizationItem := range modelMap["organization"].([]interface{}) {
			organization = append(organization, organizationItem.(string))
		}
		model.Organization = organization
	}
	if modelMap["country"] != nil {
		country := []string{}
		for _, countryItem := range modelMap["country"].([]interface{}) {
			country = append(country, countryItem.(string))
		}
		model.Country = country
	}
	if modelMap["locality"] != nil {
		locality := []string{}
		for _, localityItem := range modelMap["locality"].([]interface{}) {
			locality = append(locality, localityItem.(string))
		}
		model.Locality = locality
	}
	if modelMap["province"] != nil {
		province := []string{}
		for _, provinceItem := range modelMap["province"].([]interface{}) {
			province = append(province, provinceItem.(string))
		}
		model.Province = province
	}
	if modelMap["street_address"] != nil {
		streetAddress := []string{}
		for _, streetAddressItem := range modelMap["street_address"].([]interface{}) {
			streetAddress = append(streetAddress, streetAddressItem.(string))
		}
		model.StreetAddress = streetAddress
	}
	if modelMap["postal_code"] != nil {
		postalCode := []string{}
		for _, postalCodeItem := range modelMap["postal_code"].([]interface{}) {
			postalCode = append(postalCode, postalCodeItem.(string))
		}
		model.PostalCode = postalCode
	}
	if modelMap["serial_number"] != nil && modelMap["serial_number"].(string) != "" {
		model.SerialNumber = core.StringPtr(modelMap["serial_number"].(string))
	}
	if modelMap["data"] != nil {
	
	}
	return model, nil
}

func ResourceIBMSmCertConfigurationMapToIntermediateCertificateAuthorityConfig(modelMap map[string]interface{}) (*secretsmanagerv1.IntermediateCertificateAuthorityConfig, error) {
	model := &secretsmanagerv1.IntermediateCertificateAuthorityConfig{}

	model.SigningMethod = core.StringPtr(modelMap["signing_method"].(string))
	if modelMap["issuer"] != nil && modelMap["issuer"].(string) != "" {
		model.Issuer = core.StringPtr(modelMap["issuer"].(string))
	}
	if modelMap["crl_expiry"] != nil {
	
	}
	if modelMap["crl_disable"] != nil {
		model.CrlDisable = core.BoolPtr(modelMap["crl_disable"].(bool))
	}
	if modelMap["crl_distribution_points_encoded"] != nil {
		model.CrlDistributionPointsEncoded = core.BoolPtr(modelMap["crl_distribution_points_encoded"].(bool))
	}
	if modelMap["issuing_certificates_urls_encoded"] != nil {
		model.IssuingCertificatesUrlsEncoded = core.BoolPtr(modelMap["issuing_certificates_urls_encoded"].(bool))
	}
	model.CommonName = core.StringPtr(modelMap["common_name"].(string))
	if modelMap["status"] != nil && modelMap["status"].(string) != "" {
		model.Status = core.StringPtr(modelMap["status"].(string))
	}
	if modelMap["expiration_date"] != nil {
	
	}
	if modelMap["alt_names"] != nil && modelMap["alt_names"].(string) != "" {
		model.AltNames = core.StringPtr(modelMap["alt_names"].(string))
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
	if modelMap["key_type"] != nil && modelMap["key_type"].(string) != "" {
		model.KeyType = core.StringPtr(modelMap["key_type"].(string))
	}
	if modelMap["key_bits"] != nil {
		model.KeyBits = core.Int64Ptr(int64(modelMap["key_bits"].(int)))
	}
	if modelMap["exclude_cn_from_sans"] != nil {
		model.ExcludeCnFromSans = core.BoolPtr(modelMap["exclude_cn_from_sans"].(bool))
	}
	if modelMap["ou"] != nil {
		ou := []string{}
		for _, ouItem := range modelMap["ou"].([]interface{}) {
			ou = append(ou, ouItem.(string))
		}
		model.Ou = ou
	}
	if modelMap["organization"] != nil {
		organization := []string{}
		for _, organizationItem := range modelMap["organization"].([]interface{}) {
			organization = append(organization, organizationItem.(string))
		}
		model.Organization = organization
	}
	if modelMap["country"] != nil {
		country := []string{}
		for _, countryItem := range modelMap["country"].([]interface{}) {
			country = append(country, countryItem.(string))
		}
		model.Country = country
	}
	if modelMap["locality"] != nil {
		locality := []string{}
		for _, localityItem := range modelMap["locality"].([]interface{}) {
			locality = append(locality, localityItem.(string))
		}
		model.Locality = locality
	}
	if modelMap["province"] != nil {
		province := []string{}
		for _, provinceItem := range modelMap["province"].([]interface{}) {
			province = append(province, provinceItem.(string))
		}
		model.Province = province
	}
	if modelMap["street_address"] != nil {
		streetAddress := []string{}
		for _, streetAddressItem := range modelMap["street_address"].([]interface{}) {
			streetAddress = append(streetAddress, streetAddressItem.(string))
		}
		model.StreetAddress = streetAddress
	}
	if modelMap["postal_code"] != nil {
		postalCode := []string{}
		for _, postalCodeItem := range modelMap["postal_code"].([]interface{}) {
			postalCode = append(postalCode, postalCodeItem.(string))
		}
		model.PostalCode = postalCode
	}
	if modelMap["serial_number"] != nil && modelMap["serial_number"].(string) != "" {
		model.SerialNumber = core.StringPtr(modelMap["serial_number"].(string))
	}
	if modelMap["data"] != nil {
	
	}
	return model, nil
}

func ResourceIBMSmCertConfigurationMapToCertificateTemplateConfig(modelMap map[string]interface{}) (*secretsmanagerv1.CertificateTemplateConfig, error) {
	model := &secretsmanagerv1.CertificateTemplateConfig{}
	model.CertificateAuthority = core.StringPtr(modelMap["certificate_authority"].(string))
	if modelMap["allowed_secret_groups"] != nil && modelMap["allowed_secret_groups"].(string) != "" {
		model.AllowedSecretGroups = core.StringPtr(modelMap["allowed_secret_groups"].(string))
	}
	if modelMap["max_ttl"] != nil {
	
	}
	if modelMap["ttl"] != nil {
	
	}
	if modelMap["allow_localhost"] != nil {
		model.AllowLocalhost = core.BoolPtr(modelMap["allow_localhost"].(bool))
	}
	if modelMap["allowed_domains"] != nil {
		allowedDomains := []string{}
		for _, allowedDomainsItem := range modelMap["allowed_domains"].([]interface{}) {
			allowedDomains = append(allowedDomains, allowedDomainsItem.(string))
		}
		model.AllowedDomains = allowedDomains
	}
	if modelMap["allowed_domains_template"] != nil {
		model.AllowedDomainsTemplate = core.BoolPtr(modelMap["allowed_domains_template"].(bool))
	}
	if modelMap["allow_bare_domains"] != nil {
		model.AllowBareDomains = core.BoolPtr(modelMap["allow_bare_domains"].(bool))
	}
	if modelMap["allow_subdomains"] != nil {
		model.AllowSubdomains = core.BoolPtr(modelMap["allow_subdomains"].(bool))
	}
	if modelMap["allow_glob_domains"] != nil {
		model.AllowGlobDomains = core.BoolPtr(modelMap["allow_glob_domains"].(bool))
	}
	if modelMap["allow_any_name"] != nil {
		model.AllowAnyName = core.BoolPtr(modelMap["allow_any_name"].(bool))
	}
	if modelMap["enforce_hostnames"] != nil {
		model.EnforceHostnames = core.BoolPtr(modelMap["enforce_hostnames"].(bool))
	}
	if modelMap["allow_ip_sans"] != nil {
		model.AllowIPSans = core.BoolPtr(modelMap["allow_ip_sans"].(bool))
	}
	if modelMap["allowed_uri_sans"] != nil {
		allowedURISans := []string{}
		for _, allowedURISansItem := range modelMap["allowed_uri_sans"].([]interface{}) {
			allowedURISans = append(allowedURISans, allowedURISansItem.(string))
		}
		model.AllowedURISans = allowedURISans
	}
	if modelMap["allowed_other_sans"] != nil {
		allowedOtherSans := []string{}
		for _, allowedOtherSansItem := range modelMap["allowed_other_sans"].([]interface{}) {
			allowedOtherSans = append(allowedOtherSans, allowedOtherSansItem.(string))
		}
		model.AllowedOtherSans = allowedOtherSans
	}
	if modelMap["server_flag"] != nil {
		model.ServerFlag = core.BoolPtr(modelMap["server_flag"].(bool))
	}
	if modelMap["client_flag"] != nil {
		model.ClientFlag = core.BoolPtr(modelMap["client_flag"].(bool))
	}
	if modelMap["code_signing_flag"] != nil {
		model.CodeSigningFlag = core.BoolPtr(modelMap["code_signing_flag"].(bool))
	}
	if modelMap["email_protection_flag"] != nil {
		model.EmailProtectionFlag = core.BoolPtr(modelMap["email_protection_flag"].(bool))
	}
	if modelMap["key_type"] != nil && modelMap["key_type"].(string) != "" {
		model.KeyType = core.StringPtr(modelMap["key_type"].(string))
	}
	if modelMap["key_bits"] != nil {
		model.KeyBits = core.Int64Ptr(int64(modelMap["key_bits"].(int)))
	}
	if modelMap["key_usage"] != nil {
		keyUsage := []string{}
		for _, keyUsageItem := range modelMap["key_usage"].([]interface{}) {
			keyUsage = append(keyUsage, keyUsageItem.(string))
		}
		model.KeyUsage = keyUsage
	}
	if modelMap["ext_key_usage"] != nil {
		extKeyUsage := []string{}
		for _, extKeyUsageItem := range modelMap["ext_key_usage"].([]interface{}) {
			extKeyUsage = append(extKeyUsage, extKeyUsageItem.(string))
		}
		model.ExtKeyUsage = extKeyUsage
	}
	if modelMap["ext_key_usage_oids"] != nil {
		extKeyUsageOids := []string{}
		for _, extKeyUsageOidsItem := range modelMap["ext_key_usage_oids"].([]interface{}) {
			extKeyUsageOids = append(extKeyUsageOids, extKeyUsageOidsItem.(string))
		}
		model.ExtKeyUsageOids = extKeyUsageOids
	}
	if modelMap["use_csr_common_name"] != nil {
		model.UseCsrCommonName = core.BoolPtr(modelMap["use_csr_common_name"].(bool))
	}
	if modelMap["use_csr_sans"] != nil {
		model.UseCsrSans = core.BoolPtr(modelMap["use_csr_sans"].(bool))
	}
	if modelMap["ou"] != nil {
		ou := []string{}
		for _, ouItem := range modelMap["ou"].([]interface{}) {
			ou = append(ou, ouItem.(string))
		}
		model.Ou = ou
	}
	if modelMap["organization"] != nil {
		organization := []string{}
		for _, organizationItem := range modelMap["organization"].([]interface{}) {
			organization = append(organization, organizationItem.(string))
		}
		model.Organization = organization
	}
	if modelMap["country"] != nil {
		country := []string{}
		for _, countryItem := range modelMap["country"].([]interface{}) {
			country = append(country, countryItem.(string))
		}
		model.Country = country
	}
	if modelMap["locality"] != nil {
		locality := []string{}
		for _, localityItem := range modelMap["locality"].([]interface{}) {
			locality = append(locality, localityItem.(string))
		}
		model.Locality = locality
	}
	if modelMap["province"] != nil {
		province := []string{}
		for _, provinceItem := range modelMap["province"].([]interface{}) {
			province = append(province, provinceItem.(string))
		}
		model.Province = province
	}
	if modelMap["street_address"] != nil {
		streetAddress := []string{}
		for _, streetAddressItem := range modelMap["street_address"].([]interface{}) {
			streetAddress = append(streetAddress, streetAddressItem.(string))
		}
		model.StreetAddress = streetAddress
	}
	if modelMap["postal_code"] != nil {
		postalCode := []string{}
		for _, postalCodeItem := range modelMap["postal_code"].([]interface{}) {
			postalCode = append(postalCode, postalCodeItem.(string))
		}
		model.PostalCode = postalCode
	}
	if modelMap["serial_number"] != nil && modelMap["serial_number"].(string) != "" {
		model.SerialNumber = core.StringPtr(modelMap["serial_number"].(string))
	}
	if modelMap["require_cn"] != nil {
		model.RequireCn = core.BoolPtr(modelMap["require_cn"].(bool))
	}
	if modelMap["policy_identifiers"] != nil {
		policyIdentifiers := []string{}
		for _, policyIdentifiersItem := range modelMap["policy_identifiers"].([]interface{}) {
			policyIdentifiers = append(policyIdentifiers, policyIdentifiersItem.(string))
		}
		model.PolicyIdentifiers = policyIdentifiers
	}
	if modelMap["basic_constraints_valid_for_non_ca"] != nil {
		model.BasicConstraintsValidForNonCa = core.BoolPtr(modelMap["basic_constraints_valid_for_non_ca"].(bool))
	}
	if modelMap["not_before_duration"] != nil {
	
	}
	return model, nil
}

func ResourceIBMSmCertConfigurationConfigElementDefConfigToMap(model secretsmanagerv1.ConfigElementDefConfigIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv1.ConfigElementDefConfigLetsEncryptConfig); ok {
		return ResourceIBMSmCertConfigurationConfigElementDefConfigLetsEncryptConfigToMap(model.(*secretsmanagerv1.ConfigElementDefConfigLetsEncryptConfig))
	} else if _, ok := model.(*secretsmanagerv1.ConfigElementDefConfigCloudInternetServicesConfig); ok {
		return ResourceIBMSmCertConfigurationConfigElementDefConfigCloudInternetServicesConfigToMap(model.(*secretsmanagerv1.ConfigElementDefConfigCloudInternetServicesConfig))
	} else if _, ok := model.(*secretsmanagerv1.ConfigElementDefConfigClassicInfrastructureConfig); ok {
		return ResourceIBMSmCertConfigurationConfigElementDefConfigClassicInfrastructureConfigToMap(model.(*secretsmanagerv1.ConfigElementDefConfigClassicInfrastructureConfig))
	} else if _, ok := model.(*secretsmanagerv1.RootCertificateAuthorityConfig); ok {
		return ResourceIBMSmCertConfigurationRootCertificateAuthorityConfigToMap(model.(*secretsmanagerv1.RootCertificateAuthorityConfig))
	} else if _, ok := model.(*secretsmanagerv1.IntermediateCertificateAuthorityConfig); ok {
		return ResourceIBMSmCertConfigurationIntermediateCertificateAuthorityConfigToMap(model.(*secretsmanagerv1.IntermediateCertificateAuthorityConfig))
	} else if _, ok := model.(*secretsmanagerv1.CertificateTemplateConfig); ok {
		return ResourceIBMSmCertConfigurationCertificateTemplateConfigToMap(model.(*secretsmanagerv1.CertificateTemplateConfig))
	} else if _, ok := model.(*secretsmanagerv1.ConfigElementDefConfig); ok {
		modelMap := make(map[string]interface{})
		model := model.(*secretsmanagerv1.ConfigElementDefConfig)
		if model.PrivateKey != nil {
			modelMap["private_key"] = model.PrivateKey
		}
		if model.CisCRN != nil {
			modelMap["cis_crn"] = model.CisCRN
		}
		if model.CisApikey != nil {
			modelMap["cis_apikey"] = model.CisApikey
		}
		if model.ClassicInfrastructureUsername != nil {
			modelMap["classic_infrastructure_username"] = model.ClassicInfrastructureUsername
		}
		if model.ClassicInfrastructurePassword != nil {
			modelMap["classic_infrastructure_password"] = model.ClassicInfrastructurePassword
		}
		if model.MaxTTL != nil {
			modelMap["max_ttl"] = model.MaxTTL
		}
		if model.CrlExpiry != nil {
			modelMap["crl_expiry"] = model.CrlExpiry
		}
		if model.CrlDisable != nil {
			modelMap["crl_disable"] = model.CrlDisable
		}
		if model.CrlDistributionPointsEncoded != nil {
			modelMap["crl_distribution_points_encoded"] = model.CrlDistributionPointsEncoded
		}
		if model.IssuingCertificatesUrlsEncoded != nil {
			modelMap["issuing_certificates_urls_encoded"] = model.IssuingCertificatesUrlsEncoded
		}
		if model.CommonName != nil {
			modelMap["common_name"] = model.CommonName
		}
		if model.Status != nil {
			modelMap["status"] = model.Status
		}
		if model.ExpirationDate != nil {
			modelMap["expiration_date"] = model.ExpirationDate.String()
		}
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
		if model.KeyType != nil {
			modelMap["key_type"] = model.KeyType
		}
		if model.KeyBits != nil {
			modelMap["key_bits"] = flex.IntValue(model.KeyBits)
		}
		if model.MaxPathLength != nil {
			modelMap["max_path_length"] = flex.IntValue(model.MaxPathLength)
		}
		if model.ExcludeCnFromSans != nil {
			modelMap["exclude_cn_from_sans"] = model.ExcludeCnFromSans
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
			modelMap["serial_number"] = model.SerialNumber
		}
		if model.Data != nil {
			modelMap["data"] = model.Data
		}
		if model.SigningMethod != nil {
			modelMap["signing_method"] = model.SigningMethod
		}
		if model.Issuer != nil {
			modelMap["issuer"] = model.Issuer
		}
		if model.CertificateAuthority != nil {
			modelMap["certificate_authority"] = model.CertificateAuthority
		}
		if model.AllowedSecretGroups != nil {
			modelMap["allowed_secret_groups"] = model.AllowedSecretGroups
		}
		if model.AllowLocalhost != nil {
			modelMap["allow_localhost"] = model.AllowLocalhost
		}
		if model.AllowedDomains != nil {
			modelMap["allowed_domains"] = model.AllowedDomains
		}
		if model.AllowedDomainsTemplate != nil {
			modelMap["allowed_domains_template"] = model.AllowedDomainsTemplate
		}
		if model.AllowBareDomains != nil {
			modelMap["allow_bare_domains"] = model.AllowBareDomains
		}
		if model.AllowSubdomains != nil {
			modelMap["allow_subdomains"] = model.AllowSubdomains
		}
		if model.AllowGlobDomains != nil {
			modelMap["allow_glob_domains"] = model.AllowGlobDomains
		}
		if model.AllowAnyName != nil {
			modelMap["allow_any_name"] = model.AllowAnyName
		}
		if model.EnforceHostnames != nil {
			modelMap["enforce_hostnames"] = model.EnforceHostnames
		}
		if model.AllowIPSans != nil {
			modelMap["allow_ip_sans"] = model.AllowIPSans
		}
		if model.AllowedURISans != nil {
			modelMap["allowed_uri_sans"] = model.AllowedURISans
		}
		if model.AllowedOtherSans != nil {
			modelMap["allowed_other_sans"] = model.AllowedOtherSans
		}
		if model.ServerFlag != nil {
			modelMap["server_flag"] = model.ServerFlag
		}
		if model.ClientFlag != nil {
			modelMap["client_flag"] = model.ClientFlag
		}
		if model.CodeSigningFlag != nil {
			modelMap["code_signing_flag"] = model.CodeSigningFlag
		}
		if model.EmailProtectionFlag != nil {
			modelMap["email_protection_flag"] = model.EmailProtectionFlag
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
			modelMap["use_csr_common_name"] = model.UseCsrCommonName
		}
		if model.UseCsrSans != nil {
			modelMap["use_csr_sans"] = model.UseCsrSans
		}
		if model.RequireCn != nil {
			modelMap["require_cn"] = model.RequireCn
		}
		if model.PolicyIdentifiers != nil {
			modelMap["policy_identifiers"] = model.PolicyIdentifiers
		}
		if model.BasicConstraintsValidForNonCa != nil {
			modelMap["basic_constraints_valid_for_non_ca"] = model.BasicConstraintsValidForNonCa
		}
		if model.NotBeforeDuration != nil {
			modelMap["not_before_duration"] = model.NotBeforeDuration
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv1.ConfigElementDefConfigIntf subtype encountered")
	}
}

func ResourceIBMSmCertConfigurationConfigElementDefConfigLetsEncryptConfigToMap(model *secretsmanagerv1.ConfigElementDefConfigLetsEncryptConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["private_key"] = model.PrivateKey
	return modelMap, nil
}

func ResourceIBMSmCertConfigurationConfigElementDefConfigCloudInternetServicesConfigToMap(model *secretsmanagerv1.ConfigElementDefConfigCloudInternetServicesConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["cis_crn"] = model.CisCRN
	if model.CisApikey != nil {
		modelMap["cis_apikey"] = model.CisApikey
	}
	return modelMap, nil
}

func ResourceIBMSmCertConfigurationConfigElementDefConfigClassicInfrastructureConfigToMap(model *secretsmanagerv1.ConfigElementDefConfigClassicInfrastructureConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["classic_infrastructure_username"] = model.ClassicInfrastructureUsername
	modelMap["classic_infrastructure_password"] = model.ClassicInfrastructurePassword
	return modelMap, nil
}

func ResourceIBMSmCertConfigurationRootCertificateAuthorityConfigToMap(model *secretsmanagerv1.RootCertificateAuthorityConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["max_ttl"] = model.MaxTTL
	if model.CrlExpiry != nil {
		modelMap["crl_expiry"] = model.CrlExpiry
	}
	if model.CrlDisable != nil {
		modelMap["crl_disable"] = model.CrlDisable
	}
	if model.CrlDistributionPointsEncoded != nil {
		modelMap["crl_distribution_points_encoded"] = model.CrlDistributionPointsEncoded
	}
	if model.IssuingCertificatesUrlsEncoded != nil {
		modelMap["issuing_certificates_urls_encoded"] = model.IssuingCertificatesUrlsEncoded
	}
	modelMap["common_name"] = model.CommonName
	if model.Status != nil {
		modelMap["status"] = model.Status
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
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
	if model.KeyType != nil {
		modelMap["key_type"] = model.KeyType
	}
	if model.KeyBits != nil {
		modelMap["key_bits"] = flex.IntValue(model.KeyBits)
	}
	if model.MaxPathLength != nil {
		modelMap["max_path_length"] = flex.IntValue(model.MaxPathLength)
	}
	if model.ExcludeCnFromSans != nil {
		modelMap["exclude_cn_from_sans"] = model.ExcludeCnFromSans
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
		modelMap["serial_number"] = model.SerialNumber
	}
	if model.Data != nil {
		modelMap["data"] = model.Data
	}
	return modelMap, nil
}

func ResourceIBMSmCertConfigurationIntermediateCertificateAuthorityConfigToMap(model *secretsmanagerv1.IntermediateCertificateAuthorityConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["max_ttl"] = model.MaxTTL
	modelMap["signing_method"] = model.SigningMethod
	if model.Issuer != nil {
		modelMap["issuer"] = model.Issuer
	}
	if model.CrlExpiry != nil {
		modelMap["crl_expiry"] = model.CrlExpiry
	}
	if model.CrlDisable != nil {
		modelMap["crl_disable"] = model.CrlDisable
	}
	if model.CrlDistributionPointsEncoded != nil {
		modelMap["crl_distribution_points_encoded"] = model.CrlDistributionPointsEncoded
	}
	if model.IssuingCertificatesUrlsEncoded != nil {
		modelMap["issuing_certificates_urls_encoded"] = model.IssuingCertificatesUrlsEncoded
	}
	modelMap["common_name"] = model.CommonName
	if model.Status != nil {
		modelMap["status"] = model.Status
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
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
	if model.Format != nil {
		modelMap["format"] = model.Format
	}
	if model.PrivateKeyFormat != nil {
		modelMap["private_key_format"] = model.PrivateKeyFormat
	}
	if model.KeyType != nil {
		modelMap["key_type"] = model.KeyType
	}
	if model.KeyBits != nil {
		modelMap["key_bits"] = flex.IntValue(model.KeyBits)
	}
	if model.ExcludeCnFromSans != nil {
		modelMap["exclude_cn_from_sans"] = model.ExcludeCnFromSans
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
		modelMap["serial_number"] = model.SerialNumber
	}
	if model.Data != nil {
		modelMap["data"] = model.Data
	}
	return modelMap, nil
}

func ResourceIBMSmCertConfigurationCertificateTemplateConfigToMap(model *secretsmanagerv1.CertificateTemplateConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["certificate_authority"] = model.CertificateAuthority
	if model.AllowedSecretGroups != nil {
		modelMap["allowed_secret_groups"] = model.AllowedSecretGroups
	}
	if model.MaxTTL != nil {
		modelMap["max_ttl"] = model.MaxTTL
	}
	if model.TTL != nil {
		modelMap["ttl"] = model.TTL
	}
	if model.AllowLocalhost != nil {
		modelMap["allow_localhost"] = model.AllowLocalhost
	}
	if model.AllowedDomains != nil {
		modelMap["allowed_domains"] = model.AllowedDomains
	}
	if model.AllowedDomainsTemplate != nil {
		modelMap["allowed_domains_template"] = model.AllowedDomainsTemplate
	}
	if model.AllowBareDomains != nil {
		modelMap["allow_bare_domains"] = model.AllowBareDomains
	}
	if model.AllowSubdomains != nil {
		modelMap["allow_subdomains"] = model.AllowSubdomains
	}
	if model.AllowGlobDomains != nil {
		modelMap["allow_glob_domains"] = model.AllowGlobDomains
	}
	if model.AllowAnyName != nil {
		modelMap["allow_any_name"] = model.AllowAnyName
	}
	if model.EnforceHostnames != nil {
		modelMap["enforce_hostnames"] = model.EnforceHostnames
	}
	if model.AllowIPSans != nil {
		modelMap["allow_ip_sans"] = model.AllowIPSans
	}
	if model.AllowedURISans != nil {
		modelMap["allowed_uri_sans"] = model.AllowedURISans
	}
	if model.AllowedOtherSans != nil {
		modelMap["allowed_other_sans"] = model.AllowedOtherSans
	}
	if model.ServerFlag != nil {
		modelMap["server_flag"] = model.ServerFlag
	}
	if model.ClientFlag != nil {
		modelMap["client_flag"] = model.ClientFlag
	}
	if model.CodeSigningFlag != nil {
		modelMap["code_signing_flag"] = model.CodeSigningFlag
	}
	if model.EmailProtectionFlag != nil {
		modelMap["email_protection_flag"] = model.EmailProtectionFlag
	}
	if model.KeyType != nil {
		modelMap["key_type"] = model.KeyType
	}
	if model.KeyBits != nil {
		modelMap["key_bits"] = flex.IntValue(model.KeyBits)
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
		modelMap["use_csr_common_name"] = model.UseCsrCommonName
	}
	if model.UseCsrSans != nil {
		modelMap["use_csr_sans"] = model.UseCsrSans
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
		modelMap["serial_number"] = model.SerialNumber
	}
	if model.RequireCn != nil {
		modelMap["require_cn"] = model.RequireCn
	}
	if model.PolicyIdentifiers != nil {
		modelMap["policy_identifiers"] = model.PolicyIdentifiers
	}
	if model.BasicConstraintsValidForNonCa != nil {
		modelMap["basic_constraints_valid_for_non_ca"] = model.BasicConstraintsValidForNonCa
	}
	if model.NotBeforeDuration != nil {
		modelMap["not_before_duration"] = model.NotBeforeDuration
	}
	return modelMap, nil
}
