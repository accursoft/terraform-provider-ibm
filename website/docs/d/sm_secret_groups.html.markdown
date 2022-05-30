---
layout: "ibm"
page_title: "IBM : ibm_sm_secret_groups"
description: |-
  Get information about sm_secret_groups
subcategory: "Secrets Manager"
---

# ibm_sm_secret_groups

Provides a read-only data source for sm_secret_groups. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.

## Example Usage

```hcl
data "ibm_sm_secret_groups" "sm_secret_groups" {
	id = "id"
}
```

## Argument Reference

Review the argument reference that you can specify for your data source.

* `id` - (Required, Forces new resource, String) The ID of the secret group.
  * Constraints: The value must match regular expression `/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/`.

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your data source is created.

* `id` - The unique identifier of the sm_secret_groups.
* `metadata` - (Required, List) The metadata that describes the resource array.
Nested scheme for **metadata**:
	* `collection_total` - (Required, Integer) The number of elements in the resource array.
	* `collection_type` - (Required, String) The type of resources in the resource array.
	  * Constraints: Allowable values are: `application/vnd.ibm.secrets-manager.config+json`, `application/vnd.ibm.secrets-manager.secret+json`, `application/vnd.ibm.secrets-manager.secret.version+json`, `application/vnd.ibm.secrets-manager.secret.policy+json`, `application/vnd.ibm.secrets-manager.secret.group+json`, `application/vnd.ibm.secrets-manager.error+json`.

* `resources` - (Required, List) The properties in JSON format to define, such as the name and description. For more information, see the docs: https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-cli-plugin-secrets-manager-cli#secrets-manager-cli-secret-group-create-command
Nested scheme for **resources**:
	* `creation_date` - (Optional, String) The date the secret group was created. The date format follows RFC 3339.
	* `description` - (Optional, String) An extended description of your secret group.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.
	  * Constraints: The maximum length is `1024` characters.
	* `id` - (Optional, String) The v4 UUID that uniquely identifies the secret group.
	* `last_update_date` - (Optional, String) Updates when the metadata of the secret group is modified. The date format follows RFC 3339.
	* `name` - (Optional, String) The type of policy. Allowable values are: rotation
	  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters.
	* `type` - (Optional, String) The MIME type that represents the secret group.

