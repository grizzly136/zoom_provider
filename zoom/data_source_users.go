package zoom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users", "https://api.zoom.us/v2"), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Add("authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOm51bGwsImlzcyI6ImxOR0pCSGp1Uk9PRktDTTY4TGpIMGciLCJleHAiOjE2MTkwNjg0MTksImlhdCI6MTYxODk4MjAxOX0.6joHYk8c5ROOvkcLy2yCLcaJ9zIbor6b0E-jRvyNd24")

	r, err := client.Do(req)

	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	w_b := whole_body{}
	err = json.NewDecoder(r.Body).Decode(&w_b)
	if err != nil {
		return diag.FromErr(err)
	}
	//usersli:=flatternUsers(&users)

	ois := make([]interface{}, len(w_b.Users))

	for i, uItem := range w_b.Users {
		oi := make(map[string]interface{})

		oi["id"] = uItem.Id
		oi["first_name"] = uItem.First_name
		oi["last_name"] = uItem.Last_name
		oi["email"] = uItem.Email

		ois[i] = oi
	}

	if err := d.Set("users", ois); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
