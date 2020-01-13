## [v2.1.2] - 2020-01-13

*Added Support for User Groups*

*Added more test cases for import logic*

- User Groups
- Users
- Derived Metrics

*Fixed issue where Derived Metrics were not reading Tags*

## [v2.1.1] - 2019-12-19

*Add support for Derived Metrics*

*Add Support for Alert Target Routes*

*Add Support for Users*

*Fixed issue where certain deleted resources would not properly detect to recreate resource on plan/apply*

## [v2.1.0] - 2019-07-03

*Add support for Threshold Alerts*

*Support creating dashboards from JSON*

## [v2.0.0] - 2019-06-11

*Upgrade to Terraform 0.12 to support new language features*

*May cause breaking changes due to new syntax ([See Upgrading to 0.12](https://www.terraform.io/upgrade-guides/0-12.html))*

In testing `values_to_readable_strings {` needed to change to `values_to_readable_strings = {` and `is_html_content = 1` changed to `is_html_content = true`

## [v1.0.1] - 2018-01-08

*Sort parameter details alphabetically to ensure no changes they are always evaluated in the correct order*

- Sort the parameter details when we 'read' a Dashboard from Wavefront

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
