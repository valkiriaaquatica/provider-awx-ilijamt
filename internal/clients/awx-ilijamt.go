/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/pkg/terraform"

	"github.com/valkiriaaquatica/provider-awx-ilijamt/apis/v1beta1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal awx-ilijamt credentials as JSON"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, c client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		ref := mg.GetProviderConfigReference()
		if ref == nil {
			return ps, errors.New(errNoProviderConfig)
		}
		pc := &v1beta1.ProviderConfig{}
		if err := c.Get(ctx, types.NamespacedName{Name: ref.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		tracker := resource.NewProviderConfigUsageTracker(c, &v1beta1.ProviderConfigUsage{})
		if err := tracker.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, c, pc.Spec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		creds := map[string]string{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		cfg := map[string]any{}
		if v := creds["hostname"]; v != "" {
			cfg["hostname"] = v
		}
		// token tiene precedencia
		if v := creds["token"]; v != "" {
			cfg["token"] = v
		} else {
			if v := creds["username"]; v != "" {
				cfg["username"] = v
			}
			if v := creds["password"]; v != "" {
				cfg["password"] = v
			}
		}
		if v := creds["verify_ssl"]; v != "" {
			l := strings.ToLower(strings.TrimSpace(v))
			cfg["verify_ssl"] = !(l == "false" || l == "0" || l == "no")
		}

		ps.Configuration = cfg
		return ps, nil
	}
}
