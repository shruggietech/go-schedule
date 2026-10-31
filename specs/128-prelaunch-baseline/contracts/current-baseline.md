# S128 Fresh/Current-State Contract

| Surface | Empty state | Current state | Unsupported state |
|---|---|---|---|
| Application DB | One atomic complete schema 24 initialization | Preserve current data/identity and constraints | Actionable explicit reset, no conversion/deletion |
| Preferences | Valid current defaults, independent of Fyne files | Validate/load/patch/reset atomically | Current validation error; explicit full reset repairs |
| Metrics DB | Atomic complete format 4 for either writer | Append/read current records, including combined use | Explicit empty-path/reset instruction without version rewriting |
| Qualification packet | Current revision 2 required observations | Validate current candidate/environment/token/hash/ownership | Reject obsolete packet contract; preserve archive |

Current application bootstrap never performs ALTER/DROP/data-copy upgrades. Preference outputs omit transition fields and migration UI. Explicit reset is a documented development action, not startup behavior or a new cleanup service.

Current installer replacement mechanics and guards that prevent wipe/completion launch during replacement remain. Historical baseline artifact/session requirements do not.

No current contract change rewrites 209 protected artifacts, public release data or prior QA attribution. Old schema/packet input is not silently treated as current acceptance.
