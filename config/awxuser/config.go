/*
Copyright 2021 Upbound Inc.
*/

package awxuser

import (
	ujconfig "github.com/crossplane/upjet/pkg/config"
)

// Configure configures the null group
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("awx_user", func(r *ujconfig.Resource) {
		r.ShortGroup = "AwxUser"
		// And other overrides.
	})
}
