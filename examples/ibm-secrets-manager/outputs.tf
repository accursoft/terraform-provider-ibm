// This allows sm_secret_group data to be referenced by other resources and the terraform CLI
// Modify this if only certain data should be exposed
output "ibm_sm_secret_group" {
  value       = ibm_sm_secret_group.sm_secret_group_instance
  description = "sm_secret_group resource instance"
}
// This allows sm_secret data to be referenced by other resources and the terraform CLI
// Modify this if only certain data should be exposed
output "ibm_sm_secret" {
  value       = ibm_sm_secret.sm_secret_instance
  description = "sm_secret resource instance"
}
// This allows sm_event_notification data to be referenced by other resources and the terraform CLI
// Modify this if only certain data should be exposed
output "ibm_sm_event_notification" {
  value       = ibm_sm_event_notification.sm_event_notification_instance
  description = "sm_event_notification resource instance"
}
// This allows sm_cert_configuration data to be referenced by other resources and the terraform CLI
// Modify this if only certain data should be exposed
output "ibm_sm_cert_configuration" {
  value       = ibm_sm_cert_configuration.sm_cert_configuration_instance
  description = "sm_cert_configuration resource instance"
}
