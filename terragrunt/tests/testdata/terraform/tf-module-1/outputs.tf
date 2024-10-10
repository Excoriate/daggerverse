output "random_string_1" {
  description = "The first generated random string"
  value       = random_string.random_string_1.result
}

output "random_string_2" {
  description = "The second generated random string"
  value       = random_string.random_string_2.result
}