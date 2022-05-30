provider "ibm" {
  ibmcloud_api_key = var.ibmcloud_api_key
}

// Provision sm_secret_group resource instance
resource "ibm_sm_secret_group" "sm_secret_group_instance" {
  secret_group_resource {
    id = "bc656587-8fda-4d05-9ad8-b1de1ec7e712"
    name = "my-secret-group"
    description = "Extended description for this group."
    creation_date = 2018-04-12T23:20:50.520Z
    last_update_date = 2018-05-12T23:20:50.520Z
    type = "application/vnd.ibm.secrets-manager.secret.group+json"
  }
}

// Provision sm_secret resource instance
resource "ibm_sm_secret" "sm_secret_instance" {
  secret_type = var.sm_secret_secret_type
  secret_resource {
    id = "id"
    name = "name"
    description = "description"
    secret_group_id = "secret_group_id"
    labels = [ "labels" ]
    state = 1
    state_description = "Active"
    secret_type = "arbitrary"
    crn = "crn:v1:bluemix:public:secrets-manager:<region>:a/<account-id>:<service-instance>:secret:<secret-id>"
    creation_date = 2018-04-12T23:20:50.520Z
    created_by = "created_by"
    last_update_date = 2018-04-12T23:20:50.520Z
    versions_total = 1
    versions = [ { "key": null } ]
    expiration_date = 2030-04-01T09:30:00.000Z
    payload = "payload"
  }
}

// Provision sm_event_notification resource instance
resource "ibm_sm_event_notification" "sm_event_notification_instance" {
  event_notifications_instance_crn = var.sm_event_notification_event_notifications_instance_crn
  event_notifications_source_name = var.sm_event_notification_event_notifications_source_name
  event_notifications_source_description = var.sm_event_notification_event_notifications_source_description
}

// Provision sm_cert_configuration resource instance
resource "ibm_sm_cert_configuration" "sm_cert_configuration_instance" {
  secret_type = var.sm_cert_configuration_secret_type
  config_element = var.sm_cert_configuration_config_element
  name = var.sm_cert_configuration_name
  type = var.sm_cert_configuration_type
  config {
    private_key = "private_key"
  }
}

// Create sm_secret_groups data source
data "ibm_sm_secret_groups" "sm_secret_groups_instance" {
}

// Create sm_secrets data source
data "ibm_sm_secrets" "sm_secrets_instance" {
}

// Create sm_cert_configurations data source
data "ibm_sm_cert_configurations" "sm_cert_configurations_instance" {
  secret_type = var.sm_cert_configurations_secret_type
  config_element = var.sm_cert_configurations_config_element
}
