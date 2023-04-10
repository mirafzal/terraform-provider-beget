package beget

import (
	"context"

	begetOpenapiVps "github.com/LTD-Beget/openapi-vps-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &softwareDataSource{}
	_ datasource.DataSourceWithConfigure = &softwareDataSource{}
)

func NewSoftwareDataSource() datasource.DataSource {
	return &softwareDataSource{}
}

type softwareDataSource struct {
	client *begetOpenapiVps.APIClient
}

// softwareDataSourceModel maps the data source schema data.
type softwareDataSourceModel struct {
	Software []softwareModel `tfsdk:"software"`
}

// softwareModel maps software schema data.
type softwareModel struct {
	ID                types.Int64             `tfsdk:"id"`
	Name              types.String            `tfsdk:"name"`
	DisplayName       types.String            `tfsdk:"display_name"`
	Description       types.String            `tfsdk:"description"`
	DescriptionEn     types.String            `tfsdk:"description_en"`
	Slug              types.String            `tfsdk:"slug"`
	DocumentationSlug types.String            `tfsdk:"documentation_slug"`
	Category          []softwareCategoryModel `tfsdk:"category"`
}

// softwareCategoryModel maps software ingredients data
type softwareCategoryModel struct {
	SysName types.String `tfsdk:"sys_name"`
	Name    types.String `tfsdk:"name"`
	NameEn  types.String `tfsdk:"name_en"`
	IsMain  types.Bool   `tfsdk:"is_main"`
}

func (d *softwareDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*begetOpenapiVps.APIClient)
}

func (d *softwareDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_software"
}

// Schema defines the schema for the data source.
func (d *softwareDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"software": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"display_name": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"description_en": schema.StringAttribute{
							Computed: true,
						},
						"slug": schema.StringAttribute{
							Computed: true,
						},
						"documentation_slug": schema.StringAttribute{
							Computed: true,
						},
						"category": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"sys_name": schema.StringAttribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
									"name_en": schema.StringAttribute{
										Computed: true,
									},
									"is_main": schema.BoolAttribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *softwareDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state softwareDataSourceModel

	software, _, err := d.client.MarketplaceServiceApi.MarketplaceServiceGetSoftwareList(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Beget Software",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, softwareItem := range software.GetSoftware() {
		softwareState := softwareModel{
			ID:                types.Int64Value(int64(softwareItem.GetId())),
			Name:              types.StringValue(softwareItem.GetName()),
			DisplayName:       types.StringValue(softwareItem.GetDisplayName()),
			Description:       types.StringValue(softwareItem.GetDescription()),
			DescriptionEn:     types.StringValue(softwareItem.GetDescriptionEn()),
			Slug:              types.StringValue(softwareItem.GetSlug()),
			DocumentationSlug: types.StringValue(softwareItem.GetDocumentationSlug()),
		}

		for _, category := range softwareItem.GetCategory() {
			softwareState.Category = append(softwareState.Category, softwareCategoryModel{
				SysName: types.StringValue(category.GetSysName()),
				Name:    types.StringValue(category.GetName()),
				NameEn:  types.StringValue(category.GetNameEn()),
				IsMain:  types.BoolValue(category.GetIsMain()),
			})
		}

		state.Software = append(state.Software, softwareState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
