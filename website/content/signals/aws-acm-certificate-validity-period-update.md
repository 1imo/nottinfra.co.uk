---
title: AWS – ACM: Public Certificate Validity Period Change
description: AWS Certificate Manager reduces public certificate validity to 198 days, increasing renewal cadence and automation dependency across public endpoints.
keywords: NottInfra, Nottingham Infrastructure, Nottingham IT, Nottingham Platform, Nottingham Security, East Midlands Infrastructure, AWS, ACM, Certificate Manager, TLS, certificate lifecycle, CA Browser Forum
created: 2026-02-20
updated: 2026-02-20
---

Category: Governance  
Impact: High

Amazon Web Services has reduced the maximum validity period for newly issued public certificates in AWS Certificate Manager ( ACM ) from 395 days to 198 days.

The change aligns with the CA / Browser Forum requirement that publicly trusted TLS certificates must not exceed 200 days in validity from 15 March 2026. The mandate applies across publicly trusted certificate authorities and is not AWS-specific.

All newly issued and renewed public certificates now default to 198-day validity. Existing certificates issued under the previous 395-day regime remain valid until renewal or expiry.

ACM continues to renew certificates automatically. Under the revised lifecycle, renewal now occurs approximately 45 days prior to expiry. Pricing for exportable public certificates has been reduced to reflect the shortened validity period.

Exposure:

This change affects estates using ACM-issued public certificates for:

- Application Load Balancers  
- CloudFront distributions  
- API Gateway endpoints  
- Public-facing workloads  
- Hybrid environments exporting certificates externally  

Shorter validity increases renewal frequency across all publicly trusted endpoints. Certificate lifecycle monitoring, CMDB records, automation workflows, and compliance artefacts must reflect the revised cadence.

Estates with manual certificate handling, fragmented automation, or incomplete monitoring will experience elevated expiry risk under the 198-day standard.

Impact:

The operational consequence is governance and lifecycle maturity exposure.

- Renewal cadence increases across public infrastructure  
- Automation dependency becomes more pronounced  
- Monitoring tolerances narrow  
- Manual processes become structurally less sustainable  

There is no service outage associated with this announcement. The impact relates to control-plane governance and operational sustainability rather than availability.

Next Steps:

Infrastructure leaders should:

- Confirm reliance on ACM public certificates  
- Validate renewal automation and monitoring under 198-day validity  
- Review expiry alert thresholds  
- Update lifecycle governance documentation  
- Reassess manual certificate issuance or export workflows  

This revision reinforces that public certificate management is a continuous control-plane function, not an administrative event.

Sources:

Amazon Web Services – AWS Certificate Manager updates default certificate validity to comply with new guidelines \
[https://aws.amazon.com/about-aws/whats-new/2026/02/aws-certificate-manager-updates-default/](https://aws.amazon.com/about-aws/whats-new/2026/02/aws-certificate-manager-updates-default/)

CA / Browser Forum – Baseline Requirements: Maximum Validity Periods for Public TLS Certificates \
[https://cabforum.org/working-groups/server/baseline-requirements/requirements/](https://cabforum.org/working-groups/server/baseline-requirements/requirements/)