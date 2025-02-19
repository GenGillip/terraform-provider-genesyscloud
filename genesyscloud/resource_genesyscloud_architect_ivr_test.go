package genesyscloud

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mypurecloud/platform-client-sdk-go/v95/platformclientv2"
)

type ivrConfigStruct struct {
	resourceID  string
	name        string
	description string
	dnis        []string
	depends_on  string
	divisionId  string
}

func deleteIvrStartingWith(name string) {
	archAPI := platformclientv2.NewArchitectApiWithConfig(sdkConfig)

	for pageNum := 1; ; pageNum++ {
		const pageSize = 100
		ivrs, _, getErr := archAPI.GetArchitectIvrs(pageNum, pageSize, "", "", "", "", "")
		if getErr != nil {
			return
		}

		if ivrs.Entities == nil || len(*ivrs.Entities) == 0 {
			break
		}

		for _, ivr := range *ivrs.Entities {
			if strings.HasPrefix(*ivr.Name, name) {
				archAPI.DeleteArchitectIvr(*ivr.Id)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

func TestAccResourceIvrConfigBasic(t *testing.T) {
	ivrConfigResource1 := "test-ivrconfig1"
	ivrConfigName := "terraform-ivrconfig-" + uuid.NewString()
	ivrConfigDescription := "Terraform IVR config"
	number1 := "+14175550011"
	number2 := "+14175550012"
	err := authorizeSdk()
	if err != nil {
		t.Fatal(err)
	}
	deleteIvrStartingWith("terraform-ivrconfig-")
	deleteDidPoolWithNumber(number1)
	deleteDidPoolWithNumber(number2)
	ivrConfigDnis := []string{number1, number2}
	didPoolResource1 := "test-didpool1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: generateIvrConfigResource(&ivrConfigStruct{
					resourceID:  ivrConfigResource1,
					name:        ivrConfigName,
					description: ivrConfigDescription,
					dnis:        nil, // No dnis
					depends_on:  "",  // No depends_on
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_architect_ivr."+ivrConfigResource1, "name", ivrConfigName),
					resource.TestCheckResourceAttr("genesyscloud_architect_ivr."+ivrConfigResource1, "description", ivrConfigDescription),
					hasEmptyDnis("genesyscloud_architect_ivr."+ivrConfigResource1),
				),
			},
			{
				// Update with new DNIS
				Config: generateDidPoolResource(&didPoolStruct{
					didPoolResource1,
					ivrConfigDnis[0],
					ivrConfigDnis[1],
					nullValue, // No description
					nullValue, // No comments
					nullValue, // No provider
				}) + generateIvrConfigResource(&ivrConfigStruct{
					resourceID:  ivrConfigResource1,
					name:        ivrConfigName,
					description: ivrConfigDescription,
					dnis:        ivrConfigDnis,
					depends_on:  "genesyscloud_telephony_providers_edges_did_pool." + didPoolResource1,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_architect_ivr."+ivrConfigResource1, "name", ivrConfigName),
					resource.TestCheckResourceAttr("genesyscloud_architect_ivr."+ivrConfigResource1, "description", ivrConfigDescription),
					validateStringInArray("genesyscloud_architect_ivr."+ivrConfigResource1, "dnis", ivrConfigDnis[0]),
					validateStringInArray("genesyscloud_architect_ivr."+ivrConfigResource1, "dnis", ivrConfigDnis[1]),
				),
			},
			{
				// Import/Read
				ResourceName:      "genesyscloud_architect_ivr." + ivrConfigResource1,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testVerifyIvrConfigsDestroyed,
	})
}

func TestAccResourceIvrConfigDivision(t *testing.T) {
	ivrConfigResource1 := "test-ivrconfig1"
	ivrConfigName := "terraform-ivrconfig-" + uuid.NewString()
	ivrConfigDescription := "Terraform IVR config"
	number1 := "+14175550011"
	number2 := "+14175550012"
	divResource1 := "auth-division1"
	divResource2 := "auth-division2"
	divName1 := "TerraformDiv-" + uuid.NewString()
	divName2 := "TerraformDiv-" + uuid.NewString()
	err := authorizeSdk()
	if err != nil {
		t.Fatal(err)
	}
	deleteIvrStartingWith("terraform-ivrconfig-")
	deleteDidPoolWithNumber(number1)
	deleteDidPoolWithNumber(number2)
	ivrConfigDnis := []string{number1, number2}
	didPoolResource1 := "test-didpool1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create
				Config: generateAuthDivisionResource(
					divResource1,
					divName1,
					nullValue, // No description
					nullValue, // Not home division
				) + generateIvrConfigResource(&ivrConfigStruct{
					resourceID:  ivrConfigResource1,
					name:        ivrConfigName,
					description: ivrConfigDescription,
					dnis:        nil, // No dnis
					depends_on:  "",  // No depends_on
					divisionId:  "genesyscloud_auth_division." + divResource1 + ".id",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_architect_ivr."+ivrConfigResource1, "name", ivrConfigName),
					resource.TestCheckResourceAttr("genesyscloud_architect_ivr."+ivrConfigResource1, "description", ivrConfigDescription),
					resource.TestCheckResourceAttrPair("genesyscloud_architect_ivr."+ivrConfigResource1, "division_id", "genesyscloud_auth_division."+divResource1, "id"),
					hasEmptyDnis("genesyscloud_architect_ivr."+ivrConfigResource1),
				),
			},
			{
				// Update with new DNIS and division
				Config: generateAuthDivisionResource(
					divResource1,
					divName1,
					nullValue, // No description
					nullValue, // Not home division
				) + generateAuthDivisionResource(
					divResource2,
					divName2,
					nullValue, // No description
					nullValue, // Not home division
				) + generateDidPoolResource(&didPoolStruct{
					didPoolResource1,
					ivrConfigDnis[0],
					ivrConfigDnis[1],
					nullValue, // No description
					nullValue, // No comments
					nullValue, // No provider
				}) + generateIvrConfigResource(&ivrConfigStruct{
					resourceID:  ivrConfigResource1,
					name:        ivrConfigName,
					description: ivrConfigDescription,
					dnis:        ivrConfigDnis,
					depends_on:  "genesyscloud_telephony_providers_edges_did_pool." + didPoolResource1,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_architect_ivr."+ivrConfigResource1, "name", ivrConfigName),
					resource.TestCheckResourceAttr("genesyscloud_architect_ivr."+ivrConfigResource1, "description", ivrConfigDescription),
					resource.TestCheckResourceAttrPair("genesyscloud_architect_ivr."+ivrConfigResource1, "division_id", "genesyscloud_auth_division."+divResource1, "id"),
					validateStringInArray("genesyscloud_architect_ivr."+ivrConfigResource1, "dnis", ivrConfigDnis[0]),
					validateStringInArray("genesyscloud_architect_ivr."+ivrConfigResource1, "dnis", ivrConfigDnis[1]),
				),
			},
			{
				// Import/Read
				ResourceName:      "genesyscloud_architect_ivr." + ivrConfigResource1,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: generateAuthDivisionResource(
					divResource1,
					divName1,
					nullValue, // No description
					nullValue, // Not home division
				) + generateAuthDivisionResource(
					divResource2,
					divName2,
					nullValue, // No description
					nullValue, // Not home division
				),
			},
		},
		CheckDestroy: testVerifyIvrConfigsDestroyed,
	})
}

