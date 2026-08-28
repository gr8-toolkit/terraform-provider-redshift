# Redshift Serverless — temporary credentials via IAM
#
# The provider calls redshift-serverless:GetCredentials to obtain a short-lived
# password for the workgroup. The database user is derived from the caller's IAM
# identity (role/user name with an "IAMR:" prefix).
#
# Note: binary_parameters=no is automatically added to the connection string for
# Serverless to avoid "unexpected Bind response" errors on GRANT/REVOKE statements
# caused by Serverless's limited support of the PostgreSQL extended query protocol.
provider "redshift" {
  host          = var.redshift_serverless_host
  username      = var.redshift_user
  database      = var.redshift_database
  is_serverless = true

  temporary_credentials {
    workgroup_name = "my-workgroup"
  }
}
