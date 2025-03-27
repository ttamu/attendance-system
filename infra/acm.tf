resource "aws_acm_certificate" "api_cert" {
  domain_name       = "${var.api_subdomain}.${var.root_domain}"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Project = var.project_name
  }
}

resource "aws_route53_record" "api_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.api_cert.domain_validation_options :
    dvo.domain_name => dvo
  }

  allow_overwrite = true
  name            = each.value.resource_record_name
  type            = each.value.resource_record_type
  zone_id         = data.aws_route53_zone.primary.zone_id
  records = [each.value.resource_record_value]
  ttl             = 60
}

resource "aws_acm_certificate_validation" "api_cert_validation" {
  certificate_arn         = aws_acm_certificate.api_cert.arn
  validation_record_fqdns = [
    for record in aws_route53_record.api_cert_validation : record.fqdn
  ]
}

resource "aws_acm_certificate" "frontend_cert" {
  provider          = aws.us_east_1
  domain_name       = "${var.frontend_subdomain}.${var.root_domain}"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Project = var.project_name
  }
}

resource "aws_route53_record" "frontend_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.frontend_cert.domain_validation_options :
    dvo.domain_name => dvo
  }

  allow_overwrite = true
  name            = each.value.resource_record_name
  type            = each.value.resource_record_type
  zone_id         = data.aws_route53_zone.primary.zone_id
  records = [each.value.resource_record_value]
  ttl             = 60
}

resource "aws_acm_certificate_validation" "frontend_cert_validation" {
  provider                = aws.us_east_1
  certificate_arn         = aws_acm_certificate.frontend_cert.arn
  validation_record_fqdns = [
    for record in aws_route53_record.frontend_cert_validation : record.fqdn
  ]
}
