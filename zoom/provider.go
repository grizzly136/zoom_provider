package zoom

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"jwt": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ZOOM_JWT", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zoom_User_instance": resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"zoom_users": dataSourceUsers(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type tok struct {
	token string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	// 	// req.Header.Add("authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOm51bGwsImlzcyI6ImxOR0pCSGp1Uk9PRktDTTY4TGpIMGciLCJleHAiOjE2MTg5MzgzMjMsImlhdCI6MTYxODg1MTkyNH0.ngd_dOTYMp5ftwP2W-R8XpxHU1dX0i2o6B5xslwLDJ8")
	jwt := "Bearer " + d.Get("jwt").(string)
	c := tok{
		token: jwt,
	}
	os.Setenv("bearer", jwt)
	// 	// c := Client{
	// 	// 	HTTPClient: &http.Client{Timeout: 10 * time.Second},
	// 	// 	// Default Hashicups URL

	// 	// }

	// 	// // Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// 	// req.Header.Add("authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOm51bGwsImlzcyI6ImxOR0pCSGp1Uk9PRktDTTY4TGpIMGciLCJleHAiOjE2MTg5MzgzMjMsImlhdCI6MTYxODg1MTkyNH0.ngd_dOTYMp5ftwP2W-R8XpxHU1dX0i2o6B5xslwLDJ8")

	return c, diags
}
