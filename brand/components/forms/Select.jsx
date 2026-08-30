export function Select({ label, children, className = "", ...props }) {
  return <label className="gs-field"><span>{label}</span><select className={`gs-select ${className}`.trim()} {...props}>{children}</select></label>;
}
