"""Generate every go-schedule SVG asset from one geometry source."""

import math
import os
import sys

sys.path.insert(0, os.path.dirname(__file__))
from geometry import *  # noqa: F401,F403
from typeset import Typesetter

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
SVG = os.path.join(ROOT, "logos", "svg")
FAV = os.path.join(ROOT, "favicons")
FONTS = os.path.join(ROOT, "fonts", "ttf")

DISPLAY = Typesetter(os.path.join(FONTS, "SpaceGrotesk-Bold.ttf"))
BODY = Typesetter(os.path.join(FONTS, "Geist-Regular.ttf"))
MONO = Typesetter(os.path.join(FONTS, "GeistMono-Regular.ttf"))

WM_SIZE = 210
WM_BASELINE = 184
WM_PAD = 8
WM_W = math.ceil(DISPLAY.width("go-schedule", WM_SIZE)) + WM_PAD * 2
WM_H = 226
WM_PATH = DISPLAY.path_data("go-schedule", WM_SIZE, WM_PAD, WM_BASELINE)


def wordmark_group(fill=INTERVAL):
    return '<path class="go-schedule-wordmark" fill="%s" d="%s"/>' % (fill, WM_PATH)


def wordmark_only(fill):
    return svg(WM_W, WM_H, "0 0 %d %d" % (WM_W, WM_H), "  " + wordmark_group(fill))


def mark_only(interval=INTERVAL, anchor=ANCHOR, line=LINE, bg=None, small=False):
    body = ('  <rect width="512" height="512" rx="96" fill="%s"/>\n' % bg) if bg else ""
    body += "  " + mark_group(MONO, interval, anchor, line, small)
    return svg(512, 512, "0 0 512 512", body, "go-schedule mark")


def horizontal(interval, anchor, line, wordmark, bg=None, w=1200, h=320):
    mark_size = 220.0
    gap = 54.0
    wm_scale = 0.63
    wm_w = WM_W * wm_scale
    ink_w = mark_size + gap + wm_w
    x0 = (w - ink_w) / 2
    y_mark = (h - mark_size) / 2
    y_wm = (h - WM_H * wm_scale) / 2
    body = ('  <rect width="%d" height="%d" fill="%s"/>\n' % (w, h, bg)) if bg else ""
    body += '  <g transform="translate(%.2f %.2f) scale(%.6f)">%s</g>\n' % (
        x0, y_mark, mark_size / 512, mark_group(MONO, interval, anchor, line)
    )
    body += '  <g transform="translate(%.2f %.2f) scale(%.6f)">%s</g>' % (
        x0 + mark_size + gap, y_wm, wm_scale, wordmark_group(wordmark)
    )
    return svg(w, h, "0 0 %d %d" % (w, h), body, "go-schedule horizontal lockup")


def stacked(interval, anchor, line, wordmark, bg=None, size=1024):
    mark_size = 420.0
    wm_width = 660.0
    wm_scale = wm_width / WM_W
    gap = 58.0
    ink_h = mark_size + gap + WM_H * wm_scale
    y0 = (size - ink_h) / 2
    body = ('  <rect width="%d" height="%d" fill="%s"/>\n' % (size, size, bg)) if bg else ""
    body += '  <g transform="translate(%.2f %.2f) scale(%.6f)">%s</g>\n' % (
        (size - mark_size) / 2, y0, mark_size / 512, mark_group(MONO, interval, anchor, line)
    )
    body += '  <g transform="translate(%.2f %.2f) scale(%.6f)">%s</g>' % (
        (size - wm_width) / 2, y0 + mark_size + gap, wm_scale, wordmark_group(wordmark)
    )
    return svg(size, size, "0 0 %d %d" % (size, size), body, "go-schedule stacked lockup")


def social_preview(w=1280, h=640):
    mark_size = 300.0
    wm_scale = 0.58
    x_mark, y_mark = 92.0, (h - mark_size) / 2
    x_text = x_mark + mark_size + 86
    y_wm = 206
    tagline = "Readable schedules. Explicit execution policy."
    endorsement = "A SHRUGGIETECH PROJECT"
    body = (
        '  <rect width="%d" height="%d" fill="%s"/>\n'
        '  <defs><radialGradient id="wash" cx="28%%" cy="46%%" r="72%%">'
        '<stop offset="0" stop-color="%s" stop-opacity=".13"/>'
        '<stop offset="1" stop-color="%s" stop-opacity="0"/>'
        '</radialGradient></defs>\n  <rect width="%d" height="%d" fill="url(#wash)"/>\n'
        % (w, h, NIGHT, ANCHOR, NIGHT, w, h)
    )
    body += '  <g transform="translate(%.2f %.2f) scale(%.6f)">%s</g>\n' % (
        x_mark, y_mark, mark_size / 512, mark_group(MONO)
    )
    body += '  <g transform="translate(%.2f %.2f) scale(%.6f)">%s</g>\n' % (
        x_text, y_wm, wm_scale, wordmark_group(TEXT)
    )
    body += '  <path fill="%s" d="%s"/>\n' % (
        INTERVAL, BODY.path_data(tagline, 28, x_text, 408)
    )
    body += '  <path fill="%s" d="%s"/>' % (
        MUTED, MONO.path_data(endorsement, 18, x_text, 464, 2.3)
    )
    return svg(w, h, "0 0 %d %d" % (w, h), body, "go-schedule social preview")


