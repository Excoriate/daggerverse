variable "is_enabled" {
  type        = bool
  description = <<EOF
  Whether this module will be created or not. It is useful, for stack-composite
modules that conditionally includes resources provided by this module..
EOF
  default     = true
}

variable "value_from_another_module" {
  type        = string
  description = "This is a value that will be used in this module."
}
