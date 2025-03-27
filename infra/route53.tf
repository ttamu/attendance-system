data "aws_route53_zone" "primary" {
  name         = var.root_domain
  private_zone = false
}

resource "aws_route53_record" "attendance_frontend" {
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = "${var.frontend_subdomain}.${var.root_domain}"
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.frontend_distribution.domain_name
    zone_id                = "Z2FDTNDATAQYW2"
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api_attendance" {
  zone_id = data.aws_route53_zone.primary.zone_id
  name    = "${var.api_subdomain}.${var.root_domain}"
  type    = "A"

  alias {
    name                   = aws_lb.app_alb.dns_name
    zone_id                = aws_lb.app_alb.zone_id
    evaluate_target_health = false
  }
}