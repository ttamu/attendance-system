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
  name    = tolist(aws_acm_certificate.api_cert.domain_validation_options)[0].resource_record_name
  type    = tolist(aws_acm_certificate.api_cert.domain_validation_options)[0].resource_record_type
  zone_id = data.aws_route53_zone.primary.zone_id
  records = [tolist(aws_acm_certificate.api_cert.domain_validation_options)[0].resource_record_value]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "api_cert_validation" {
  certificate_arn = aws_acm_certificate.api_cert.arn
  validation_record_fqdns = [aws_route53_record.api_cert_validation.fqdn]
}
