resource "null_resource" "example_null_resource" {
  triggers = {
    example_name = var.value_from_another_module
  }

  provisioner "local-exec" {
    command = "echo ${var.value_from_another_module}"
  }
}
