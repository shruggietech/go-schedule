"""Generate a starter scheduling icon set."""

import os

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
OUT = os.path.join(ROOT, "icons")

ICONS = {
    "calendar": '<rect x="3" y="4.5" width="18" height="16" rx="2"/><path d="M7 2.5v4M17 2.5v4M3 9h18M7 13h3M14 13h3M7 17h3"/>',
    "clock": '<circle cx="12" cy="12" r="9"/><path d="M12 6.5V12l4 2.5"/>',
    "terminal": '<rect x="2.5" y="4" width="19" height="16" rx="2"/><path d="M6.5 8l3 3-3 3M12 15h5"/>',
    "repeat": '<path d="M4 8h12l-3-3M20 16H8l3 3M16 8l3-3M8 16l-3 3"/>',
    "play": '<path d="M8 5l11 7-11 7z"/>',
    "history": '<path d="M4 7V3M4 7h4M4.7 6.2A9 9 0 1 1 3 14"/><path d="M12 7v5l3 2"/>',
}

TEMPLATE = (
    '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" '
    'fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" '
    'stroke-linejoin="round" role="img" aria-labelledby="title-%s">'
    '<title id="title-%s">%s</title>%s</svg>\n'
)


def main():
    os.makedirs(OUT, exist_ok=True)
    for name, body in ICONS.items():
        with open(os.path.join(OUT, name + ".svg"), "w", encoding="utf-8", newline="\n") as fh:
            label = name.replace("-", " ").title()
            fh.write(TEMPLATE % (name, name, label, body))
    print("%d starter icons written" % len(ICONS))


if __name__ == "__main__":
    main()
