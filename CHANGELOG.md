## [v1.0.0] - 2017-12-29

*Breaking Change - Add support for Dynamic and List parameter types*

- string_key and string_value have been removed from parameter_detail
- values_to_readable_strings replaces string_key and string_value as a map[string]string. Each key in the map is 
effectively a separate string_key and the value is a separate string_value.
- The value of default_value must equal one of the keys (not value) within the values_to_readable_string map.

*Add missing fields to source*

- Allow disabled, scatter_plot_source, query_builder_enabled, source_description to optionally be applied to sources

## [v0.2.0] - 2018-01-03
*Updated README section on handling the creation of multiple alerts*
*Trim spaces on alert fields that support complex string types*

- Fixed #11 - condition, display_expression, and additional_information have been updated to call TrimSpaces. Preventing multiple plan/applies from re-applying the same state.

## [v0.1.2] - 2017-10-13

*Allow optional Alert attributes (as defined by the API) to be omitted from Terraform.*

- display_expression and resolve_after_minutes are now optional.

## [v0.1.1] - 2017-10-12

*Update Release process*

- Builds both linux and darwin versions of the plugin and uploads them all to github releases.

## [v0.1.0] - 2017-09-15

*First Release - Supports a limited Set of the Wavefront API*

- Support for Alerts, Dashboards and Alert Targets.
- Integration Testing of each resource
