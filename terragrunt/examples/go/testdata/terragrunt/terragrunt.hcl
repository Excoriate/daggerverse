terraform {
  source = "modules/random-string"
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
provider "random" {
    version = "3.5.1"
}
EOF
}

inputs = {
  include_special_characters = true
  include_uppercase          = true
  string_length              = 16
}
