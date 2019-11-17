package wavefront_plugin

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func suppressCase(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}

func suppressSpaces(k, old, new string, d *schema.ResourceData) bool {
	if strings.TrimSpace(old) == strings.TrimSpace(new) {
		return true
	}
	return false
}

func trimSpaces(d interface{}) string {
	if s, ok := d.(string); ok {
		return strings.TrimSpace(s)
	}

	return ""
}

func trimSpacesMap(m map[string]interface{}) map[string]string {
	trimmed := map[string]string{}
	for key, v := range m {
		trimmed[key] = trimSpaces(v)
	}
	return trimmed
}
