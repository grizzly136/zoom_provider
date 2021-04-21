package zoom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type user_info struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Type       int    `json:"type"`
}
type Main_body struct {
	Action    string    `json:"action"`
	User_info user_info `json:"user_info"`
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := len(val.(string))
					if v > 64 && v < 1 {
						errs = append(errs, fmt.Errorf("%q max length 64, got: %d", key, v))
					}
					return
				},
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := len(val.(string))
					if v > 64 && v < 1 {
						errs = append(errs, fmt.Errorf("%q max length 64, got: %d", key, v))
					}
					return
				},
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
			"type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if !(v == 99 || v == 1 || v == 2 || v == 3) {
						errs = append(errs, fmt.Errorf("%q must be 1 or 2 or 3 or 99, got: %d", key, v))
					}
					return
				},
			},
		},
	}

}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	F_name := d.Get("first_name").(string)
	L_name := d.Get("last_name").(string)
	Email := d.Get("email").(string)
	Type := d.Get("type").(int)
	url := "https://api.zoom.us/v2/users"
	main_json := Main_body{
		Action: "create",
		User_info: user_info{
			First_name: F_name,
			Last_name:  L_name,
			Email:      Email,
			Type:       Type,
		},
	}
	rb, _ := json.Marshal(main_json)

	req, _ := http.NewRequest("POST", url, strings.NewReader(string(rb)))
	jwt := os.Getenv("bearer")
	req.Header.Add("authorization", jwt)
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	d.SetId(main_json.User_info.Email)
	defer res.Body.Close()
	resourceUserRead(ctx, d, m)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	client := &http.Client{Timeout: 10 * time.Second}
	var diags diag.Diagnostics

	UserID := d.Id()
	url := "https://api.zoom.us/v2/users/" + UserID

	req, _ := http.NewRequest("GET", url, nil)
	jwt := os.Getenv("bearer")
	req.Header.Add("authorization", jwt)

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	User_i := user_info{}
	json.NewDecoder(r.Body).Decode(&User_i)
	d.Set("email", User_i.Email)
	d.Set("first_name", User_i.First_name)
	d.Set("last_name", User_i.Last_name)
	d.Set("type", User_i.Type)

	defer r.Body.Close()
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
