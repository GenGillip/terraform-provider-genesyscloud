package genesyscloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v95/platformclientv2"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud Groups. Select a group by name.",
		ReadContext: readWithPooledClient(dataSourceGroupRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Group name.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdkConfig := m.(*ProviderMeta).ClientConfig
	groupsAPI := platformclientv2.NewGroupsApiWithConfig(sdkConfig)

	exactSearchType := "EXACT"
	nameField := "name"
	nameStr := d.Get("name").(string)

	searchCriteria := platformclientv2.Groupsearchcriteria{
		VarType: &exactSearchType,
		Value:   &nameStr,
		Fields:  &[]string{nameField},
	}

	return withRetries(ctx, 15*time.Second, func() *resource.RetryError {
		groups, _, getErr := groupsAPI.PostGroupsSearch(platformclientv2.Groupsearchrequest{
			Query: &[]platformclientv2.Groupsearchcriteria{searchCriteria},
		})
		if getErr != nil {
			return resource.NonRetryableError(fmt.Errorf("Error requesting group %s: %s", nameStr, getErr))
		}

		if *groups.Total == 0 {
			return resource.RetryableError(fmt.Errorf("No groups found with search criteria %v ", searchCriteria))
		}

		// Select first group in the list
		group := (*groups.Results)[0]
		d.SetId(*group.Id)
		return nil
	})
}
