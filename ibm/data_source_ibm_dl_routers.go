package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/networking-go-sdk/directlinkapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	dlCrossConnectRouters = "cross_connect_routers"
	dlRouterName          = "router_name"
	dlTotalConns          = "total_connections"
	dlLocation            = "location_name"
)

func dataSourceIBMDLRouters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMDLRoutersRead,
		Schema: map[string]*schema.Schema{
			dlOfferingType: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Direct Link offering type",
				ValidateFunc: InvokeValidator("ibm_dl_routers", dlOfferingType),
			},
			dlLocation: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Direct Link location",
			},
			dlCrossConnectRouters: {
				Type:        schema.TypeList,
				Description: "Collection of Direct Link cross connect routers",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlRouterName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Router",
						},
						dlTotalConns: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Count of existing Direct Link Dedicated gateways on this router for this account",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMDLRoutersRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	dlType := d.Get(dlOfferingType).(string)
	dlLocName := d.Get(dlLocation).(string)
	listRoutersOptionsModel := &directlinkapisv1.ListOfferingTypeLocationCrossConnectRoutersOptions{}
	listRoutersOptionsModel.OfferingType = &dlType
	listRoutersOptionsModel.LocationName = &dlLocName

	listRouters, detail, err := directLink.ListOfferingTypeLocationCrossConnectRouters(listRoutersOptionsModel)

	if err != nil {
		return fmt.Errorf("Error Getting Direct Link Location Cross Connect Routers: %s\n%s", err, detail)
	}

	routers := make([]map[string]interface{}, 0)
	for _, instance := range listRouters.CrossConnectRouters {
		route := map[string]interface{}{}
		if instance.RouterName != nil {
			route[dlRouterName] = *instance.RouterName
		}
		if instance.TotalConnections != nil {
			route[dlTotalConns] = *instance.TotalConnections
		}
		routers = append(routers, route)
	}
	d.SetId(dataSourceIBMDLRoutersID(d))
	d.Set(dlCrossConnectRouters, routers)
	return nil
}

// dataSourceIBMDLSpeedsID returns a reasonable ID for a direct link speeds list.
func dataSourceIBMDLRoutersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func datasourceIBMDLRoutersValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 2)
	dlTypeAllowedValues := "dedicated"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 dlOfferingType,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              dlTypeAllowedValues})

	ibmDLRoutersDatasourceValidator := ResourceValidator{ResourceName: "ibm_dl_routers", Schema: validateSchema}
	return &ibmDLRoutersDatasourceValidator
}
