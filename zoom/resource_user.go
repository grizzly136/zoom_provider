package zoom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	client := &http.Client{Timeout: 10 * time.Second}
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	//getting the data from the resource block
	F_name := d.Get("first_name").(string)
	L_name := d.Get("last_name").(string)
	Email := d.Get("email").(string)
	Type := d.Get("type").(int)
	d.SetId(Email)
	url1 := "https://api.zoom.us/v2/users/" + Email
	//generating req obj
	req1, _ := http.NewRequest("GET", url1, nil)
	//adding jwt token
	jwt := os.Getenv("bearer")
	req1.Header.Add("authorization", jwt)
	//sending the req
	r, _ := client.Do(req1)
	if r.StatusCode == 404 {
		//basic url
		url := "https://api.zoom.us/v2/users"
		//embedding the data into defined structure suitable for sending request
		main_json := Main_body{
			Action: "create",
			User_info: user_info{
				First_name: F_name,
				Last_name:  L_name,
				Email:      Email,
				Type:       Type,
			},
		}
		// converting the struct to json object
		rb, _ := json.Marshal(main_json)
		//generating the request with payload
		req, _ := http.NewRequest("POST", url, strings.NewReader(string(rb)))

		//adding the jwt token to the header of the request
		req.Header.Add("authorization", jwt)
		//adding content type to the header
		req.Header.Add("content-type", "application/json")
		//sending the request
		res, _ := http.DefaultClient.Do(req)
		//setting id to the resorce for reading

		//closing connection
		defer res.Body.Close()
		//calling the reading function for acknowledgement
		resourceUserRead(ctx, d, m)
	} else {

		return diag.Errorf("user alrady created")
	}
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	client := &http.Client{Timeout: 10 * time.Second}
	var diags diag.Diagnostics
	//getting the id from the resorce created by the create function
	UserID := d.Id()
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

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userID := d.Id()

	if d.HasChange("email") || d.HasChange("first_name") || d.HasChange("last_name") || d.HasChange("type") {

		Fir_name := d.Get("first_name").(string)
		Las_name := d.Get("last_name").(string)
		E_mail := d.Get("email").(string)
		Type_1 := d.Get("type").(int)
		//basic url
		url := "https://api.zoom.us/v2/users/" + userID
		//embedding the data into defined structure suitable for sending request

		User_info := user_info{
			First_name: Fir_name,
			Last_name:  Las_name,
			Email:      E_mail,
			Type:       Type_1,
		}

		// converting the struct to json object
		rb, _ := json.Marshal(User_info)
		//generating the request with payload
		req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(rb)))
		//getting the jwt token from the environmental variables
		jwt := os.Getenv("bearer")
		//adding the jwt token to the header of the request
		req.Header.Add("authorization", jwt)
		//adding content type to the header
		req.Header.Add("content-type", "application/json")
		//sending the request
		res, _ := http.DefaultClient.Do(req)

		//closing connection
		defer res.Body.Close()
		d.Set("last_updated", time.Now().Format(time.RFC850))

	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	userID := d.Id()
	url := "https://api.zoom.us/v2/users/" + userID + "?action=disassociate"

	req, _ := http.NewRequest("DELETE", url, nil)

	jwt := os.Getenv("bearer")
	//adding the jwt token to the header of the request
	req.Header.Add("authorization", jwt)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	return diags
}
