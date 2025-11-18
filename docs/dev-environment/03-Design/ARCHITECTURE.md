# Architecture - spec-05
## Docker Compose Profiles
```yaml
services:
  postgres:
    profiles: [full, db-only]
  api-mobile:
    profiles: [full, mobile-only]
  api-admin:
    profiles: [full, admin-only]
  worker:
    profiles: [full, worker-only]
```