func generateIvrConfigResource(ivrConfig *ivrConfigStruct) string {
	dnisStrs := make([]string, len(ivrConfig.dnis))
	for i, val := range ivrConfig.dnis {
		dnisStrs[i] = fmt.Sprintf("\"%s\"", val)
	}

	divisionId := ""
	if ivrConfig.divisionId != "" {
		divisionId = "division_id = " + ivrConfig.divisionId
	}

	return fmt.Sprintf(`resource "genesyscloud_architect_ivr" "%s" {
		name = "%s"
		description = "%s"
		dnis = [%s]
		depends_on=[%s]
		%s
	}
	`, ivrConfig.resourceID,
		ivrConfig.name,
		ivrConfig.description,
		strings.Join(dnisStrs, ","),
		ivrConfig.depends_on,
		divisionId,
	)
}

func testVerifyIvrConfigsDestroyed(state *terraform.State) error {
	architectApi := platformclientv2.NewArchitectApi()
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "genesyscloud_architect_ivr" {
			continue
		}

		ivrConfig, resp, err := architectApi.GetArchitectIvr(rs.Primary.ID)
		if ivrConfig != nil && ivrConfig.State != nil && *ivrConfig.State == "deleted" {
			continue
		}

		if ivrConfig != nil {
			return fmt.Errorf("IVR config (%s) still exists", rs.Primary.ID)
		}

		if isStatus404(resp) {
			// IVR Config not found as expected
			continue
		}

		// Unexpected error
		return fmt.Errorf("Unexpected error: %s", err)
	}
	// Success. All IVR Config pool destroyed
	return nil
}

func hasEmptyDnis(ivrResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		ivrResource, ok := state.RootModule().Resources[ivrResourceName]
		if !ok {
			return fmt.Errorf("Failed to find ivr config %s in state", ivrResourceName)
		}
		ivrID := ivrResource.Primary.ID

		dnisCountStr, ok := ivrResource.Primary.Attributes["dnis.#"]
		if !ok {
			return fmt.Errorf("No dnis found for %s in state", ivrID)
		}

		dnisCount, err := strconv.Atoi(dnisCountStr)
		if err != nil {
			return fmt.Errorf("Error while converting dnis count")
		}

		if dnisCount > 0 {
			return fmt.Errorf("Dnis is not empty.")
		}

		return nil
	}
}
