terraform {
  required_providers {
    astronomer = {
      source = "registry.terraform.io/gk-consulting/astronomer"
    }
  }
}

provider "astronomer" {
  token           = "your token here"
  organization_id = "cabcabcabcabcabcabcabac"
}
