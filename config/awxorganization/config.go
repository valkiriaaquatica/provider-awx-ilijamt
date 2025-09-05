package awxorganization

import (
	ujconfig "github.com/crossplane/upjet/pkg/config"
)

// Configure configures resources for the virtual environment group
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("awx_organization ", func(r *ujconfig.Resource) {
		r.ShortGroup = "AwxOrganization"
	})
}