def header(w=1600, h=400):
    body = (
        '  <rect width="%d" height="%d" fill="%s"/>\n'
        '  <defs><linearGradient id="header-wash" x1="0" x2="1">'
        '<stop stop-color="%s" stop-opacity=".12"/>'
        '<stop offset="1" stop-color="%s" stop-opacity="0"/>'
        '</linearGradient></defs>\n  <rect width="%d" height="%d" fill="url(#header-wash)"/>\n'
        % (w, h, NIGHT, INTERVAL, NIGHT, w, h)
    )
    mark_size = 220.0
    wm_scale = 0.63
    gap = 54.0
    wm_w = WM_W * wm_scale
    ink_w = mark_size + gap + wm_w
    x0 = (w - ink_w) / 2
    body += '  <g transform="translate(%.2f 90) scale(%.6f)">%s</g>\n' % (
        x0, mark_size / 512, mark_group(MONO)
    )
    body += '  <g transform="translate(%.2f %.2f) scale(%.6f)">%s</g>' % (
        x0 + mark_size + gap, (h - WM_H * wm_scale) / 2, wm_scale, wordmark_group(TEXT)
    )
    return svg(w, h, "0 0 %d %d" % (w, h), body, "go-schedule repository header")


def write(path, content):
    with open(path, "w", encoding="utf-8", newline="\n") as fh:
        fh.write(content)
    return path


def main():
    os.makedirs(SVG, exist_ok=True)
    os.makedirs(FAV, exist_ok=True)
    files = []
    files += [
        write(os.path.join(SVG, "go-schedule-horizontal-dark.svg"), horizontal(INTERVAL, ANCHOR, LINE, TEXT, NIGHT)),
        write(os.path.join(SVG, "go-schedule-horizontal-light.svg"), horizontal(LIGHT_INTERVAL, LIGHT_ANCHOR, LINE, INK, PAPER)),
        write(os.path.join(SVG, "go-schedule-horizontal-color.svg"), horizontal(INTERVAL, ANCHOR, LINE, TEXT)),
        write(os.path.join(SVG, "go-schedule-horizontal-white.svg"), horizontal(WHITE, WHITE, WHITE, WHITE)),
        write(os.path.join(SVG, "go-schedule-horizontal-black.svg"), horizontal(BLACK, BLACK, BLACK, BLACK)),
        write(os.path.join(SVG, "go-schedule-stacked-dark.svg"), stacked(INTERVAL, ANCHOR, LINE, TEXT, NIGHT)),
        write(os.path.join(SVG, "go-schedule-stacked-light.svg"), stacked(LIGHT_INTERVAL, LIGHT_ANCHOR, LINE, INK, PAPER)),
        write(os.path.join(SVG, "go-schedule-stacked-white.svg"), stacked(WHITE, WHITE, WHITE, WHITE)),
        write(os.path.join(SVG, "go-schedule-stacked-black.svg"), stacked(BLACK, BLACK, BLACK, BLACK)),
        write(os.path.join(SVG, "go-schedule-wordmark-interval.svg"), wordmark_only(INTERVAL)),
        write(os.path.join(SVG, "go-schedule-wordmark-white.svg"), wordmark_only(WHITE)),
        write(os.path.join(SVG, "go-schedule-wordmark-black.svg"), wordmark_only(BLACK)),
        write(os.path.join(SVG, "go-schedule-mark-color.svg"), mark_only()),
        write(os.path.join(SVG, "go-schedule-mark-white.svg"), mark_only(WHITE, WHITE, WHITE)),
        write(os.path.join(SVG, "go-schedule-mark-black.svg"), mark_only(BLACK, BLACK, BLACK)),
        write(os.path.join(SVG, "go-schedule-mark-dark-background.svg"), mark_only(bg=NIGHT)),
        write(os.path.join(SVG, "go-schedule-mark-light-background.svg"), mark_only(LIGHT_INTERVAL, LIGHT_ANCHOR, LINE, PAPER)),
        write(os.path.join(SVG, "go-schedule-mark-reduced.svg"), mark_only(small=True)),
        write(os.path.join(SVG, "go-schedule-social-preview.svg"), social_preview()),
        write(os.path.join(SVG, "go-schedule-header.svg"), header()),
        write(os.path.join(FAV, "favicon.svg"), mark_only(bg=NIGHT, small=True)),
    ]
    for path in files:
        print("  ", os.path.relpath(path, ROOT), os.path.getsize(path), "bytes")
    print("%d SVG files written" % len(files))


if __name__ == "__main__":
    main()
