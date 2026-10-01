package maas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MAAS only offers subnets on the interface's own VLAN unless the link is
// forced, so the provider has to ask for force every time or a subnet on
// another VLAN is rejected with a misleading "None found with id=...".
func TestGetNetworkInterfaceLinkParamsForcesLink(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceMAASNetworkInterfaceLink().Schema, map[string]any{
		"machine":           "abc123",
		"network_interface": "aa:bb:cc:dd:ee:ff",
		"subnet":            "10.0.0.0/24",
		"mode":              "STATIC",
		"ip_address":        "10.0.0.7",
		"default_gateway":   true,
	})

	params := getNetworkInterfaceLinkParams(d, 42)

	if !params.Force {
		t.Error("Force is false; MAAS will reject a subnet outside the interface's VLAN")
	}

	if params.Subnet != 42 {
		t.Errorf("Subnet = %d, want 42", params.Subnet)
	}

	if params.Mode != "STATIC" {
		t.Errorf("Mode = %q, want STATIC", params.Mode)
	}

	if params.IPAddress != "10.0.0.7" {
		t.Errorf("IPAddress = %q, want 10.0.0.7", params.IPAddress)
	}

	if !params.DefaultGateway {
		t.Error("DefaultGateway is false, want true")
	}
}
