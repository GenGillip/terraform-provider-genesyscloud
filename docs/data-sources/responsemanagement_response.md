---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "genesyscloud_responsemanagement_response Data Source - terraform-provider-genesyscloud"
subcategory: ""
description: |-
  Data source for Genesys Cloud Responsemanagement Response. Select a Responsemanagement Response by name.
---

# genesyscloud_responsemanagement_response (Data Source)

Data source for Genesys Cloud Responsemanagement Response. Select a Responsemanagement Response by name.

## Example Usage

```terraform
data "genesyscloud_responsemanagement_response" "example_responsemanagement_response" {
  name    = "Responsemanagement response"
  library = genesyscloud_responsemanagement_library.library_1.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `library_id` (String) ID of the library that contains the response.
- `name` (String) Responsemanagement Response name.

### Read-Only

- `id` (String) The ID of this resource.


