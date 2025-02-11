---
page_title: "genesyscloud_knowledge_category Resource - terraform-provider-genesyscloud"
subcategory: ""
description: |-
  Genesys Cloud Knowledge Category
---
# genesyscloud_knowledge_category (Resource)

Genesys Cloud Knowledge Category

## API Usage
The following Genesys Cloud APIs are used by this resource. Ensure your OAuth Client has been granted the necessary scopes and permissions to perform these operations:

* [GET /api/v2/knowledge/knowledgebases/{knowledgeBaseId}/languages/{languageCode}/categories](https://developer.genesys.cloud/api/rest/v2/knowledge/#post-api-v2-knowledge-knowledgebases--knowledgeBaseId--languages--languageCode--categories)
* [POST /api/v2/knowledge/knowledgebases/{knowledgeBaseId}/languages/{languageCode}/categories](https://developer.genesys.cloud/api/rest/v2/knowledge/#post-api-v2-knowledge-knowledgebases--knowledgeBaseId--languages--languageCode--categories)
* [GET /api/v2/knowledge/knowledgebases/{knowledgeBaseId}/languages/{languageCode}/categories/{categoryId}](https://developer.mypurecloud.com/api/rest/v2/knowledge/#get-api-v2-knowledge-knowledgebases--knowledgeBaseId--languages--languageCode--categories--categoryId-)
* [PATCH /api/v2/knowledge/knowledgebases/{knowledgeBaseId}/languages/{languageCode}/categories](https://developer.mypurecloud.com/api/rest/v2/knowledge/#patch-api-v2-knowledge-knowledgebases--knowledgeBaseId--languages--languageCode--categories)
* [DELETE /api/v2/knowledge/knowledgebases/{knowledgeBaseId}/languages/{languageCode}/categories/{categoryId}](https://developer.mypurecloud.com/api/rest/v2/knowledge/#delete-api-v2-knowledge-knowledgebases--knowledgeBaseId--languages--languageCode--categories--categoryId-)

## Example Usage

```terraform
resource "genesyscloud_knowledge_category" "example_category" {
  knowledge_base_id = genesyscloud_knowledge_knowledgebase.example_knowledgebase.id
  language_code     = "en-US"
  knowledge_category {
    name        = "ExampleCategory"
    description = "An example category"
    parent_id   = genesyscloud_knowledge_category.parent_category.id
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `knowledge_base_id` (String) Knowledge base id of the category
- `knowledge_category` (Block List, Min: 1, Max: 1) Knowledge category parent id (see [below for nested schema](#nestedblock--knowledge_category))
- `language_code` (String) language code of the category

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--knowledge_category"></a>
### Nested Schema for `knowledge_category`

Optional:

- `description` (String) Knowledge base description
- `name` (String) Knowledge base name. Changing the name attribute will cause the knowledge_category resource to be dropped and recreated with a new ID.
- `parent_id` (String) Knowledge category parent id

