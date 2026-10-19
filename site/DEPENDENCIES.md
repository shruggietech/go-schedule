# Website dependencies

The website exports static HTML using Next.js 16.4.2, React 19.2.8, Fumadocs core/UI 16.17.6 and Fumadocs MDX 15.4.7. Exact direct versions and the npm lockfile make the selected dependency graph reproducible. The Fumadocs packages were selected together after the older MDX combination failed to compile the current Markdown AST dependencies.

Production direct dependencies Next.js, React, React DOM, Fumadocs core/UI, Mermaid and next-themes use MIT licenses; YAML uses ISC. Local fonts retain their OFL files in the canonical brand asset tree. The build copies those assets verbatim. This project retains its Apache License 2.0 and the included dependency licenses.

Run `npm ci`, `npm audit --audit-level=high`, `npm test` and `npm run build` in this directory. Dependency proposals for `/site` are enabled in Dependabot. Update the Fumadocs core, UI and MDX combination together when compatibility requires it, then rebuild all canonical documents and exercise the production browser suite.
