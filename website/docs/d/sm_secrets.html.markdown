---
layout: "ibm"
page_title: "IBM : ibm_sm_secrets"
description: |-
  Get information about sm_secrets
subcategory: "Secrets Manager"
---

# ibm_sm_secrets

Provides a read-only data source for sm_secrets. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "ibm_sm_secrets" "sm_secrets" {
}
```


## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your data source is created.

* `id` - The unique identifier of the sm_secrets.
* `metadata` - (Required, List) The metadata that describes the resource array.
Nested scheme for **metadata**:
	* `collection_total` - (Required, Integer) The number of elements in the resource array.
	* `collection_type` - (Required, String) The type of resources in the resource array.
	  * Constraints: Allowable values are: `application/vnd.ibm.secrets-manager.config+json`, `application/vnd.ibm.secrets-manager.secret+json`, `application/vnd.ibm.secrets-manager.secret.version+json`, `application/vnd.ibm.secrets-manager.secret.policy+json`, `application/vnd.ibm.secrets-manager.secret.group+json`, `application/vnd.ibm.secrets-manager.error+json`.

* `resources` - (Optional, List) A collection of resources.
Nested scheme for **resources**:
	* `access_groups` - (Optional, List) The access groups that define the capabilities of the service ID and API key that are generated for an`iam_credentials` secret. If you prefer to use an existing service ID that is already assigned the access policies that you require, you can omit this parameter and use the `service_id` field instead.**Tip:** To list the access groups that are available in an account, you can use the [IAM Access Groups API](https://cloud.ibm.com/apidocs/iam-access-groups#list-access-groups). To find the ID of an access group in the console, go to **Manage > Access (IAM) > Access groups**. Select the access group to inspect, and click **Details** to view its ID.
	  * Constraints: The maximum length is `10` items. The minimum length is `1` item.
	* `algorithm` - (Optional, String) The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign the certificate.
	* `alt_names` - (Optional, Map) The alternative names that are defined for the certificate.For public certificates, this value is provided as an array of strings. For private certificates, this value is provided as a comma-delimited list (string). In the API response, this value is returned as an array of strings for all the types of certificate secrets.
	* `api_key` - (Optional, String) The API key that is generated for this secret.After the secret reaches the end of its lease (see the `ttl` field), the API key is deleted automatically. If you want to continue to use the same API key for future read operations, see the `reuse_api_key` field.
	* `api_key_id` - (Optional, String) The ID of the API key that is generated for this secret.
	* `bundle_certs` - (Optional, Boolean) Determines whether your issued certificate is bundled with intermediate certificates.Set to `false` for the certificate file to contain only the issued certificate.
	  * Constraints: The default value is `true`.
	* `ca` - (Optional, String) The name of the certificate authority configuration.
	* `certificate` - (Optional, String) The contents of your certificate.
	  * Constraints: The maximum length is `100000` characters. The minimum length is `50` characters.
	* `certificate_authority` - (Optional, String) The intermediate certificate authority that signed this certificate.
	* `certificate_template` - (Optional, String) The name of the certificate template.
	  * Constraints: The maximum length is `64` characters.
	* `common_name` - (Optional, String) The fully qualified domain name or host domain name that is defined for the certificate.
	* `created_by` - (Optional, String) The unique identifier for the entity that created the secret.
	* `creation_date` - (Optional, String) The date the secret was created. The date format follows RFC 3339.
	* `crn` - (Optional, String) The Cloud Resource Name (CRN) that uniquely identifies your Secrets Manager resource.
	* `description` - (Optional, String) An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret.
	  * Constraints: The maximum length is `1024` characters. The minimum length is `2` characters.
	* `dns` - (Optional, String) The name of the DNS provider configuration.
	* `exclude_cn_from_sans` - (Optional, Boolean) Controls whether the common name is excluded from Subject Alternative Names (SANs). If set to `true`, the common name is is not included in DNS or Email SANs if they apply. This field can be useful if the common name is not a hostname or an email address, but is instead a human-readable identifier.
	  * Constraints: The default value is `false`.
	* `expiration_date` - (Optional, String) The date the secret material expires. The date format follows RFC 3339.You can set an expiration date on supported secret types at their creation. If you create a secret without specifying an expiration date, the secret does not expire. The `expiration_date` field is supported for the following secret types:- `arbitrary`- `username_password`.
	* `format` - (Optional, String) The format of the returned data.
	  * Constraints: The default value is `pem`. Allowable values are: `pem`, `pem_bundle`.
	* `id` - (Optional, String) The v4 UUID that uniquely identifies the secret.
	* `intermediate` - (Optional, String) (Optional) The intermediate certificate to associate with the root certificate.
	  * Constraints: The maximum length is `100000` characters. The minimum length is `50` characters.
	* `intermediate_included` - (Optional, Boolean) Indicates whether the certificate was imported with an associated intermediate certificate.
	* `ip_sans` - (Optional, String) The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	  * Constraints: The maximum length is `2048` characters.
	* `issuance_info` - (Optional, List) Issuance information that is associated with your certificate.
	Nested scheme for **issuance_info**:
		* `auto_rotated` - (Optional, Boolean) Indicates whether the issued certificate is configured with an automatic rotation policy.
		* `bundle_certs` - (Optional, Boolean) Indicates whether the issued certificate is bundled with intermediate certificates.
		* `ca` - (Optional, String) The name that was assigned to the certificate authority configuration.
		* `dns` - (Optional, String) The name that was assigned to the DNS provider configuration.
		* `error_code` - (Optional, String) A code that identifies an issuance error.This field, along with `error_message`, is returned when Secrets Manager successfully processes your request, but a certificate is unable to be issued by the certificate authority.
		* `error_message` - (Optional, String) A human-readable message that provides details about the issuance error.
		* `ordered_on` - (Optional, String) The date the certificate was ordered. The date format follows RFC 3339.
		* `state` - (Optional, Integer) The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,  Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
		  * Constraints: Allowable values are: `0`, `1`, `2`, `3`, `5`.
		* `state_description` - (Optional, String) A text representation of the secret state.
	* `issuer` - (Optional, String) The distinguished name that identifies the entity that signed and issued the certificate.
	* `key_algorithm` - (Optional, String) The identifier for the cryptographic algorithm that was used to generate the public and private keys that are associated with the certificate.
	* `labels` - (Optional, List) Labels that you can use to filter for secrets in your instance.Up to 30 labels can be created. Labels can be 2 - 30 characters, including spaces. Special characters that are not permitted include the angled bracket, comma, colon, ampersand, and vertical pipe character (|).To protect your privacy, do not use personal data, such as your name or location, as a label for your secret.
	* `last_update_date` - (Optional, String) Updates when the actual secret is modified. The date format follows RFC 3339.
	* `name` - (Optional, String) A human-readable alias to assign to your secret.To protect your privacy, do not use personal data, such as your name or location, as an alias for your secret.
	  * Constraints: The maximum length is `256` characters. The minimum length is `2` characters. The value must match regular expression `/^\\w(([\\w-.]+)?\\w)?$/`.
	* `next_rotation_date` - (Optional, String) The date that the secret is scheduled for automatic rotation.The service automatically creates a new version of the secret on its next rotation date. This field exists only for secrets that can be auto-rotated and have an existing rotation policy.
	* `other_sans` - (Optional, List) The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`.
	* `password` - (Optional, String) The password to assign to this secret.
	  * Constraints: The maximum length is `64` characters.
	* `payload` - (Optional, String) The new secret data to assign to the secret.
	* `private_key` - (Optional, String) (Optional) The private key to associate with the certificate.
	  * Constraints: The maximum length is `100000` characters. The minimum length is `50` characters.
	* `private_key_format` - (Optional, String) The format of the generated private key.
	  * Constraints: The default value is `der`. Allowable values are: `der`, `pkcs8`.
	* `private_key_included` - (Optional, Boolean) Indicates whether the certificate was imported with an associated private key.
	* `reuse_api_key` - (Optional, Boolean) (IAM credentials) Reuse the service ID and API key for future read operations.
	  * Constraints: The default value is `false`.
	* `revocation_time` - (Optional, Integer) The timestamp of the certificate revocation.
	* `revocation_time_rfc3339` - (Optional, String) The date and time that the certificate was revoked. The date format follows RFC 3339.
	* `rotation` - (Optional, List)
	Nested scheme for **rotation**:
		* `auto_rotate` - (Optional, Boolean) Determines whether Secrets Manager rotates your certificate automatically.For public certificates, if `auto_rotate` is set to `true` the service reorders your certificate 31 days before it expires. For private certificates, the certificate is rotated according to the time interval specified in the `interval` and `unit` fields.To access the previous version of the certificate, you can use the[Get a version of a secret](#get-secret-version) method.
		  * Constraints: The default value is `false`.
		* `interval` - (Optional, Integer) Used together with the `unit` field to specify the rotation interval. The minimum interval is one day, and the maximum interval is 3 years (1095 days). Required in case `auto_rotate` is set to `true`.**Note:** Use this field only for private certificates. It is ignored for public certificates.
		* `rotate_keys` - (Optional, Boolean) Determines whether Secrets Manager rotates the private key for your certificate automatically.If set to `true`, the service generates and stores a new private key for your rotated certificate.**Note:** Use this field only for public certificates. It is ignored for private certificates.
		  * Constraints: The default value is `false`.
		* `unit` - (Optional, String) The time unit of the rotation interval.**Note:** Use this field only for private certificates. It is ignored for public certificates.
		  * Constraints: Allowable values are: `day`, `month`.
	* `secret_data` - (Optional, Map) The data that is associated with the secret version.The data object contains the field `payload`.
	* `secret_group_id` - (Optional, String) The v4 UUID that uniquely identifies the secret group to assign to this secret.If you omit this parameter, your secret is assigned to the `default` secret group.
	  * Constraints: The value must match regular expression `/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/`.
	* `secret_type` - (Optional, String) The secret type.
	  * Constraints: Allowable values are: `arbitrary`, `username_password`, `iam_credentials`, `imported_cert`, `public_cert`, `private_cert`, `kv`.
	* `serial_number` - (Optional, String) The unique serial number that was assigned to the certificate by the issuing certificate authority.
	* `service_id` - (Optional, String) (IAM credentials) The service ID under which the API key is created. To have Secrets Manager generate a new service ID, omit this option and include 'access_groups'.
	* `service_id_is_static` - (Optional, Boolean) Indicates whether an `iam_credentials` secret was created with a static service ID.If `true`, the service ID for the secret was provided by the user at secret creation. If `false`, the service ID was generated by Secrets Manager.
	* `state` - (Optional, Integer) The secret state based on NIST SP 800-57. States are integers and correspond to the Pre-activation = 0, Active = 1,  Suspended = 2, Deactivated = 3, and Destroyed = 5 values.
	  * Constraints: Allowable values are: `0`, `1`, `2`, `3`, `5`.
	* `state_description` - (Optional, String) A text representation of the secret state.
	* `ttl` - (Optional, Map) The time-to-live (TTL) or lease duration to assign to generated credentials.For `iam_credentials` secrets, the TTL defines for how long each generated API key remains valid. The value can be either an integer that specifies the number of seconds, or the string representation of a duration, such as `120m` or `24h`.Minimum duration is 1 minute. Maximum is 90 days.
	* `uri_sans` - (Optional, String) The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.
	  * Constraints: The maximum length is `2048` characters.
	* `username` - (Optional, String) The username to assign to this secret.
	  * Constraints: The maximum length is `64` characters.
	* `validity` - (Optional, List)
	Nested scheme for **validity**:
		* `not_after` - (Optional, String) The date and time that the certificate validity period ends.
		* `not_before` - (Optional, String) The date and time that the certificate validity period begins.
	* `versions` - (Optional, List) An array that contains metadata for each secret version. For more information on the metadata properties, see [Get secret version metadata](#get-secret-version-metadata).
	* `versions_total` - (Optional, Integer) The number of versions that are associated with a secret.

