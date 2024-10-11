variable "string_length" {
  description = "The length of the random string to generate"
  type        = number
  default     = 16

  validation {
    condition     = var.string_length > 0
    error_message = "String length must be greater than 0."
  }
}

variable "include_special_characters" {
  description = "Whether to include special characters in the random string"
  type        = bool
  default     = true
}

variable "include_uppercase" {
  description = "Whether to include uppercase letters in the random string"
  type        = bool
  default     = true
}
