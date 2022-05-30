variable "ibmcloud_api_key" {
  description = "IBM Cloud API key"
  type        = string
}

// Resource arguments for sm_secret_group

// Resource arguments for sm_secret
variable "sm_secret_secret_type" {
  description = "The secret type. Allowable values are: arbitrary, iam_credentials, imported_cert, public_cert, username_password, kv."
  type        = string
  default     = "arbitrary"
}

// Resource arguments for sm_event_notification
variable "sm_event_notification_event_notifications_instance_crn" {
  description = "The Cloud Resource Name (CRN) of the connected Event Notifications instance."
  type        = string
  default     = "crn:v1:bluemix:public:event-notifications:us-south:a/<account-id>:<service-instance>::"
}
variable "sm_event_notification_event_notifications_source_name" {
  description = "The name that is displayed as a source in your Event Notifications instance."
  type        = string
  default     = "My Secrets Manager"
}
variable "sm_event_notification_event_notifications_source_description" {
  description = "An optional description for the source in your Event Notifications instance."
  type        = string
  default     = "Optional description of this source in an Event Notifications instance."
}

// Resource arguments for sm_cert_configuration
variable "sm_cert_configuration_secret_type" {
  description = "The secret type. Allowable values are: public_cert, private_cert"
  type        = string
  default     = "public_cert"
}
variable "sm_cert_configuration_config_element" {
  description = "The configuration element to define or manage. Allowable values are: certificate_authorities, dns_providers, root_certificate_authorities, intermediate_certificate_authorities, certificate_templates"
  type        = string
  default     = "certificate_authorities"
}
variable "sm_cert_configuration_name" {
  description = "The human-readable name to assign to your configuration."
  type        = string
  default     = "cis-example-config"
}
variable "sm_cert_configuration_type" {
  description = "The type of configuration. Value options differ depending on the 'config_element'. Allowable values are: letsencrypt, letsencrypt-stage, cis, classic_infrastructure, root_certificate_authority, intermediate_certificate_authority"
  type        = string
  default     = "cis"
}

// Data source arguments for sm_secret_groups
variable "sm_secret_groups_id" {
  description = "The ID of the secret group."
  type        = string
  default     = "id"
}


// Data source arguments for sm_cert_configurations
variable "sm_cert_configurations_secret_type" {
  description = "The secret type. Allowable values are: public_cert, private_cert"
  type        = string
  default     = "public_cert"
}
variable "sm_cert_configurations_config_element" {
  description = "The configuration element to define or manage. Allowable values are: certificate_authorities, dns_providers, root_certificate_authorities, intermediate_certificate_authorities, certificate_templates"
  type        = string
  default     = "certificate_authorities"
}
