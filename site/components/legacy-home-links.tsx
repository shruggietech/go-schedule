'use client';

import { useEffect } from 'react';

const oldHomeAnchors = ['go-schedule-documentation', 'using-go-schedule', 'verifying-an-installation', 'developing-go-schedule', 'reporting-problems'];

export function LegacyHomeLinks() {
  useEffect(() => {
    if (oldHomeAnchors.includes(window.location.hash.slice(1))) window.location.replace(`/docs/${window.location.search}${window.location.hash}`);
  }, []);
  return <>{oldHomeAnchors.map(id => <a key={id} id={id} className="legacy-home-link" href={`/docs/#${id}`}>Continue to this section of the documentation index.</a>)}</>;
}
