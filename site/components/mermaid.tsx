'use client';

import { useEffect, useId, useRef, useState } from 'react';

export function Mermaid({ chart }: { chart: string }) {
  const id = useId().replace(/[^a-zA-Z0-9]/g, '');
  const ref = useRef<HTMLDivElement>(null);
  const [failed, setFailed] = useState(false);
  useEffect(() => {
    let cancelled = false;
    import('mermaid').then(async ({ default: mermaid }) => {
      mermaid.initialize({ startOnLoad: false, securityLevel: 'strict', theme: 'neutral' });
      const { svg } = await mermaid.render(`diagram${id}`, chart);
      if (!cancelled && ref.current) ref.current.innerHTML = svg;
    }).catch(() => { if (!cancelled) setFailed(true); });
    return () => { cancelled = true; };
  }, [chart, id]);
  return <figure className="diagram"><div ref={ref} aria-label="Architecture diagram" />{failed && <p>Diagram unavailable. The source is shown below.</p>}<details><summary>Diagram source</summary><pre>{chart}</pre></details></figure>;
}
