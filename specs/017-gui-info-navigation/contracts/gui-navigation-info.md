# Contract: GUI Navigation and Information

## Navigation Collection

The leading application navigation exposes exactly this ordered contract:

| Position | Base label | Content |
|---:|---|---|
| 1 | Tasks | Existing task management view |
| 2 | Groups | Existing group management view |
| 3 | Schedule | Existing calendar view |
| 4 | Activity | Existing activity view; may append a bounded alert count |
| 5 | Info | Local application information view |

Activity label updates MUST preserve positions 4 and 5.

## Info Content

Info MUST contain:

1. the existing full embedded application mark, scaled with its aspect ratio;
2. the exact running build version, visibly prefixed as a version;
3. text identifying ShruggieTech as builder and maintainer;
4. a `ShruggieTech` hyperlink to `https://shruggie.tech`;
5. a `Source repository` hyperlink to `https://github.com/shruggietech/go-schedule`;
6. a `Documentation` hyperlink to `https://shruggietech.github.io/go-schedule/`.

All content is local. Merely constructing or selecting Info MUST NOT call the daemon, access storage, or contact a network resource.

## Presentation Boundary

- Visible text conveys product identity and link purpose without depending on color or imagery alone.
- The mark uses aspect-preserving containment.
- Vertically constrained content remains reachable through the view's layout.
- Activating a hyperlink follows the standard desktop toolkit and operating system behavior; custom navigation, retry, or network handling is out of scope.
