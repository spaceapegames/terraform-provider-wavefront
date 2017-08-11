[![Build Status](https://travis-ci.com/spaceapegames/terraform-provider-wavefront.svg?token=bQMpYkjkzKf94BWWKiAA&branch=master)](https://travis-ci.com/spaceapegames/terraform-provider-wavefront)

# Wavefront Terraform Provider

A Terraform Provider to manage resources in Wavefront. Currently supports Alerts, the hope is to support Dashboards in the future.

## Requirements

Go
Terraform 0.10.0 or higher (Custom providers were released at 0.10.0)
govendor for dependency management

## Building and Testing

### Build the plugin.

`make build`

 This will create the plugin binary ./terraform-provider-wavefront.
 The naming is important, terraform looks up providers using terraform-provider-<name>

### Unit Test
`make test`

### Acceptance Tests
Acceptance tests are run against the Wavefront API so you'll need an account to use. Run at your own risk.

To run acceptance tests you must set the `TF_ACC` environment variable
`export TF_ACC=true`

You also need to supply the `WAVEFRONT_TOKEN` and `WAVEFRONT_ADDRESS` environment variables

To run the tests run
`make acceptance`

### Integration Tests

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

 export your wavefront token `export WAVEFRONT_TOKEN=<token>` You could also configure the `token` in the provider section of main.tf, but best not to.

 Run `terraform init` to load your provider.

 Run `terraform plan` to show the plan.

 Run `terraform apply` to apply the test configuration and then check the results in Wavefront.

 Update main.tf to change a value, the run plan and apply again to check that updates work.

 Run `terraform destroy` to test deleting resources.

## Contributing

* [This is a good blog post to get started](https://www.terraform.io/guides/writing-custom-terraform-providers.html?)
* [Some example Providers](https://github.com/terraform-providers)

### Creating a new Resource

Create a new file `resource_<resource_type>.go`.

```
package main

// import the terraform helper schema
import (
    "github.com/hashicorp/terraform/helper/schema"
)

func resourceServer() *schema.Resource {
    return &schema.Resource{

        // Register the CRUD operations

        Create: resourceServerCreate,
        Read:   resourceServerRead,
        Update: resourceServerUpdate,
        Delete: resourceServerDelete,

        // Register the attributes of the resource

        Schema: map[string]*schema.Schema{
            "address": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
            },
        },
    }
}

// Implement the various crud operations

// Create a resource and set the terraform ID to be used as state - d.SetId("")
func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
    return nil
}
// Using the Id set by create locate the resource and it's current state.
// Set the ID to "" if it cannot be found ( has been deleted out of band)
func resourceServerRead(d *schema.ResourceData, m interface{}) error {
    return nil
}

// Apply any updates to the state of the resrouce
func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
    return nil
}

// Delete the resource
func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
    return nil
}
```