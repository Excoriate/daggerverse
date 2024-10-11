locals {
  # Ensure string length is at least 8 characters
  string_length = max(var.string_length, 8)
}
