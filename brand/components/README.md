# go-schedule components

`components.css` defines a small, framework-neutral UI vocabulary. The JSX files are thin wrappers for projects that already use React.

The signature component is `ScheduleRow`: a task name, its authored schedule expression, and the next observable state. Preserve the expression as technical text and keep state labels explicit.

Load `styles.css` before `components/components.css`. Use the `gs-light` class on a container for light reading surfaces.
