[![Build Status](https://travis-ci.com/spaceapegames/terraform-provider-wavefront.svg?token=bQMpYkjkzKf94BWWKiAA&branch=master)](https://travis-ci.com/spaceapegames/terraform-provider-wavefront)

# Wavefront Terraform Provider

A Terraform Provider to manage resources in Wavefront. Currently supports Alerts, Alert Targets and Dashboards.

## Requirements
* Go version 1.8 or higher
* Terraform 0.10.0 or higher (Custom providers were released at 0.10.0)
* [govendor](https://github.com/kardianos/govendor) for dependency management

## Installing the Plugin

We release darwin and linux amd64 packages on the [releases page](https://github.com/spaceapegames/terraform-provider-wavefront/releases). If you require a different architecture you will need to build the plugin from source, see below for more details:

Once you have the plugin you should remove the `_os_arch` from the end of the file name and place it in `~/.terraform.d/plugins` which is where `terraform init` will look for plugins.

Valid provider filenames are `terraform-provider-NAME_X.X.X` or `terraform-provider-NAME_vX.X.X`

## Known Issues

There was an issue with the wavefront API when applying tagged alerts that it will cause a race condition. This has been fixed by Wavefront, but it's possible that not all clusters have been upgraded to the version containing the fix. If you run into this issue you can use the following solution until your cluster has been upgraded.

To ensure that applies of more than one Alert are successful you can use  the `-parallelism` flag to prevent parallel resource creations
`terraform apply -parallelism=1`

## Building and Testing

### Build the plugin.

`make build`

This will build amd64 arch binaries for darwin and linux in the format terraform-provider-wavefront_<version>_<targetOS>_<arch>

### Unit Test
`make test`

### Acceptance Tests
Acceptance tests are run against the Wavefront API so you'll need an account to use. Run at your own risk.

You need to supply the `WAVEFRONT_TOKEN` and `WAVEFRONT_ADDRESS` environment variables

To run the tests run
`make acceptance`

### Running the Plugin

Use a main.tf to create some test config, such as

```
 provider "wavefront" {
   address = "spaceape.wavefront.com"
 }

 resource "wavefront_alert" "test_alert" {
   name = "Terraform Test Alert"
   target = "test@example.com"
   condition = "100-ts(\"cpu.usage_idle\", environment=flamingo-int and cpu=cpu-total and service=game-service) > 80"
   display_expression = "100-ts(\"cpu.usage_idle\", environment=flamingo-int and cpu=cpu-total and service=game-service)"
   minutes = 5
   resolve_after_minutes = 5
   severity = "WARN"
   tags = [
     "terraform",
     "flamingo"
   ]
 }
```

Export your wavefront token `export WAVEFRONT_TOKEN=<token>` You could also configure the `token` in the provider section of main.tf, but best not to.

Run `terraform init` to load your provider.

Run `terraform plan` to show the plan.

Run `terraform apply` to apply the test configuration and then check the results in Wavefront.

Update main.tf to change a value, the run plan and apply again to check that updates work.

Run `terraform destroy` to test deleting resources.

## Contributing

Please review the CONTRIBUTOR.md document for more information on contributing.