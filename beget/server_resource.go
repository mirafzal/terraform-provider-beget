package beget

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"

	begetOpenapiVps "github.com/LTD-Beget/openapi-vps-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &serverResource{}
	_ resource.ResourceWithConfigure = &serverResource{}
)

// NewServerResource is a helper function to simplify the provider implementation.
func NewServerResource() resource.Resource {
	return &serverResource{}
}

// serverResource is the resource implementation.
type serverResource struct {
	client *begetOpenapiVps.APIClient
}

// serverResourceModel maps the resource schema data.
type serverResourceModel struct {
	ID            types.String             `tfsdk:"id"`
	LastUpdated   types.String             `tfsdk:"last_updated"`
	Slug          types.String             `tfsdk:"slug"`
	DisplayName   types.String             `tfsdk:"display_name"`
	Hostname      types.String             `tfsdk:"hostname"`
	Configuration serverConfigurationModel `tfsdk:"configuration"`
	Status        types.String             `tfsdk:"status"`
	SshKeys       []serverSshKeyModel      `tfsdk:"ssh_keys"`
	HasPassword   types.String             `tfsdk:"has_password"`
}

// serverConfigurationModel maps server config data.
type serverConfigurationModel struct {
	ID           types.String  `tfsdk:"id"`
	Name         types.String  `tfsdk:"name"`
	CpuCount     types.Int64   `tfsdk:"cpu_count"`
	DiskSize     types.Int64   `tfsdk:"disk_size"`
	Memory       types.Int64   `tfsdk:"memory"`
	PriceDay     types.Float64 `tfsdk:"price_day"`
	PriceMonth   types.Float64 `tfsdk:"price_month"`
	Available    types.Bool    `tfsdk:"available"`
	Custom       types.Bool    `tfsdk:"custom"`
	Configurable types.Bool    `tfsdk:"configurable"`
}

// serverConfigurationModel maps server ssh keys data.
type serverSshKeyModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Fingerprint types.String `tfsdk:"fingerprint"`
}

// Configure adds the provider configured client to the resource.
func (r *serverResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*begetOpenapiVps.APIClient)
}

// Metadata returns the resource type name.
func (r *serverResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

// Schema defines the schema for the resource.
func (r *serverResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"slug": schema.StringAttribute{
				Computed: true,
			},
			"display_name": schema.StringAttribute{
				Computed: true,
			},
			"hostname": schema.StringAttribute{
				Computed: true,
			},
			"configuration": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
					"cpu_count": schema.Int64Attribute{
						Computed: true,
					},
					"disk_size": schema.Int64Attribute{
						Computed: true,
					},
					"memory": schema.Int64Attribute{
						Computed: true,
					},
					"price_day": schema.Float64Attribute{
						Computed: true,
					},
					"price_month": schema.Float64Attribute{
						Computed: true,
					},
					"available": schema.BoolAttribute{
						Computed: true,
					},
					"custom": schema.BoolAttribute{
						Computed: true,
					},
					"configurable": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"ssh_keys": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"fingerprint": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"has_password": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

// Create a new resource
func (r *serverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan serverResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var items []hashicups.OrderItem
	for _, item := range plan.Items {
		items = append(items, hashicups.OrderItem{
			Coffee: hashicups.Coffee{
				ID: int(item.Coffee.ID.ValueInt64()),
			},
			Quantity: int(item.Quantity.ValueInt64()),
		})
	}

	// Create new order
	server, err := r.client.CreateOrder(items)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue()
	for orderItemIndex, orderItem := range order.Items {
		plan.Items[orderItemIndex] = orderItemModel{
			Coffee: orderItemCoffeeModel{
				ID:          types.Int64Value(int64(orderItem.Coffee.ID)),
				Name:        types.StringValue(orderItem.Coffee.Name),
				Teaser:      types.StringValue(orderItem.Coffee.Teaser),
				Description: types.StringValue(orderItem.Coffee.Description),
				Price:       types.Float64Value(orderItem.Coffee.Price),
				Image:       types.StringValue(orderItem.Coffee.Image),
			},
			Quantity: types.Int64Value(int64(orderItem.Quantity)),
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
