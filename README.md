[![Build Status](https://travis-ci.com/spaceapegames/terraform-provider-wavefront.svg?token=bQMpYkjkzKf94BWWKiAA&branch=master)](https://travis-ci.com/spaceapegames/terraform-provider-wavefront)

# Wavefront Terraform Provider

A Terraform Provider to manage resources in Wavefront. Currently supports Alerts, the hope is to support Dashboards and more in the future.

## Requirements
* Go version 1.8 or higher
* Terraform 0.10.0 or higher (Custom providers were released at 0.10.0)
* [govendor](https://github.com/kardianos/govendor) for dependency management

## Known Issues

There is an issue with the wavefront API when applying tagged alerts that it will cause a race condition. They are working on fixing this.

To ensure that applies of more than one Alert are successful you can use  the `-parallelism` flag to prevent parallel resource creations
`terraform apply -parallelism=1`

## Building and Testing

### Build the plugin.

`make build`

This will create the plugin binary ./terraform-provider-wavefront

### Unit Test
`make test`

### Acceptance Tests
Acceptance tests are run against the Wavefront API so you'll need an account to use. Run at your own risk.

You need to supply the `WAVEFRONT_TOKEN` and `WAVEFRONT_ADDRESS` environment variables

To run the tests run
`make acceptance`

### Running the Plugin

Use the main.tf to create some test config, such as

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

* [This is a good blog post to get started](https://www.terraform.io/guides/writing-custom-terraform-providers.html?)
* [Some example Providers](https://github.com/terraform-providers)
