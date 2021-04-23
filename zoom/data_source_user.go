package zoom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func data_user() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceSingleUserRead,
		Schema: map[string]*schema.Schema{
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := len(val.(string))
					if v > 128 && v < 1 {
						errs = append(errs, fmt.Errorf("%q max length 64, got: %d", key, v))
					}
					return
				},
			},
		},
	}

}
func resourceSingleUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	client := &http.Client{Timeout: 10 * time.Second}
	var diags diag.Diagnostics
	//getting the id from the resorce created by the create function
	UserID := d.Get("email").(string)
	//url
	url := "https://api.zoom.us/v2/users/" + UserID
	//generating req obj
	req, _ := http.NewRequest("GET", url, nil)
	//adding jwt token
	jwt := os.Getenv("bearer")
	req.Header.Add("authorization", jwt)
	//sending the req
	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	//initalizing defined structure for saving the response data
	User_i := user_info{}
	//decoding the body and adding it into User_i
	json.NewDecoder(r.Body).Decode(&User_i)
	//setting the values to their respective schema varibles
	d.Set("email", User_i.Email)
	d.Set("first_name", User_i.First_name)
	d.Set("last_name", User_i.Last_name)
	d.Set("type", User_i.Type)
	//closing the body
	defer r.Body.Close()
	//setting the id
	d.SetId(User_i.Email)
	return diags
}
