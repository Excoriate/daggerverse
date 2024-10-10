locals {
  string_length = var.string_length < 8 ? 8 : var.string_length
}