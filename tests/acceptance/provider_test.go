package tests

import (
	"fmt"
	"strconv"

	"github.com/criblio/terraform-provider-criblio/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	providerFactory = map[string]func() (tfprotov6.ProviderServer, error){
		"criblio": providerserver.NewProtocol6WithError(provider.New("999.99.9")()),
	}
)

func testCheckListDataSourceHasItems(name string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(name, "items.#", func(value string) error {
		count, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("%s items.# is not a number: %w", name, err)
		}
		if count == 0 {
			return fmt.Errorf("%s returned no items", name)
		}
		return nil
	})
}
