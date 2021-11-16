package vault

import (
	"fmt"
)

// EnginesFromVault takes a path and returns a list of engines from vault.
func (vc *vaultClient) EnginesFromVault(path string) ([]string, error) {
	vc.tracer.trace(fmt.Sprintf("%s/EnginesFromVault", vc.config.tracePrefix))

	engines, err := vc.client.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("listing engines from Vault for %s", path)
	}

	if engines == nil {
		return nil, fmt.Errorf("engines returned from Vault are <nil> for %s", path)
	}

	engineData, _ := extractListData(engines)

	result := []string{}

	for _, value := range engineData {
		switch v := value.(type) {
		case string:
			result = append(result, v)
		default:
			return nil, fmt.Errorf("unexpected type, expected string, got: %T, value: %v", v, result)
		}
	}
	return result, nil
}
