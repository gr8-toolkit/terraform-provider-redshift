package redshift

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func forceNewIfListSizeChanged(key string) schema.CustomizeDiffFunc {
	return customdiff.ForceNewIfChange(key, listSizeChanged)
}

func listSizeChanged(_ context.Context, oldVal, newVal, _ interface{}) bool {
	old, ok := oldVal.([]interface{})
	if !ok {
		return false
	}
	nw, ok := newVal.([]interface{})
	if !ok {
		return false
	}
	return len(old) != len(nw)
}
