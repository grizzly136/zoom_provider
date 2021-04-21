terraform {
  required_providers {
    zoom = {
      version = "1.0"
      source  = "tharun/edu/zoom"
    }
  }
}

provider "zoom" {

  jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOm51bGwsImlzcyI6ImxOR0pCSGp1Uk9PRktDTTY4TGpIMGciLCJleHAiOjE2MTkwNjg0MTksImlhdCI6MTYxODk4MjAxOX0.6joHYk8c5ROOvkcLy2yCLcaJ9zIbor6b0E-jRvyNd24"

}
//
data "zoom_users" "all" {


}

output "users" {
  value = data.zoom_users.all
}
