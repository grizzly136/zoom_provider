terraform {
  required_providers {
    zoom = {
      version = "1.0"
      source  = "tharun/edu/zoom"
    }
  }
}

provider "zoom" {

  jwt = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOm51bGwsImlzcyI6ImxOR0pCSGp1Uk9PRktDTTY4TGpIMGciLCJleHAiOjE2MTg5MzgzMjMsImlhdCI6MTYxODg1MTkyNH0.ngd_dOTYMp5ftwP2W-R8XpxHU1dX0i2o6B5xslwLDJ8"

}
//
data "zoom_users" "all" {


}

output "users" {
  value = data.zoom_users.all
}
