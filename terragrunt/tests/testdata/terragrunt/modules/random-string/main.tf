resource "random_string" "this" {
  length  = local.string_length
  special = var.include_special_characters
  upper   = var.include_uppercase
}