---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "genesyscloud_telephony_providers_edges_site Data Source - terraform-provider-genesyscloud"
subcategory: ""
description: |-
  Data source for Genesys Cloud Sites. Select a site by name
---

# genesyscloud_telephony_providers_edges_site (Data Source)

Data source for Genesys Cloud Sites. Select a site by name

## Example Usage

```terraform
data "genesyscloud_telephony_providers_edges_site" "site" {
  name = "example site name"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Site name.

### Read-Only

- `id` (String) The ID of this resource.


