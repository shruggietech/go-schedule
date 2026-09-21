# Quickstart: Verify Agent Grant Controls

1. Start the local daemon and desktop with all optional listeners off.
2. Open Agent Access and verify the summary reports MCP Off, stdio Available on demand, Localhost HTTP Off, and Remote HTTPS Off or Not configured independently.
3. Create an Observe grant for 1 hour. Verify the native clipboard receives one enrollment bundle and no phrase appears in the interface, console, logs, snapshots, or audit output.
4. Exchange the pairing as an MCP client and refresh Agent Access. Verify the row shows client, this daemon, Observe, Remote HTTPS, creation, last use when available, the fixed deadline, and Active.
5. Create Operate and Manage grants with different durations. Verify Manage uses definition-changing language and is not distinguished by color alone.
6. Narrow Manage to Operate, then Observe. Keep a connection open and verify a now-disallowed tool fails on its next request.
7. Impose or shorten an expiry and verify expiry stops the existing connection without restart.
8. Revoke the grant through the confirmation dialog and verify the next request fails while recent audit evidence remains available.
9. Exercise the complete workflow by keyboard and an accessibility query, including dialog focus return and announcements.
10. Run focused Go and frontend tests, race tests, MCP conformance and hostile-content suites, formatting, and `sh scripts/verify.sh all`.
