package genesyscloud

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mypurecloud/platform-client-sdk-go/v95/platformclientv2"
)

func dataSourceWebDeploymentsDeployment() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for Genesys Cloud Web Deployments. Select a deployment by name.",
		ReadContext: readWithPooledClient(dataSourceDeploymentRead),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the deployment",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdkConfig := m.(*ProviderMeta).ClientConfig
	api := platformclientv2.NewWebDeploymentsApiWithConfig(sdkConfig)

	name := d.Get("name").(string)

	return withRetries(ctx, 15*time.Second, func() *resource.RetryError {
		deployments, resp, err := api.GetWebdeploymentsDeployments([]string{})

		if err != nil && resp.StatusCode == http.StatusNotFound {
			return resource.RetryableError(fmt.Errorf("No web deployment record found %s: %s. Correlation id: %s", name, err, resp.CorrelationID))
		}

		if err != nil && resp.StatusCode != http.StatusNotFound {
			return resource.NonRetryableError(fmt.Errorf("Error retrieving web deployment %s: %s. Correlation id: %s", name, err, resp.CorrelationID))
		}

		for _, deployment := range *deployments.Entities {
			if name == *deployment.Name {
				d.SetId(*deployment.Id)
				return nil
			}
		}

		return resource.NonRetryableError(fmt.Errorf("No web deployment was found with the name %s", name))
	})
}
