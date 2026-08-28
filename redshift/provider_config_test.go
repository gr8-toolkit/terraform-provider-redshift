package redshift

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerSchema returns the provider's top-level schema map for use in tests.
func providerSchema() map[string]*schema.Schema {
	return Provider().Schema
}

// TestConnParams_Serverless verifies that binary_parameters=no is added to the
// DSN when IsServerless is true. Without this lib/pq uses the extended query
// protocol, which causes "unexpected Bind response" errors on GRANT/REVOKE
// statements against Redshift Serverless.
func TestConnParams_Serverless(t *testing.T) {
	cfg := Config{
		SSLMode:      "require",
		IsServerless: true,
	}
	joined := strings.Join(cfg.connParams(), "&")

	if !strings.Contains(joined, "binary_parameters=no") {
		t.Errorf("expected binary_parameters=no in Serverless DSN params, got: %s", joined)
	}
}

// TestConnParams_Provisioned verifies that binary_parameters is NOT added for
// provisioned clusters — the extended protocol works fine there.
func TestConnParams_Provisioned(t *testing.T) {
	cfg := Config{
		SSLMode:      "require",
		IsServerless: false,
	}
	joined := strings.Join(cfg.connParams(), "&")

	if strings.Contains(joined, "binary_parameters") {
		t.Errorf("binary_parameters should not be set for provisioned clusters, got: %s", joined)
	}
}

// TestConnStr_Serverless_ContainsBinaryParametersNo verifies the full DSN
// produced for a Serverless config includes binary_parameters=no.
func TestConnStr_Serverless_ContainsBinaryParametersNo(t *testing.T) {
	cfg := Config{
		Host:         "my-workgroup.123.us-east-1.redshift-serverless.amazonaws.com",
		Port:         5439,
		Username:     "admin",
		Password:     "secret",
		Database:     "dev",
		SSLMode:      "require",
		IsServerless: true,
	}
	dsn := cfg.connStr("dev")

	if !strings.Contains(dsn, "binary_parameters=no") {
		t.Errorf("expected binary_parameters=no in Serverless DSN, got: %s", dsn)
	}
}

// TestConnStr_Provisioned_NoBinaryParameters verifies the full DSN for a
// provisioned cluster does not include binary_parameters.
func TestConnStr_Provisioned_NoBinaryParameters(t *testing.T) {
	cfg := Config{
		Host:         "my-cluster.us-east-1.redshift.amazonaws.com",
		Port:         5439,
		Username:     "admin",
		Password:     "secret",
		Database:     "dev",
		SSLMode:      "require",
		IsServerless: false,
	}
	dsn := cfg.connStr("dev")

	if strings.Contains(dsn, "binary_parameters") {
		t.Errorf("binary_parameters should not be in provisioned DSN, got: %s", dsn)
	}
}

// TestResolveCredentials_PasswordAuth verifies that when no temporary_credentials
// block is present the provider falls through to plain password authentication.
func TestResolveCredentials_PasswordAuth(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, providerSchema(), map[string]interface{}{
		"username": "admin",
		"password": "secret",
	})

	user, pass, err := resolveCredentials(rd)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if user != "admin" {
		t.Errorf("expected username 'admin', got '%s'", user)
	}
	if pass != "secret" {
		t.Errorf("expected password 'secret', got '%s'", pass)
	}
}

// TestResolveCredentials_MissingUsername verifies that omitting username returns
// a clear error rather than panicking or silently producing empty credentials.
func TestResolveCredentials_MissingUsername(t *testing.T) {
	// Use a minimal ResourceData with username explicitly set to empty string,
	// overriding the DefaultFunc ("root") so GetOk returns false.
	rd := schema.TestResourceDataRaw(t, providerSchema(), map[string]interface{}{
		"username": "",
		"password": "secret",
	})

	_, _, err := resolveCredentials(rd)
	if err == nil {
		t.Fatal("expected error for missing username, got nil")
	}
	if !strings.Contains(err.Error(), "username is required") {
		t.Errorf("unexpected error message: %s", err)
	}
}

// TestTemporaryCredentials_MissingIdentifier verifies that a temporary_credentials
// block with neither cluster_identifier nor workgroup_name returns a clear error
// that mentions both fields.
func TestTemporaryCredentials_MissingIdentifier(t *testing.T) {
	rd := schema.TestResourceDataRaw(t, providerSchema(), map[string]interface{}{
		"username": "admin",
		"temporary_credentials": []interface{}{
			map[string]interface{}{
				"cluster_identifier": "",
				"workgroup_name":     "",
				"region":             "",
				"auto_create_user":   false,
				"duration_seconds":   0,
				"db_groups":          []interface{}{},
				"assume_role":        []interface{}{},
			},
		},
	})

	_, _, err := temporaryCredentials("admin", rd)
	if err == nil {
		t.Fatal("expected error when neither cluster_identifier nor workgroup_name is set")
	}
	if !strings.Contains(err.Error(), "workgroup_name") || !strings.Contains(err.Error(), "cluster_identifier") {
		t.Errorf("error should mention both fields, got: %s", err)
	}
}
