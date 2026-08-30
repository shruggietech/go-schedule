export function Button({ variant = "primary", className = "", ...props }) {
  const modifier = variant === "secondary" ? " gs-button--secondary" : "";
  return <button className={`gs-button${modifier} ${className}`.trim()} {...props} />;
}
