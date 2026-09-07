---
name: aof-secure-sdlc
description: Apply AOF-Secure-SDLC domain-profile controls to software requirements, architecture, coding, testing, security verification, release and deployment.
---
# AOF Secure SDLC

## Trigger
Use for requirements, architecture/security review, threat modeling, code/test/security Verification, release/deployment Authority gates, Human control gates, or security Evidence.

## Procedure
1. Read `aof/secure-sdlc/profile.yaml` and `.aof/requirements.yaml`.
2. Bind software requirements to acceptance criteria.
3. Identify architecture/security review and threat-model needs.
4. Define code/test/security Verification and required Evidence.
5. Keep release/deployment Authority independent from technical deployment Capability.
6. Add Human control gates according to Risk.
7. Use SAST/SCA/DAST or equivalent Evidence when applicable.
8. Preserve all mandatory base-profile controls.

AOF-Secure-SDLC is a domain profile, not a linear maturity level and not permission to weaken base governance.
