'use client';

import { useState } from 'react';

type Product = { id: string; name: string };
type Cell = { status: string; note: string; source: string; reviewed_at: string };
type Capability = { id: string; category: string; name: string; cells: Record<string, Cell> };

export function ComparisonGrid({ products, capabilities }: { products: Product[]; capabilities: Capability[] }) {
  const [category, setCategory] = useState('all');
  const [competitor, setCompetitor] = useState('all');
  const categories = [...new Set(capabilities.map(row => row.category))];
  const columns = products.filter(product => competitor === 'all' || product.id === 'goschedule' || product.id === competitor);
  const rows = capabilities.filter(row => category === 'all' || row.category === category);
  return <>
    <div className="comparison-filters">
      <label>Capability category<select value={category} onChange={event => setCategory(event.target.value)}><option value="all">All categories</option>{categories.map(value => <option key={value}>{value}</option>)}</select></label>
      <label>Compare with<select value={competitor} onChange={event => setCompetitor(event.target.value)}><option value="all">All products</option>{products.filter(product => product.id !== 'goschedule').map(product => <option key={product.id} value={product.id}>{product.name}</option>)}</select></label>
      <button className="button secondary" type="button" onClick={() => { setCategory('all'); setCompetitor('all'); }}>Reset filters</button>
    </div>
    <p className="fine-print" aria-live="polite">Showing {rows.length} capabilities across {columns.length} products.</p>
    <div className="table-scroll" role="region" aria-label="Capability comparison, scroll horizontally for all products" tabIndex={0}>
      <table className="comparison-table" style={{ minWidth: Math.max(600, columns.length * 210 + 210) }}>
        <caption>Released capability and vendor-documented alternatives. Each cell includes its evidence.</caption>
        <thead><tr><th scope="col">Capability</th>{columns.map(product => <th key={product.id} scope="col">{product.name}</th>)}</tr></thead>
        <tbody>{rows.map(row => <tr key={row.id}><th scope="row">{row.name}<br /><small>{row.category}</small></th>{columns.map(product => { const cell = row.cells[product.id]; return <td key={product.id}><strong>{cell.status}</strong><p>{cell.note}</p><a href={cell.source}>Official source</a><br /><small>Reviewed {cell.reviewed_at}</small></td>; })}</tr>)}</tbody>
      </table>
    </div>
  </>;
}
