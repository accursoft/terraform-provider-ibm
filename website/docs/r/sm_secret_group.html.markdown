---
layout: "ibm"
page_title: "IBM : ibm_sm_secret_group"
description: |-
  Manages sm_secret_group.
subcategory: "Secrets Manager"
---

# ibm_sm_secret_group

Provides a resource for sm_secret_group. This allows sm_secret_group to be created, updated and deleted.

## Example Usage

```hcl
resource "ibm_sm_secret_group" "sm_secret_group" {
  secret_group_resource {
		id = "bc656587-8fda-4d05-9ad8-b1de1ec7e712"
		name = "my-secret-group"
		description = "Extended description for this group."
		creation_date = 2018-04-12T23:20:50.520Z
		last_update_date = 2018-05-12T23:20:50.520Z
		type = "application/vnd.ibm.secrets-manager.secret.group+json"
  }
}
```

## Argument Reference

Review the argument reference that you can specify for your resource.

* `secret_group_resource` - (Required, List) Properties that describe a secret group.
Nested scheme for **secret_group_resource**:
	* `creation_date` - (Optional, String) The date the secret group was created. The date format follows RFC 3339.
	* `description` - (Optional, String) An extended description of your secret group.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.
	  * Constraints: The maximum length is `1024` characters.
	* `id` - (Optional, String) The v4 UUID that uniquely identifies the secret group.
	* `last_update_date` - (Optional, String) Updates when the metadata of the secret group is modified. The date format follows RFC 3339.
	* `name` - (Optional, String) The type of policy. Allowable values are: rotation
	  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters.
	* `type` - (Optional, String) The MIME type that represents the secret group.

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your resource is created.

* `id` - The unique identifier of the sm_secret_group.
* `creation_date` - (Optional, String) The date the secret group was created. The date format follows RFC 3339.
* `description` - (Optional, String) An extended description of your secret group.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.
  * Constraints: The maximum length is `1024` characters.
* `last_update_date` - (Optional, String) Updates when the metadata of the secret group is modified. The date format follows RFC 3339.
* `name` - (Optional, String) The type of policy. Allowable values are: rotation
  * Constraints: The maximum length is `64` characters. The minimum length is `2` characters.
* `type` - (Optional, String) The MIME type that represents the secret group.

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

You can import the `ibm_sm_secret_group` resource by using `id`. The v4 UUID that uniquely identifies the secret group.
For more information, see [the documentation](https://cloud.ibm.com/docs/secrets-manager)

# Syntax
```
$ terraform import ibm_sm_secret_group.sm_secret_group <id>
```

# Example
```
$ terraform import ibm_sm_secret_group.sm_secret_group bc656587-8fda-4d05-9ad8-b1de1ec7e712
```
