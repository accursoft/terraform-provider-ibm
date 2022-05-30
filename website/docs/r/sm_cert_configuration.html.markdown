---
layout: "ibm"
page_title: "IBM : ibm_sm_cert_configuration"
description: |-
  Manages sm_cert_configuration.
subcategory: "Secrets Manager"
---

# ibm_sm_cert_configuration

Provides a resource for sm_cert_configuration. This allows sm_cert_configuration to be created, updated and deleted.

## Example Usage

```hcl
resource "ibm_sm_cert_configuration" "sm_cert_configuration" {
  config {
		private_key = "private_key"
  }
  config_element = "certificate_authorities"
  name = "cis-example-config"
  secret_type = "public_cert"
  type = "cis"
}
```

## Argument Reference

Review the argument reference that you can specify for your resource.

* `config` - (Required, List) The configuration to define for the specified secret type.
Nested scheme for **config**:
	* `allow_any_name` - (Optional, Boolean) Determines whether to allow clients to request a private certificate that matches any common name.
	  * Constraints: The default value is `false`.
	* `allow_bare_domains` - (Optional, Boolean) Determines whether to allow clients to request private certificates that match the value of the actual domains on the final certificate.For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a certificate that contains the name `example.com` as one of the DNS values on the final certificate.**Important:** In some scenarios, allowing bare domains can be considered a security risk.
	  * Constraints: The default value is `false`.
	* `allow_glob_domains` - (Optional, Boolean) Determines whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are specified in the `allowed_domains` field.If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.
	  * Constraints: The default value is `false`.
	* `allow_ip_sans` - (Optional, Boolean) Determines whether to allow clients to request a private certificate with IP Subject Alternative Names.
	  * Constraints: The default value is `true`.
	* `allow_localhost` - (Optional, Boolean) Determines whether to allow `localhost` to be included as one of the requested common names.
	  * Constraints: The default value is `true`.
	* `allow_subdomains` - (Optional, Boolean) Determines whether to allow clients to request private certificates with common names (CN) that are subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.**Note:** This field is redundant if you use the `allow_any_name` option.
	  * Constraints: The default value is `false`.
	* `allowed_domains` - (Optional, List) The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and `allow_subdomains` options.
	* `allowed_domains_template` - (Optional, Boolean) Determines whether to allow the domains that are supplied in the `allowed_domains` field to contain access control list (ACL) templates.
	  * Constraints: The default value is `false`.
	* `allowed_other_sans` - (Optional, List) The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private certificates. The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any `other_sans` input.
	* `allowed_secret_groups` - (Optional, String) Scopes the creation of private certificates to only the secret groups that you specify. This field can be supplied as a comma-delimited list of secret group IDs.
	* `allowed_uri_sans` - (Optional, List) The URI Subject Alternative Names to allow for private certificates.Values can contain glob patterns, for example `spiffe://hostname/_*`.
	* `alt_names` - (Optional, String) The Subject Alternative Names to define for the CA certificate, in a comma-delimited list.The alternative names can be host names or email addresses.
	  * Constraints: The maximum length is `2048` characters.
	* `basic_constraints_valid_for_non_ca` - (Optional, Boolean) Determines whether to mark the Basic Constraints extension of an issued private certificate as valid for non-CA certificates.
	* `certificate_authority` - (Optional, String) The name of the intermediate certificate authority.
	* `cis_apikey` - (Optional, String) An IBM Cloud API key that can to list domains in your CIS instance.To grant Secrets Manager the ability to view the CIS instance and all of its domains, the API key must be assigned the Reader service role on Internet Services (`internet-svcs`).If you need to manage specific domains, you can assign the Manager role. For production environments, it is recommended that you assign the Reader access role, and then use the[IAM Policy Management API](https://cloud.ibm.com/apidocs/iam-policy-management#create-policy) to control specific domains. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#authorize-specific-domains).
	* `cis_crn` - (Optional, String) The Cloud Resource Name (CRN) that is associated with the CIS instance.
	* `classic_infrastructure_password` - (Optional, String) Your classic infrastructure API key.For information about viewing and accessing your classic infrastructure API key, see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	* `classic_infrastructure_username` - (Optional, String) The username that is associated with your classic infrastructure account.In most cases, your classic infrastructure username is your `<account_id>_<email_address>`. For more information, see the [docs](https://cloud.ibm.com/docs/account?topic=account-classic_keys).
	* `client_flag` - (Optional, Boolean) Determines whether private certificates are flagged for client use.
	  * Constraints: The default value is `true`.
	* `code_signing_flag` - (Optional, Boolean) Determines whether private certificates are flagged for code signing use.
	  * Constraints: The default value is `false`.
	* `common_name` - (Optional, String) The fully qualified domain name or host domain name for the certificate.
	* `country` - (Optional, List) The Country (C) values to define in the subject field of the resulting certificate.
	* `crl_disable` - (Optional, Boolean) Disables or enables certificate revocation list (CRL) building. If CRL building is disabled, a signed but zero-length CRL is returned when downloading the CRL. If CRL building is enabled,  it will rebuild the CRL.
	  * Constraints: The default value is `false`.
	* `crl_distribution_points_encoded` - (Optional, Boolean) Determines whether to encode the certificate revocation list (CRL) distribution points in the certificates that are issued by this certificate authority.
	  * Constraints: The default value is `false`.
	* `crl_expiry` - (Optional, Map) The time until the certificate revocation list (CRL) expires.The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours. In the API response, this value is returned in seconds (integer).**Note:** The CRL is rotated automatically before it expires.
	* `data` - (Optional, Map) The data that is associated with the root certificate authority. The data object contains the following fields:- `certificate`: The root certificate content.- `issuing_ca`: The certificate of the certificate authority that signed and issued this certificate.- `serial_number`: The unique serial number of the root certificate.
	* `email_protection_flag` - (Optional, Boolean) Determines whether private certificates are flagged for email protection use.
	  * Constraints: The default value is `false`.
	* `enforce_hostnames` - (Optional, Boolean) Determines whether to enforce only valid host names for common names, DNS Subject Alternative Names, and the host section of email addresses.
	  * Constraints: The default value is `true`.
	* `exclude_cn_from_sans` - (Optional, Boolean) Controls whether the common name is excluded from Subject Alternative Names (SANs). If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if the common name is not a hostname or an email address, but is instead a human-readable identifier.
	  * Constraints: The default value is `false`.
	* `expiration_date` - (Optional, String) The date that the certificate expires. The date format follows RFC 3339.
	* `ext_key_usage` - (Optional, List) The allowed extended key usage constraint on private certificates.You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage). Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to an empty list.
	* `ext_key_usage_oids` - (Optional, List) A list of extended key usage Object Identifiers (OIDs).
	* `format` - (Optional, String) The format of the returned data.
	  * Constraints: The default value is `pem`. Allowable values are: `pem`, `pem_bundle`.
	* `ip_sans` - (Optional, String) The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	  * Constraints: The maximum length is `2048` characters.
	* `issuer` - (Optional, String) The certificate authority that signed and issued the certificate.If the certificate is signed internally, the `issuer` field is required and must match the name of a certificate authority that is configured in the Secrets Manager service instance.
	* `issuing_certificates_urls_encoded` - (Optional, Boolean) Determines whether to encode the URL of the issuing certificate in the certificates that are issued by this certificate authority.
	  * Constraints: The default value is `false`.
	* `key_bits` - (Optional, Integer) The number of bits to use when generating the private key.Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.
	* `key_type` - (Optional, String) The type of private key to generate.
	  * Constraints: Allowable values are: `rsa`, `ec`.
	* `key_usage` - (Optional, List) The allowed key usage constraint to define for private certificates.You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage).  Omit the `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to an empty list.
	* `locality` - (Optional, List) The Locality (L) values to define in the subject field of the resulting certificate.
	* `max_path_length` - (Optional, Integer) The maximum path length to encode in the generated certificate. `-1` means no limit.If the signing certificate has a maximum path length set, the path length is set to one less than that of the signing certificate. A limit of `0` means a literal path length of zero.
	* `max_ttl` - (Optional, Map) The maximum time-to-live (TTL) for certificates that are created by this CA.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).
	* `not_before_duration` - (Optional, Map) The duration in seconds by which to backdate the `not_before` property of an issued private certificate.The value can be supplied as a string representation of a duration, such as `30s`. In the API response, this value is returned in seconds (integer).
	* `organization` - (Optional, List) The Organization (O) values to define in the subject field of the resulting certificate.
	* `other_sans` - (Optional, List) The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`.
	* `ou` - (Optional, List) The Organizational Unit (OU) values to define in the subject field of the resulting certificate.
	* `permitted_dns_domains` - (Optional, List) The allowed DNS domains or subdomains for the certificates to be signed and issued by this CA certificate.
	* `policy_identifiers` - (Optional, List) A list of policy Object Identifiers (OIDs).
	* `postal_code` - (Optional, List) The Postal Code values in the subject field of the resulting certificate.
	* `private_key` - (Optional, String) The private key that is associated with your Automatic Certificate Management Environment (ACME) account.If you have a working ACME client or account for Let's Encrypt, you can use the existing private key to enable communications with Secrets Manager. If you don't have an account yet, you can create one. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-prepare-order-certificates#create-acme-account).
	* `private_key_format` - (Optional, String) The format of the generated private key.
	  * Constraints: The default value is `der`. Allowable values are: `der`, `pkcs8`.
	* `province` - (Optional, List) The Province (ST) values to define in the subject field of the resulting certificate.
	* `require_cn` - (Optional, Boolean) Determines whether to require a common name to create a private certificate.By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the `require_cn` option to `false`.
	  * Constraints: The default value is `true`.
	* `serial_number` - (Optional, String) The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.
	  * Constraints: The maximum length is `64` characters.
	* `server_flag` - (Optional, Boolean) Determines whether private certificates are flagged for server use.
	  * Constraints: The default value is `true`.
	* `signing_method` - (Optional, String) The signing method to use with this certificate authority to generate private certificates.You can choose between internal or externally signed options. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).
	  * Constraints: Allowable values are: `internal`, `external`.
	* `status` - (Optional, String) The status of the certificate authority. The status of a root certificate authority is either `configured` or `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,`signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.
	  * Constraints: Allowable values are: `signing_required`, `signed_certificate_required`, `certificate_template_required`, `configured`, `expired`, `revoked`.
	* `street_address` - (Optional, List) The Street Address values in the subject field of the resulting certificate.
	* `ttl` - (Optional, Map) The time-to-live (TTL) to assign to this CA certificate.The value can be supplied as a string representation of a duration, such as `12h`. The value can't exceed the `max_ttl` that is defined in the associated certificate template. In the API response, this value is returned in seconds (integer).
	* `uri_sans` - (Optional, String) The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	  * Constraints: The maximum length is `2048` characters.
	* `use_csr_common_name` - (Optional, Boolean) When used with the `sign_csr` action, this field determines whether to use the common name (CN) from a certificate signing request (CSR) instead of the CN that's included in the JSON data of the certificate.Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include the `use_csr_sans` property.
	  * Constraints: The default value is `true`.
	* `use_csr_sans` - (Optional, Boolean) When used with the `sign_csr` action, this field determines whether to use the Subject Alternative Names (SANs) from a certificate signing request (CSR) instead of the SANs that are included in the JSON data of the certificate.Does not include the common name in the CSR. To use the common name, include the `use_csr_common_name` property.
	  * Constraints: The default value is `true`.
