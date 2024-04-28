dependencies {
  paths = ["../tg-module"]
}

dependency "upstream_dep" {
  config_path = "../tg-module"
  skip_outputs = true

  mock_outputs = {
    random_string = "temporary-dummy-id"
  }
}

terraform {
  source = "../tf-module-2"
}

inputs =   {
  is_enabled = true
  value_from_another_module = dependency.upstream_dep.outputs.random_string
}
