terraform {
  required_providers {
    zoom = {
      version = "1.0"
      source  = "tharun/edu/zoom"
    }
  }
}
//provider
provider "zoom" {

  jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOm51bGwsImlzcyI6ImxOR0pCSGp1Uk9PRktDTTY4TGpIMGciLCJleHAiOjE2MTk2Nzc0NjAsImlhdCI6MTYxOTA3MjY2MH0.DJGamgGCF7NMKQjzE2QdfEuzkkja-dLk-wy11QUS0eU"
}
//resource

# resource "zoom_User_instance" "tharun" {

#   first_name = "tharun"
#   last_name  = "kumar"
#   email      = "ch.tharunkumar1@gmail.com"
#   type       = 1
# }

# data "zoom_users" "all" {

# }


# data "zoom_user" "single" {
#   email = "ch.tharunkumar1@gmail.com"
# }

# output "resource" {
#   value = zoom_User_instance.tharun
# }

# output "all_users" {
#   value = data.zoom_users.all
# }
# output "single_user" {
#   value = data.zoom_user.single
# }
