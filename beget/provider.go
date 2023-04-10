package beget

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"

	begetOpenapiVps "github.com/LTD-Beget/openapi-vps-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &begetProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &begetProvider{}
}

// begetProvider is the provider implementation.
type begetProvider struct{}

// begetProviderModel maps provider schema data to a Go type.
type begetProviderModel struct {
	Token types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *begetProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "beget"
}

// Schema defines the provider-level schema for configuration data.
func (p *begetProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *begetProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Beget client")

	// Retrieve provider data from configuration
	var config begetProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown beget API Token",
			"The provider cannot create the beget API client as there is an unknown configuration value for the beget API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BEGET_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	token := os.Getenv("BEGET_TOKEN")

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing beget API Token",
			"The provider cannot create the beget API client as there is a missing or empty value for the beget API token. "+
				"Set the token value in the configuration or use the BEGET_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "beget_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "beget_token")

	tflog.Debug(ctx, "Creating Beget client")

	// Create a new beget client using the configuration values
	begetConfig := begetOpenapiVps.NewConfiguration()
	begetConfig.DefaultHeader["Authorization"] = "Bearer " + token
	client := begetOpenapiVps.NewAPIClient(begetConfig)

	//ctx = context.WithValue(context.Background(), begetOpenapiVps.ContextAccessToken, token)

	//begetApiClient := begetApiClient{
	//	client: client,
	//	token:  token,
	//}

	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Unable to Create beget API Client",
	//		"An unexpected error occurred when creating the beget API client. "+
	//			"If the error is not clear, please contact the provider developers.\n\n"+
	//			"beget Client Error: "+err.Error(),
	//	)
	//	return
	//}

	// Make the beget client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Beget client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *begetProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSoftwareDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *begetProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
