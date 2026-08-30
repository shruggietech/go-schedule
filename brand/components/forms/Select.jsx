export function Select({ label, children, ...props }) {
  return <label className="gs-field"><span>{label}</span><select className="gs-select" {...props}>{children}</select></label>;
}