* `config_element` - (Required, Forces new resource, String) The configuration element to define or manage. Allowable values are: certificate_authorities, dns_providers, root_certificate_authorities, intermediate_certificate_authorities, certificate_templates
  * Constraints: Allowable values are: `certificate_authorities`, `dns_providers`, `root_certificate_authorities`, `intermediate_certificate_authorities`, `certificate_templates`.
* `name` - (Required, String) The human-readable name to assign to your configuration.
  * Constraints: The maximum length is `256` characters. The minimum length is `2` characters.
* `secret_type` - (Required, Forces new resource, String) The secret type. Allowable values are: public_cert, private_cert
  * Constraints: Allowable values are: `public_cert`, `private_cert`.
* `type` - (Required, String) The type of configuration. Value options differ depending on the 'config_element'. Allowable values are: letsencrypt, letsencrypt-stage, cis, classic_infrastructure, root_certificate_authority, intermediate_certificate_authority
  * Constraints: Allowable values are: `letsencrypt`, `letsencrypt-stage`, `cis`, `classic_infrastructure`, `root_certificate_authority`, `intermediate_certificate_authority`, `certificate_template`. The maximum length is `128` characters. The minimum length is `2` characters.

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your resource is created.

* `id` - The unique identifier of the sm_cert_configuration.

