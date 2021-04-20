terraform {
  required_providers {
    zoom = {
      version = "1.0"
      source  = "tharun/edu/zoom"
    }
  }
}
data "zoom_users" "all" {}

output "users" {
  value = data.zoom_users.all
}
