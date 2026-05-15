# AGPL Compliance Guide

Ech0 is licensed under AGPL-3.0-or-later. This guide explains the practical
requirements for using, modifying, and distributing Ech0.

This document is not legal advice. It is a project-level compliance guide for
operators, downstream developers, and vendors.

## What AGPL Requires

If you modify Ech0 and make it available to users over a network, including as
source code of your modified version under AGPL-3.0-or-later.

You must also preserve copyright notices, license notices, and source-code
access information.

## Common Cases

### Unmodified Self-Hosting

You may run an unmodified Ech0 release for yourself or your organization. Keep
the license and copyright notices intact.

### Modified Self-Hosting or SaaS

If you change Ech0 and let users access that modified service over a network,
you must provide the complete corresponding source code for the deployed
version. This includes backend changes, frontend changes, build scripts, and
other files needed to build and run the modified version.

### Container Images and Binary Distribution

If you distribute Ech0 binaries, container images, packages, or installers, you
must preserve the license notices and provide the corresponding source code for
the distributed build.

### White-Label or Proprietary Forks

White-label deployments and proprietary forks are still subject to the AGPL if
they are based on Ech0. If you do not want to publish your modifications under
AGPL-3.0-or-later, request a commercial license instead.

## Minimum Compliance Checklist

- Keep the AGPL-3.0-or-later license notice.
- Keep Ech0 copyright notices and attribution.
- Provide a clear source-code link to users of the network service.
- If modified, publish the complete corresponding source code for the exact
  deployed version.
- License your modifications under AGPL-3.0-or-later unless you have a separate
  commercial license.
- Do not remove source-code access information from the About page, API
  metadata, CLI output, or documentation.

## Commercial License

For proprietary products, closed-source hosted services, white-label
deployments, or other use cases where you do not want to release your
modifications under AGPL-3.0-or-later, see:

../../COMMERCIAL.md
