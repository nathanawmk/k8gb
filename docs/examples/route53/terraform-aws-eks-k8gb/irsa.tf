module "iam_assumable_role_admin" {
  source                        = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"
  version                       = "2.14.0"
  create_role                   = true
  role_name                     = "external-dns-${var.cluster_name}"
  provider_url                  = replace(module.eks-cluster.cluster_oidc_issuer_url, "https://", "")
  role_policy_arns              = [aws_iam_policy.external-dns-route53.arn]
  oidc_fully_qualified_subjects = ["system:serviceaccount:k8gb:external-dns"]
}

resource "aws_iam_policy" "external-dns-route53" {
  name        = "AllowExternalDNSUpdatesFor${var.cluster_name}"
  description = "Enable external-dns to update Route53"
  policy      = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "route53:ChangeResourceRecordSets"
            ],
            "Resource": [
                "arn:aws:route53:::hostedzone/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "route53:ListHostedZones",
                "route53:ListResourceRecordSets"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
EOF
  tags = {
    yor_trace = "0e535bb3-2e22-42cf-bce3-c2f30e983c00"
  }
}

