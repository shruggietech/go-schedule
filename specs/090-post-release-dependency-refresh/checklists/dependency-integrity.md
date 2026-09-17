# Dependency Integrity Checklist

- [x] Direct selections match the approved S090 baseline.
- [x] Root and desktop Go graphs tidy and verify without residual changes.
- [x] Frontend lockfile restores through `npm ci` without bypass flags.
- [x] Node 26 and Go 1.26 declarations match automation and guidance.
- [x] No test, security, accessibility, packaging, or release assertion is weakened.
- [x] Focused compatibility and canonical verification pass on the final local head.
- [ ] Hosted checks and all review dispositions refer to the final PR head.
- [ ] #240, #241, and #242 contain explicit replacement links before closure.
