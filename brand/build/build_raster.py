"""Rasterize SVG masters into web and platform assets."""

import io
import os
import shutil
import struct

import cairosvg
from PIL import Image

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
SVG = os.path.join(ROOT, "logos", "svg")
PNG = os.path.join(ROOT, "logos", "png")
FAV = os.path.join(ROOT, "favicons")
PLATFORM = os.path.join(ROOT, "platform")


def render(src, out, width, height=None, bg=None):
    os.makedirs(os.path.dirname(out), exist_ok=True)
    cairosvg.svg2png(url=src, write_to=out, output_width=width, output_height=height, background_color=bg)
    return Image.open(out).size


def render_image(src, size, bg=None):
    buf = io.BytesIO()
    cairosvg.svg2png(url=src, write_to=buf, output_width=size, output_height=size, background_color=bg)
    buf.seek(0)
    return Image.open(buf).convert("RGBA")


def write_ico(path, entries):
    payloads = []
    for size, image in entries:
        buf = io.BytesIO()
        image.convert("RGB").save(buf, format="ICO", sizes=[(size, size)])
        blob = buf.getvalue()
        offset = struct.unpack("<I", blob[18:22])[0]
        length = struct.unpack("<I", blob[14:18])[0]
        payloads.append((size, blob[offset:offset + length]))
    header = struct.pack("<HHH", 0, 1, len(payloads))
    offset = 6 + 16 * len(payloads)
    directory = body = b""
    for size, data in payloads:
        directory += struct.pack(
            "<BBBBHHII", 0 if size >= 256 else size, 0 if size >= 256 else size,
            0, 0, 1, 32, len(data), offset,
        )
        body += data
        offset += len(data)
    with open(path, "wb") as fh:
        fh.write(header + directory + body)


def main():
    os.makedirs(PNG, exist_ok=True)
    exports = [
        ("go-schedule-horizontal-dark", "go-schedule-horizontal-dark-2400.png", 2400, 640),
        ("go-schedule-horizontal-light", "go-schedule-horizontal-light-2400.png", 2400, 640),
        ("go-schedule-horizontal-white", "go-schedule-horizontal-white-2400.png", 2400, 640),
        ("go-schedule-stacked-dark", "go-schedule-stacked-dark-2048.png", 2048, 2048),
        ("go-schedule-stacked-light", "go-schedule-stacked-light-2048.png", 2048, 2048),
        ("go-schedule-mark-color", "go-schedule-mark-color-1024.png", 1024, 1024),
        ("go-schedule-mark-dark-background", "go-schedule-mark-dark-1024.png", 1024, 1024),
        ("go-schedule-mark-light-background", "go-schedule-mark-light-1024.png", 1024, 1024),
        ("go-schedule-wordmark-interval", "go-schedule-wordmark-interval-1440.png", 1440, None),
        ("go-schedule-wordmark-white", "go-schedule-wordmark-white-1440.png", 1440, None),
        ("go-schedule-wordmark-black", "go-schedule-wordmark-black-1440.png", 1440, None),
        ("go-schedule-social-preview", "go-schedule-social-preview-1280x640.png", 1280, 640),
        ("go-schedule-header", "go-schedule-header-1600x400.png", 1600, 400),
    ]
    for stem, name, width, height in exports:
        out = os.path.join(PNG, name)
        size = render(os.path.join(SVG, stem + ".svg"), out, width, height)
        print("  logos/png/%s %s" % (name, size))

    reduced = os.path.join(FAV, "favicon.svg")
    full = os.path.join(SVG, "go-schedule-mark-dark-background.svg")
    for size in (16, 32, 48, 256):
        src = reduced if size <= 32 else full
        render(src, os.path.join(FAV, "favicon-%dx%d.png" % (size, size)), size, size, "#071014")
    render(full, os.path.join(FAV, "apple-touch-icon.png"), 180, 180, "#071014")
    for size in (192, 512):
        render(full, os.path.join(FAV, "android-chrome-%dx%d.png" % (size, size)), size, size, "#071014")

    ico_sizes = (16, 24, 32, 48, 64, 128, 256)
    ico = os.path.join(FAV, "favicon.ico")
    write_ico(ico, [(s, render_image(reduced if s <= 32 else full, s, "#071014")) for s in ico_sizes])

    windows = os.path.join(PLATFORM, "windows")
    os.makedirs(windows, exist_ok=True)
    shutil.copy2(ico, os.path.join(windows, "go-schedule.ico"))

    linux = os.path.join(PLATFORM, "linux")
    for size in (16, 32, 48, 64, 128, 256, 512):
        out = os.path.join(linux, "hicolor", "%dx%d" % (size, size), "apps", "go-schedule.png")
        render(reduced if size <= 32 else full, out, size, size, "#071014")

    mac = os.path.join(PLATFORM, "macos")
    os.makedirs(mac, exist_ok=True)
    render_image(full, 1024, "#071014").save(os.path.join(mac, "go-schedule.icns"), format="ICNS")
    print("  platform assets written")


if __name__ == "__main__":
    main()
