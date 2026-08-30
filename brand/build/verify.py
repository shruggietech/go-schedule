"""Verify the complete go-schedule brand kit and write its inventory."""

from __future__ import annotations

import hashlib
import json
import os
import re
import struct
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path

import pikepdf
from PIL import Image

ROOT = Path(__file__).resolve().parent.parent
TOKENS = ROOT / "tokens" / "brand.tokens.json"
PDF = ROOT / "brand-guide.pdf"

REQUIRED = [
    "README.md", "SKILL.md", "brand-guide.pdf", "styles.css",
    "tokens/brand.tokens.json", "tokens/base.css", "tokens/colors.css",
    "tokens/spacing.css", "tokens/typography.css", "fonts/fonts.css",
    "logos/svg/go-schedule-mark-color.svg",
    "logos/svg/go-schedule-mark-reduced.svg",
    "logos/svg/go-schedule-horizontal-dark.svg",
    "logos/svg/go-schedule-horizontal-light.svg",
    "logos/svg/go-schedule-horizontal-white.svg",
    "logos/svg/go-schedule-horizontal-black.svg",
    "logos/svg/go-schedule-stacked-dark.svg",
    "logos/svg/go-schedule-stacked-light.svg",
    "logos/svg/go-schedule-wordmark-interval.svg",
    "logos/svg/go-schedule-wordmark-white.svg",
    "logos/svg/go-schedule-header.svg",
    "logos/svg/go-schedule-social-preview.svg",
    "favicons/favicon.ico", "favicons/site.webmanifest",
    "platform/windows/go-schedule.ico", "platform/macos/go-schedule.icns",
    "platform/linux/go-schedule.desktop", "guidelines/index.html",
    "ui_kits/go-schedule-web/index.html", "components/components.css",
    "icons/calendar.svg", "icons/clock.svg", "icons/play.svg",
]

PNG_SIZES = {
    "favicons/favicon-16x16.png": (16, 16),
    "favicons/favicon-32x32.png": (32, 32),
    "favicons/favicon-48x48.png": (48, 48),
    "favicons/apple-touch-icon.png": (180, 180),
    "favicons/android-chrome-192x192.png": (192, 192),
    "favicons/android-chrome-512x512.png": (512, 512),
    "logos/png/go-schedule-social-preview-1280x640.png": (1280, 640),
}

PAIRS = {
    "text_on_night": ("text", "night"),
    "muted_on_night": ("muted", "night"),
    "interval_on_night": ("interval", "night"),
    "anchor_on_night": ("anchor", "night"),
    "hold_on_night": ("hold", "night"),
    "stop_on_night": ("stop", "night"),
    "ink_on_paper": ("ink", "paper"),
    "light_muted_on_paper": ("light_muted", "paper"),
    "light_interval_on_paper": ("light_interval", "paper"),
    "light_anchor_on_paper": ("light_anchor", "paper"),
    "light_hold_on_paper": ("light_hold", "paper"),
    "light_stop_on_paper": ("light_stop", "paper"),
    "interval_on_paper": ("interval", "paper"),
}


def luminance(hex_color: str) -> float:
    channels = [int(hex_color[i:i + 2], 16) / 255 for i in (1, 3, 5)]
    channels = [c / 12.92 if c <= 0.04045 else ((c + 0.055) / 1.055) ** 2.4 for c in channels]
    return 0.2126 * channels[0] + 0.7152 * channels[1] + 0.0722 * channels[2]


def contrast(a: str, b: str) -> float:
    high, low = sorted((luminance(a), luminance(b)), reverse=True)
    return round((high + 0.05) / (low + 0.05), 2)


def text_files():
    suffixes = {".css", ".desktop", ".html", ".js", ".json", ".jsx", ".md", ".py", ".svg"}
    return [p for p in ROOT.rglob("*") if p.is_file() and p.suffix.lower() in suffixes]


def ico_sizes(path: Path):
    data = path.read_bytes()
    reserved, kind, count = struct.unpack("<HHH", data[:6])
    assert reserved == 0 and kind == 1 and count > 0, f"invalid ICO header: {path}"
    sizes = []
    for i in range(count):
        width, height = data[6 + i * 16], data[7 + i * 16]
        sizes.append((width or 256, height or 256))
    return sizes


def sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def main() -> int:
    failures = []
    notes = []

    for rel in REQUIRED:
        if not (ROOT / rel).is_file():
            failures.append(f"missing required artifact: {rel}")

    tokens = json.loads(TOKENS.read_text(encoding="utf-8"))
    measured = {}
    for claim, (foreground, background) in PAIRS.items():
        actual = contrast(tokens["color"][foreground]["value"], tokens["color"][background]["value"])
        measured[claim] = actual
        if actual != round(float(tokens["contrast"][claim]), 2):
            failures.append(f"contrast drift: {claim} claims {tokens['contrast'][claim]} but measures {actual}")
    if measured["light_interval_on_paper"] < 4.5 or measured["interval_on_paper"] >= 4.5:
        failures.append("light-surface accent policy no longer matches WCAG AA contrast")

    mojibake = re.compile(r"(?:\u00c3.|\u00c2.|\ufffd|\u00e2[\u0080-\u00bf])")
    for path in text_files():
        raw = path.read_bytes()
        rel = path.relative_to(ROOT).as_posix()
        if raw.startswith(b"\xef\xbb\xbf"):
            failures.append(f"UTF-8 BOM found: {rel}")
        try:
            content = raw.decode("utf-8")
        except UnicodeDecodeError as exc:
            failures.append(f"invalid UTF-8: {rel}: {exc}")
            continue
        if mojibake.search(content):
            failures.append(f"possible mojibake: {rel}")

    svg_files = list(ROOT.rglob("*.svg"))
    for path in svg_files:
        content = path.read_text(encoding="utf-8")
        rel = path.relative_to(ROOT).as_posix()
        if re.search(r"<text\b|font-family\s*=|font-family\s*:", content, re.I):
            failures.append(f"non-portable live SVG text: {rel}")
        if not re.search(r"<title\b", content, re.I):
            failures.append(f"SVG lacks accessible title: {rel}")

    white = ROOT / "logos/svg/go-schedule-wordmark-white.svg"
    interval = ROOT / "logos/svg/go-schedule-wordmark-interval.svg"
    if white.exists() and interval.exists() and sha256(white) == sha256(interval):
        failures.append("white and interval wordmarks are identical")

    guide_source = (ROOT / "build/brand-guide.html").read_text(encoding="utf-8")
    cover_match = re.search(r'<div class="page"[^>]*>.*?</div>\s*<div class="page">', guide_source, re.S)
    cover = cover_match.group(0) if cover_match else ""
    if "go-schedule-mark-color.svg" not in cover:
        failures.append("guide cover does not use the transparent mark master")
    if "horizontal" in cover or "stacked" in cover or "wordmark" in cover:
        failures.append("guide cover image repeats the wordmark")

    for rel, expected in PNG_SIZES.items():
        path = ROOT / rel
        if path.exists():
            with Image.open(path) as image:
                if image.size != expected:
                    failures.append(f"wrong PNG dimensions: {rel} is {image.size}, expected {expected}")
                if image.mode not in {"RGB", "RGBA"}:
                    failures.append(f"unexpected PNG mode: {rel} is {image.mode}")

    favicon_sizes = ico_sizes(ROOT / "favicons/favicon.ico") if (ROOT / "favicons/favicon.ico").exists() else []
    windows_sizes = ico_sizes(ROOT / "platform/windows/go-schedule.ico") if (ROOT / "platform/windows/go-schedule.ico").exists() else []
    if not {(16, 16), (32, 32), (48, 48)}.issubset(set(favicon_sizes)):
        failures.append(f"favicon ICO size set incomplete: {favicon_sizes}")
    if not {(16, 16), (32, 32), (48, 48), (256, 256)}.issubset(set(windows_sizes)):
        failures.append(f"Windows ICO size set incomplete: {windows_sizes}")

    pdf_info = {}
    if PDF.exists():
        with pikepdf.open(PDF) as pdf:
            with pdf.open_outline() as outline:
                bookmarks = len(outline.root)
            pdf_info = {
                "pages": len(pdf.pages),
                "bookmarks": bookmarks,
                "title": str(pdf.docinfo.get("/Title", "")),
                "author": str(pdf.docinfo.get("/Author", "")),
            }
            if len(pdf.pages) != 11:
                failures.append(f"brand guide has {len(pdf.pages)} pages, expected 11")
            if bookmarks != 11:
                failures.append(f"brand guide has {bookmarks} bookmarks, expected 11")
            if pdf_info["title"] != "go-schedule Brand System" or pdf_info["author"] != "ShruggieTech":
                failures.append(f"brand guide metadata is incomplete: {pdf_info}")
            if "/StructTreeRoot" not in pdf.Root:
                failures.append("brand guide is not a tagged PDF")
            fonts = {}
            for page in pdf.pages:
                resources = page.get("/Resources", {})
                for _, font_ref in resources.get("/Font", {}).items():
                    font = font_ref
                    name = str(font.get("/BaseFont", "unnamed"))
                    subtype = str(font.get("/Subtype", ""))
                    embedded = False
                    descriptor = font.get("/FontDescriptor")
                    if subtype == "/Type0" and font.get("/DescendantFonts"):
                        descriptor = font["/DescendantFonts"][0].get("/FontDescriptor")
                    if descriptor:
                        embedded = any(key in descriptor for key in ("/FontFile", "/FontFile2", "/FontFile3"))
                    fonts[name] = (subtype, embedded)
            for name, (subtype, embedded) in fonts.items():
                if subtype == "/Type3" or not embedded:
                    failures.append(f"PDF font is unsuitable: {name}, {subtype}, embedded={embedded}")

    node = os.environ.get("NODE_BINARY", "node")
    try:
        check = subprocess.run(
            [node, str(ROOT / "build/check_pages.js"), str(ROOT / "build/brand-guide.html")],
            cwd=ROOT, text=True, capture_output=True, check=False,
        )
        notes.append(check.stdout.strip())
        if check.returncode:
            failures.append("brand guide page-fit check failed\n" + check.stdout.strip() + check.stderr.strip())
    except OSError as exc:
        failures.append(f"could not run page-fit check: {exc}")

    generated = []
    for path in sorted(p for p in ROOT.rglob("*") if p.is_file()):
        rel = path.relative_to(ROOT).as_posix()
        if rel in {"manifest.json", "VERIFY.md"} or "__pycache__" in path.parts:
            continue
        generated.append({"path": rel, "bytes": path.stat().st_size, "sha256": sha256(path)})

    manifest = {
        "name": "go-schedule brand kit",
        "version": tokens["meta"]["version"],
        "generated_at": datetime.now(timezone.utc).isoformat(timespec="seconds"),
        "files": generated,
    }
    (ROOT / "manifest.json").write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8", newline="\n")

    report = [
        "# Verification report", "",
        f"Status: **{'PASS' if not failures else 'FAIL'}**", "",
        f"- Inventory: {len(generated)} files plus this report and manifest",
        f"- Portable SVGs: {len(svg_files)} checked, no live text allowed",
        f"- PDF: {pdf_info.get('pages', 0)} pages, {pdf_info.get('bookmarks', 0)} bookmarks",
        f"- Favicon ICO sizes: {', '.join(f'{w}x{h}' for w, h in favicon_sizes)}",
        f"- Windows ICO sizes: {', '.join(f'{w}x{h}' for w, h in windows_sizes)}",
        "- Encoding: UTF-8 without BOM, mojibake scan enabled", "",
        "## Measured contrast", "",
    ]
    report.extend(f"- `{name}`: {ratio:.2f}:1" for name, ratio in measured.items())
    report += ["", "## Page fit", "", "```text", *(notes or ["Not run"]), "```", ""]
    if failures:
        report += ["## Failures", "", *(f"- {failure}" for failure in failures), ""]
    else:
        report += ["All structural, accessibility, portability, dimension, encoding, and page-fit checks passed.", ""]
    (ROOT / "VERIFY.md").write_text("\n".join(report), encoding="utf-8", newline="\n")

    if failures:
        print("VERIFICATION FAILED")
        for failure in failures:
            print(" -", failure)
        return 1
    print(f"VERIFICATION PASSED: {len(generated)} inventoried files")
    return 0


if __name__ == "__main__":
    sys.exit(main())
