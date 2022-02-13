package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() { //nolint:gochecknoinits
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
} //nolint:wsl

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		pvd := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"installer_apt":         resourceApt(),
				"installer_brew":        resourceBrew(),
				"installer_script":      resourceScript(),
				"installer_asdf_plugin": resourceASDFPlugin(),
				"installer_asdf":        resourceASDF(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"installer_apt":  dataSourceApt(),
				"installer_brew": dataSourceBrew(),
			},
		}

		pvd.ConfigureContextFunc = configure(version, pvd)

		return pvd
	}
}

// Add whatever fields, client or connection info, etc. here
// you would need to setup to communicate with the upstream API.
type apiClient struct{}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) { // nolint:lll
	//nolint:wsl
	return func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		// userAgent := p.UserAgent("terraform-provider-scaffolding", version)
		// TODO: myClient.UserAgent = userAgent

		return &apiClient{}, nil
	}
}
