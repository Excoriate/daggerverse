variable "string_length" {
  description = "The length of the random strings to generate"
  type        = number
  default     = 16
}

variable "include_special_characters" {
  description = "Whether to include special characters in the random strings"
  type        = bool
  default     = true
}

variable "include_uppercase" {
  description = "Whether to include uppercase letters in the random strings"
  type        = bool
  default     = true
}
