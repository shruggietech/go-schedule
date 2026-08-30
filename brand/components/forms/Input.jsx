export function Input({ label, ...props }) {
  return <label className="gs-field"><span>{label}</span><input className="gs-input" {...props} /></label>;
}
