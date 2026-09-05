# Contract: S047 Desktop Release Evidence

## Canonical scenario extension

The S040 evidence schema remains version 1 and retains all 36 existing scenario identities. S047 appends exactly these eleven required identities:

```text
desktop.appearance-standard
desktop.appearance-scaled
desktop.interaction-states
desktop.interaction-states-scaled
desktop.navigation-options
desktop.navigation-options-scaled
desktop.scroll-input
desktop.tasks-table
desktop.tasks-table-scaled
desktop.schedule-activity-tables
desktop.schedule-activity-tables-scaled
```

Omission, duplication, an unknown scenario, a non-passing status, an invalid environment, an absent screenshot, or an invalid scenario metric returns release gate exit code 1 with an actionable diagnostic.

## Shared native contract

Every `desktop.` observation:

- uses a Windows 11 client environment;
- uses the intended user at medium integrity;
- observes the installed LocalSystem service;
- falls within the evidence interval;
- references at least one hashed attachment whose bytes decode as a supported raster image (declared media type and extension are not trusted);
- identifies each exercised palette, DPI, size, state, font, input, header, or semantic value through an exact normalized set where the data model requires one; and
- records reviewed outcome booleans instead of relying on summary prose.

Comma-separated exact sets are order-insensitive, whitespace-trimmed, and reject blank or duplicate members.

## Scenario contracts

### Appearance

The standard and scaled observations jointly close the native evidence gap for
#101 and the font portion of #106. Both require Dark/Light, all five configured
font families, System default/reset, persistence, sharp body and Info text, centered/unclipped labels, resize, minimize/restore, and reopen. Standard uses 96 DPI; scaled uses an environment and metric greater than 96 DPI.

### Interaction states

The standard and scaled interaction observations support #109 and consumers
#104, #105, #112, and #113. Each requires every shared control family and
state, 4.5:1 normal-text and 3:1 non-text minima, readable labels/glyphs, visible focus, persistent selected identity, and non-color cues.

### Navigation and Options

The standard and scaled navigation/options observations support #104, #105, and #106. Each covers the complete destination order, both supported content sizes, rail spacing/boundary, bottom-right Exit behavior, compact storage rows, unavailable rows, exact Copy, selector alternatives, and absence of horizontal scrolling.

### Scroll input

The scroll observation supports #111 and #106. A conventional wheel must pass at 1x, 2x, and 4x on every owned scroll surface. Touchpad behavior passes when available; otherwise a reason is mandatory. Persistence, immediate application, nested delta isolation, and keyboard scrolling remain required.

### Tables

The standard and scaled Tasks and Schedule/Activity observations support #112 and #113, consuming #109's shared state rules. All require at least 100 populated rows per view, both palettes, both sizes, exact headers, frozen-header and disclosure behavior, no horizontal scrolling, odd/even/hover/focus/selection states, and live-refresh identity. View-specific operations and semantic sets follow `data-model.md`.

## Collector contract

`Initialize` creates 47 unavailable observations and scenario templates for all `setup.`, `remove.`, and `desktop.` identities. A desktop template includes its complete metrics shape and a default screenshot path. `RecordObservation` retains the existing one-write import behavior. `Finalize` hashes attachments and invokes the same Go validator used by promotion.

`CaptureWindow` remains limited to `window.*`; desktop observations reuse the environment and screenshot conventions but require operator-reviewed fragments because native visual and physical-input outcomes cannot be derived from HWND geometry.

## Issue traceability contract

The runbook and S047 verification record include a table mapping every included issue to its required formal observations. A passing local demo is recorded as exploratory evidence only. Closing keywords are used later only for an issue whose individual criteria are actually satisfied; all other references remain non-closing.

## Compatibility

- Existing candidate identity, archive safety, task-run attachment, window, access, error, setup, and removal validation remains unchanged.
- Older 36-observation bundles intentionally fail the updated gate because they do not prove the current v1 desktop baseline.
- The synthetic fixture uses `automated-fixture`; production validation still requires `attended-windows`.
- No product preference, database, API, or command schema changes.
