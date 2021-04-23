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
			"zoom_user":  data_user(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type tok struct {
	token string
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	jwt := "Bearer " + d.Get("jwt").(string)
	c := tok{
		token: jwt,
	}
	os.Setenv("bearer", jwt)

	var diags diag.Diagnostics

	return c, diags
}
