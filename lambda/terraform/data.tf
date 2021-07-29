data "aws_acm_certificate" "primary" {
  domain   = var.domain_name
  statuses = ["ISSUED"]
}
