export function Card({ className = "", ...props }) {
  return <section className={`gs-card ${className}`.trim()} {...props} />;
}
