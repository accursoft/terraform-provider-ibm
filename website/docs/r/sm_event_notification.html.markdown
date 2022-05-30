---
layout: "ibm"
page_title: "IBM : ibm_sm_event_notification"
description: |-
  Manages sm_event_notification.
subcategory: "Secrets Manager"
---

# ibm_sm_event_notification

Provides a resource for sm_event_notification. This allows sm_event_notification to be created, updated and deleted.

## Example Usage

```hcl
resource "ibm_sm_event_notification" "sm_event_notification" {
  event_notifications_instance_crn = "crn:v1:bluemix:public:event-notifications:us-south:a/<account-id>:<service-instance>::"
  event_notifications_source_description = "Optional description of this source in an Event Notifications instance."
  event_notifications_source_name = "My Secrets Manager"
}
```

## Argument Reference

Review the argument reference that you can specify for your resource.

* `event_notifications_instance_crn` - (Required, Forces new resource, String) The Cloud Resource Name (CRN) of the connected Event Notifications instance.
* `event_notifications_source_description` - (Optional, Forces new resource, String) An optional description for the source in your Event Notifications instance.
* `event_notifications_source_name` - (Required, Forces new resource, String) The name that is displayed as a source in your Event Notifications instance.

## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your resource is created.

* `id` - The unique identifier of the sm_event_notification.

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

You can import the `ibm_sm_event_notification` resource by using `event_notifications_instance_crn`. The Cloud Resource Name (CRN) of the connected Event Notifications instance.
For more information, see [the documentation](https://cloud.ibm.com/docs/secrets-manager)

# Syntax
```
$ terraform import ibm_sm_event_notification.sm_event_notification <event_notifications_instance_crn>
```

# Example
```
$ terraform import ibm_sm_event_notification.sm_event_notification crn:v1:bluemix:public:event-notifications:us-south:a/<account-id>:<service-instance>::
```
