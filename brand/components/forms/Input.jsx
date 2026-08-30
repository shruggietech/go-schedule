export function Input({ label, className = "", ...props }) {
  return <label className="gs-field"><span>{label}</span><input className={`gs-input ${className}`.trim()} {...props} /></label>;
}
