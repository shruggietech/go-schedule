"""Build a portable, fully outlined go-schedule type specimen."""

import os
import sys

import cairosvg

sys.path.insert(0, os.path.dirname(__file__))
from geometry import ANCHOR, HOLD, INTERVAL, LINE, MUTED, NIGHT, TEXT, svg
from typeset import Typesetter

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
OUT = os.path.join(ROOT, "specimens")
FONTS = os.path.join(ROOT, "fonts", "ttf")
DISPLAY = Typesetter(os.path.join(FONTS, "SpaceGrotesk-Bold.ttf"))
BODY = Typesetter(os.path.join(FONTS, "Geist-Regular.ttf"))
MONO = Typesetter(os.path.join(FONTS, "GeistMono-Regular.ttf"))
W, H, X = 1600, 1000, 100


def line(ts, value, size, y, fill, tracking=0):
    return '  <path fill="%s" d="%s"/>\n' % (fill, ts.path_data(value, size, X, y, tracking))


def main():
    os.makedirs(OUT, exist_ok=True)
    body = '  <rect width="%d" height="%d" fill="%s"/>\n' % (W, H, NIGHT)
    body += line(MONO, "GO-SCHEDULE TYPOGRAPHY", 24, 110, INTERVAL, 4)
    body += line(DISPLAY, "Know the next run.", 104, 232, TEXT, -2.6)
    body += line(BODY, "Readable schedules with explicit execution policy.", 30, 298, MUTED)
    body += line(MONO, "SCHEDULE READABILITY", 22, 402, ANCHOR, 2)
    body += line(MONO, "0 15 * * 1-5    @reboot    at 09:30 every weekday", 33, 468, TEXT)
    body += line(MONO, "2026-08-30T12:00:00Z    America/New_York", 33, 526, TEXT)
    body += line(MONO, "GLYPH DISAMBIGUATION", 22, 618, HOLD, 2)
    body += line(MONO, "0 O   1 l I   8 B   5 S   2 Z   { } [ ] ( )", 40, 686, TEXT)
    body += '<line x1="100" y1="752" x2="1500" y2="752" stroke="%s" stroke-width="1"/>\n' % LINE
    body += line(BODY, "Display: Space Grotesk 500 / 700, tracking -0.025em", 26, 832, TEXT)
    body += line(BODY, "Body: Geist 400 / 500, line-height 1.65", 26, 878, TEXT)
    body += line(MONO, "Technical: Geist Mono 400 - schedules, paths, logs, metadata", 26, 924, TEXT)
    source = svg(W, H, "0 0 %d %d" % (W, H), body.rstrip(), "go-schedule type specimen")
    svg_path = os.path.join(OUT, "go-schedule-type-specimen.svg")
    png_path = os.path.join(OUT, "go-schedule-type-specimen.png")
    with open(svg_path, "w", encoding="utf-8", newline="\n") as fh:
        fh.write(source)
    cairosvg.svg2png(url=svg_path, write_to=png_path, output_width=W, output_height=H)
    print("type specimen written")


if __name__ == "__main__":
    main()
