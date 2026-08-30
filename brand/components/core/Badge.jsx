export function Badge({ state = "ready", className = "", children, ...props }) {
  return <span className={`gs-badge gs-badge--${state} ${className}`.trim()} {...props}>{children}</span>;
}
