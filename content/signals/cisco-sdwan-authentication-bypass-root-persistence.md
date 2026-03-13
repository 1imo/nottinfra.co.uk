---
title: Cisco – SD-WAN: Authentication Bypass Enabling Root Persistence
description: Active exploitation of Cisco SD-WAN authentication bypass enables firmware downgrade, root access, and persistent rogue peer insertion.
keywords: NottInfra, Nottingham Infrastructure, Nottingham Security, Cisco, SD-WAN, CVE-2026-20127, CVE-2022-20775, CISA KEV, NCSC advisory, control plane compromise, downgrade attack, network governance
created: 2026-02-27
updated: 2026-02-27
---

Category: Security\
Impact: Critical

A joint Five Eyes advisory including the UK National Cyber Security Centre confirms active exploitation of CVE-2026-20127 affecting Cisco Catalyst SD-WAN Controller and Manager.

The vulnerability enables remote authentication bypass. Observed attack activity includes insertion of rogue SD-WAN peers and firmware downgrade to exploit CVE-2022-20775, a previously disclosed path traversal vulnerability permitting root privilege escalation.

Following privilege escalation, attackers have been observed restoring original firmware versions to reduce detection likelihood.

Exposure:

Environments running Cisco Catalyst SD-WAN Controller or Manager, particularly those:

- Exposed to untrusted networks  
- Supporting hybrid or multi-site WAN architectures  
- Acting as centralised routing policy enforcement nodes  
- Operating without configuration integrity monitoring  

SD-WAN controllers represent Tier-0 governance infrastructure. Compromise enables manipulation of routing policy, trust boundaries, and segmentation enforcement.

Impact:

This represents control-plane compromise.

Operational consequences include:

- Traffic interception or redirection  
- Rogue peer persistence  
- Lateral movement across network zones  
- Undetected long-term administrative access  

The downgrade-then-restore technique invalidates version-based assurance models. Apparent patch compliance does not confirm configuration integrity.

For regulated or service-critical environments dependent on SD-WAN policy enforcement, this is structural exposure.

Next Steps:

Infrastructure leaders should:

- Upgrade immediately to vendor-fixed software releases.  
- Audit peer configurations for unauthorised additions.  
- Review firmware downgrade/upgrade sequencing logs.  
- Validate control-plane integrity beyond version reporting.  
- Reclassify SD-WAN controllers as Tier-0 monitoring assets if not already treated as such.  

Where SD-WAN underpins business-critical connectivity, assume potential historical exposure until validated.

Sources:

CVE Program – CVE-2026-20127 \
[https://www.cve.org/CVERecord?id=CVE-2026-20127](https://www.cve.org/CVERecord?id=CVE-2026-20127)

CVE Program – CVE-2022-20775 \
[https://www.cve.org/CVERecord?id=CVE-2022-20775](https://www.cve.org/CVERecord?id=CVE-2022-20775)

Cisco Security Advisory – Cisco Catalyst SD-WAN Controller Authentication Bypass Vulnerability
[https://sec.cloudapps.cisco.com/security/center/content/CiscoSecurityAdvisory/cisco-sa-sdwan-authbp-qwCX8D4v](https://sec.cloudapps.cisco.com/security/center/content/CiscoSecurityAdvisory/cisco-sa-sdwan-authbp-qwCX8D4v)

Cybersecurity and Infrastructure Security Agency – Emergency Directive 26-03: Mitigate Vulnerabilities in Cisco SD-WAN Systems \
[https://www.cisa.gov/news-events/directives/ed-26-03-mitigate-vulnerabilities-cisco-sd-wan-systems](https://www.cisa.gov/news-events/directives/ed-26-03-mitigate-vulnerabilities-cisco-sd-wan-systems)