config {
  module = true
  force  = false
}

plugin "terraform" {
  enabled = true
  preset  = "recommended"
}

rule "terraform_module_pinned_source" {
  enabled = true
}

rule "terraform_documented_variables" {
  enabled = true
}
rule "terraform_documented_outputs" {
  enabled = true
}
rule "terraform_unused_required_providers" {
  enabled = true
}

rule "unused_variables" {
  enabled = true
}
