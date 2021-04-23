package zoom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Users struct {
	Id         string `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
}

type whole_body struct {
	Users []Users `json:"users"`
}

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{

			"users": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	//setting timeout for http request using clinent object in http
	client := &http.Client{Timeout: 10 * time.Second}

	//for errors
	var diags diag.Diagnostics
	//creating request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users", "https://api.zoom.us/v2"), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	//getting jwt token from the environment variables
	jwt := os.Getenv("bearer")
	//adding token to the header
	req.Header.Add("authorization", jwt)
	//sending req
	r, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}

	//closing req at the end
	defer r.Body.Close()
	//initializing structure for copying the json response
	w_b := whole_body{}
	//decoding the body and saving it in w_b
	err = json.NewDecoder(r.Body).Decode(&w_b)
	if err != nil {
		return diag.FromErr(err)
	}
	//creating slice of interfaces for storing user array objects
	u_items := make([]interface{}, len(w_b.Users))
	//mapping the data into uis
	for i, uItem := range w_b.Users {
		ui := make(map[string]interface{})

		ui["id"] = uItem.Id
		ui["first_name"] = uItem.First_name
		ui["last_name"] = uItem.Last_name
		ui["email"] = uItem.Email

		u_items[i] = ui
	}
	// setting the data into users
	if err := d.Set("users", u_items); err != nil {
		return diag.FromErr(err)
	}
	//setting an id
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
