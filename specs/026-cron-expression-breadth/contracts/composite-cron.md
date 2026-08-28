# Contract: Composite Five-Field Cron

## Accepted standard forms

Each field accepts numbers, comma lists, inclusive ascending ranges, wildcard steps, range steps, and the existing case-insensitive month or weekday names.

Representative accepted expressions:

```text
0 9,17 * * *
30 8-17 * * 1-5
*/10 9-17 * * *
10-20/2 * * * *
0 0 1,15 JAN,MAR *
0 12 * * MON,WED,FRI
*/7 * * * *
```

Day-of-month and day-of-week may not both be restricted.

## Description contract

Existing supported shapes retain their concise phrases. Other accepted expressions receive a deterministic field-complete description. A description must mention every restricted field and must never be used as execution input unless it is independently accepted by the human grammar.

## Compilation contract

One accepted expression produces one recurring schedule and one task. The recurrence stores selected minute, hour, month, date, and weekday values independently of retained source text. Evaluation is strictly after the schedule anchor.

## Canonical export contract

Export renders five numeric fields from recurrence values. Full sets become `*`; full-range arithmetic sequences may become `*/n`; consecutive subsets may become `a-b`; all other sets become ordered comma lists. Re-import must reproduce the same occurrence sequence.

## Refusal contract

The following remain named refusals:

- restricted day-of-month together with restricted day-of-week;
- six-field and Quartz-specific syntax;
- `@reboot`;
- modifier lists, ranges, steps, mixtures, or restricted combinations outside the existing focused subsets;
- any export where non-default policy or recurrence state changes what cron would run.

Malformed tokens remain errors naming their field.
