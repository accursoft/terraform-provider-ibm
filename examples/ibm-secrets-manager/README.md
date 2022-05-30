# Example for SecretsManagerV1

This example illustrates how to use the SecretsManagerV1

These types of resources are supported:

* sm_secret_group
* sm_secret
* sm_event_notification
* sm_cert_configuration

## Usage

To run this example you need to execute:

```bash
$ terraform init
$ terraform plan
$ terraform apply
```

Run `terraform destroy` when you don't need these resources.


## SecretsManagerV1 resources

sm_secret_group resource:

```hcl
resource "sm_secret_group" "sm_secret_group_instance" {
  metadata = var.sm_secret_group_metadata
  resources = var.sm_secret_group_resources
}
```
sm_secret resource:

```hcl
resource "sm_secret" "sm_secret_instance" {
  secret_type = var.sm_secret_secret_type
  metadata = var.sm_secret_metadata
  resources = var.sm_secret_resources
}
```
sm_event_notification resource:

```hcl
resource "sm_event_notification" "sm_event_notification_instance" {
  event_notifications_instance_crn = var.sm_event_notification_event_notifications_instance_crn
  event_notifications_source_name = var.sm_event_notification_event_notifications_source_name
  event_notifications_source_description = var.sm_event_notification_event_notifications_source_description
}
```
sm_cert_configuration resource:

```hcl
resource "sm_cert_configuration" "sm_cert_configuration_instance" {
  secret_type = var.sm_cert_configuration_secret_type
  config_element = var.sm_cert_configuration_config_element
  name = var.sm_cert_configuration_name
  type = var.sm_cert_configuration_type
  config = var.sm_cert_configuration_config
}
```

## SecretsManagerV1 Data sources

sm_secret_groups data source:

```hcl
data "sm_secret_groups" "sm_secret_groups_instance" {
  id = var.sm_secret_groups_id
}
```
sm_secrets data source:

```hcl
data "sm_secrets" "sm_secrets_instance" {
}
```
sm_cert_configurations data source:

```hcl
data "sm_cert_configurations" "sm_cert_configurations_instance" {
  secret_type = var.sm_cert_configurations_secret_type
  config_element = var.sm_cert_configurations_config_element
}
```

## Assumptions

1. TODO

## Notes

1. TODO

## Requirements

| Name | Version |
|------|---------|
| terraform | ~> 0.12 |

## Providers

| Name | Version |
|------|---------|
| ibm | 1.13.1 |

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|---------|
| ibmcloud\_api\_key | IBM Cloud API key | `string` | true |
| metadata | The metadata that describes the resource array. | `` | true |
| resources | The properties in JSON format to define, such as the name and description. For more information, see the docs: https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-cli-plugin-secrets-manager-cli#secrets-manager-cli-secret-group-create-command | `list()` | true |
| secret_type | The secret type. Allowable values are: arbitrary, iam_credentials, imported_cert, public_cert, username_password, kv. | `string` | true |
| metadata | The metadata that describes the resource array. | `` | true |
| resources | The properties in JSON format to store for the secret. Properties differ depending on the secret type. For more information, see the docs: https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-cli-plugin-secrets-manager-cli#secrets-manager-cli-secret-create-command | `list()` | true |
| event_notifications_instance_crn | The Cloud Resource Name (CRN) of the connected Event Notifications instance. | `string` | true |
| event_notifications_source_name | The name that is displayed as a source in your Event Notifications instance. | `string` | true |
| event_notifications_source_description | An optional description for the source in your Event Notifications instance. | `string` | false |
| secret_type | The secret type. Allowable values are: public_cert, private_cert | `string` | true |
| config_element | The configuration element to define or manage. Allowable values are: certificate_authorities, dns_providers, root_certificate_authorities, intermediate_certificate_authorities, certificate_templates | `string` | true |
| name | The human-readable name to assign to your configuration. | `string` | true |
| type | The type of configuration. Value options differ depending on the 'config_element'. Allowable values are: letsencrypt, letsencrypt-stage, cis, classic_infrastructure, root_certificate_authority, intermediate_certificate_authority | `string` | true |
| config | The configuration to define for the specified secret type. | `` | true |
| id | The ID of the secret group. | `string` | true |
| secret_type | The secret type. Allowable values are: public_cert, private_cert | `string` | true |
| config_element | The configuration element to define or manage. Allowable values are: certificate_authorities, dns_providers, root_certificate_authorities, intermediate_certificate_authorities, certificate_templates | `string` | true |

## Outputs

| Name | Description |
|------|-------------|
| sm_secret_group | sm_secret_group object |
| sm_secret | sm_secret object |
| sm_event_notification | sm_event_notification object |
| sm_cert_configuration | sm_cert_configuration object |
| sm_secret_groups | sm_secret_groups object |
| sm_secrets | sm_secrets object |
| sm_cert_configurations | sm_cert_configurations object |
