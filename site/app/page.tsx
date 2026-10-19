import Link from 'next/link';
import { Header, Footer } from '@/components/chrome';
import { LegacyHomeLinks } from '@/components/legacy-home-links';

export const metadata = { alternates: { canonical: '/' } };

export default function Home() {
  return <><Header /><main id="main-content" tabIndex={-1}>
    <LegacyHomeLinks /><section className="hero">
      <div className="hero-copy"><p className="release-line"><span className="status-dot" /> Open-source workload automation</p><h1>Put your workflows <br />on schedule.</h1><p className="hero-description">From the first scheduled script to connected, event-driven workflows. Write when it should run. See what happened. Keep work moving.</p><div className="actions"><Link className="button primary" href="/docs/install/">Get started</Link><a className="button secondary" href="https://github.com/shruggietech/go-schedule/releases/latest">Download the latest release</a></div><p className="platforms">Windows, Linux and macOS <span aria-hidden="true">/</span> CLI, desktop and API</p></div>
      <div className="workflow-panel" aria-label="Example scheduled workflow">
        <div className="panel-title"><span className="status-dot" /> The daily reporting workflow <span>09:00</span></div>
        <div className="workflow-step"><span className="step-icon">↳</span><div><strong>Collect the numbers</strong><p>Weekdays at 09:00</p></div><span className="result success">Success</span></div>
        <div className="connector"><span />On completion</div>
        <div className="workflow-step"><span className="step-icon">↳</span><div><strong>Build the report</strong><p>Run after collection succeeds</p></div><span className="result success">Success</span></div>
        <div className="connector"><span />On completion</div>
        <div className="workflow-step"><span className="step-icon">↳</span><div><strong>Deliver the result</strong><p>Notify your team through a webhook</p></div><span className="result">Delivered</span></div>
        <p className="example-caption">An example of schedules, completion chains and notifications working together.</p>
      </div>
    </section>
    <section className="command-section"><div><h2>Say when. <br />Or bring your cron.</h2><p>Readable schedules and supported cron share one scheduling engine. Choose a timezone, review the next runs, and control what happens when work overlaps or a date is missed.</p><Link className="text-link" href="/docs/cron/">Explore scheduling syntax</Link></div><div className="terminal"><div className="terminal-bar"><span />gosched</div><pre><code><span className="comment"># Schedule a weekday report</span>{'\n'}gosched task add weekday-report \{'\n'}  --command /usr/local/bin/report.sh \{'\n'}  --schedule <span className="code-string">"weekdays at 09:00"</span> \{'\n'}  --tz America/New_York{'\n\n'}<span className="comment"># Prefer cron? Use "0 9 * * 1-5".</span>{'\n'}gosched runs --task &lt;task-id&gt;</code></pre><div className="terminal-note">The CLI reports the interpreted schedule and upcoming runs.</div></div></section>
    <section className="capabilities"><div className="section-intro"><h2>More ways to start. <br />One place to understand.</h2><p>Use go-schedule for recurring operations, incoming files, build hooks and dependent jobs. Keep the execution history and the next step within reach.</p></div><div className="capability-list">
      <article><span className="feature-symbol" aria-hidden="true">⌁</span><div><h3>Time and events, connected</h3><p>Run on readable schedules, cron, daemon startup, file changes, external triggers or task completion. Group related tasks and connect their outcomes.</p><Link href="/docs/cli/">Explore commands</Link></div></article>
      <article><span className="feature-symbol" aria-hidden="true">▤</span><div><h3>A control center for the work</h3><p>Author tasks, inspect the calendar, follow activity and manage explicitly enrolled remote daemons from the desktop. The daemon keeps scheduling when the window closes.</p><Link href="/docs/gui-fields/">Meet the desktop</Link></div></article>
      <article><span className="feature-symbol" aria-hidden="true">⌘</span><div><h3>Automation you can build on</h3><p>Use the CLI, JSON API or optional MCP access. Inspect permissions and audit records, move portable definitions and connect your existing tools.</p><Link href="/docs/api/">Read the API reference</Link></div></article>
    </div></section>
    <section className="comparison-invitation"><div><h2>A serious place in your automation stack.</h2><p>See how go-schedule compares with established scheduling and workflow platforms, where it fits today, and what comes next.</p></div><Link className="button secondary" href="/compare/">Compare capabilities</Link></section>
    <section className="closing"><img src="/assets/brand/go-schedule-mark-color.svg" width="80" height="80" alt="" /><h2>Your scripts already do the work. <br />Give them a schedule.</h2><Link className="button primary" href="/docs/install/">Install go-schedule</Link><p>Run locally. Connect other daemons when you need them.</p></section>
  </main><Footer /></>;
}
