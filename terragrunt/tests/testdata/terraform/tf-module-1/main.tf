provider "random" {}

resource "random_string" "random_string_1" {
  length  = local.string_length
  special = var.include_special_characters
  upper   = var.include_uppercase
}

resource "random_string" "random_string_2" {
  length  = local.string_length
  special = var.include_special_characters
  upper   = var.include_uppercase
}