## Provider Configuration

The IBM Cloud provider offers a flexible means of providing credentials for authentication. The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables

To find which credentials are required for this resource, see the service table [here](https://cloud.ibm.com/docs/ibm-cloud-provider-for-terraform?topic=ibm-cloud-provider-for-terraform-provider-reference#required-parameters).

### Static credentials

You can provide your static credentials by adding the `ibmcloud_api_key`, `iaas_classic_username`, and `iaas_classic_api_key` arguments in the IBM Cloud provider block.

Usage:
```
provider "ibm" {
    ibmcloud_api_key = ""
    iaas_classic_username = ""
    iaas_classic_api_key = ""
}
```

### Environment variables

You can provide your credentials by exporting the `IC_API_KEY`, `IAAS_CLASSIC_USERNAME`, and `IAAS_CLASSIC_API_KEY` environment variables, representing your IBM Cloud platform API key, IBM Cloud Classic Infrastructure (SoftLayer) user name, and IBM Cloud infrastructure API key, respectively.

```
provider "ibm" {}
```

Usage:
```
export IC_API_KEY="ibmcloud_api_key"
export IAAS_CLASSIC_USERNAME="iaas_classic_username"
export IAAS_CLASSIC_API_KEY="iaas_classic_api_key"
terraform plan
```

Note:

1. Create or find your `ibmcloud_api_key` and `iaas_classic_api_key` [here](https://cloud.ibm.com/iam/apikeys).
  - Select `My IBM Cloud API Keys` option from view dropdown for `ibmcloud_api_key`
  - Select `Classic Infrastructure API Keys` option from view dropdown for `iaas_classic_api_key`
2. For iaas_classic_username
  - Go to [Users](https://cloud.ibm.com/iam/users)
  - Click on user.
  - Find user name in the `VPN password` section under `User Details` tab

For more informaton, see [here](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs#authentication).

## Import

You can import the `ibm_sm_cert_configuration` resource by using `name`.
The `name` property can be formed from `secret_type`, `config_element`, and `config_name` in the following format:

```
<secret_type>/<config_element>/<config_name>
```
* `secret_type`: A string. The secret type.
* `config_element`: A string. The configuration element to define or manage.
* `config_name`: A string. The name of your configuration.
For more information, see [the documentation](https://cloud.ibm.com/docs/secrets-manager).

# Syntax
```
$ terraform import ibm_sm_cert_configuration.sm_cert_configuration <secret_type>/<config_element>/<config_name>
```
