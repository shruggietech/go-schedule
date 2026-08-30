# go-schedule brand kit build

Every distributed asset is generated from the scripts and bundled fonts in this directory.

Run from the kit root:

```text
python build/build_logos.py
python build/build_icons.py
python build/build_specimen.py
python build/build_raster.py
node build/print_pdf.js build/brand-guide.html brand-guide.pdf
python build/finish_pdf.py
python build/verify.py
```

| File | Role |
| --- | --- |
| `geometry.py` | Shared mark geometry and palette |
| `typeset.py` | HarfBuzz shaping and portable outlined SVG text |
| `build_logos.py` | Vector masters, lockups, wordmarks, header, and social preview |
| `build_icons.py` | Starter scheduling icons |
| `build_specimen.py` | Fully outlined typography and schedule specimen |
| `build_raster.py` | PNGs, favicons, ICO, ICNS, and Linux hicolor icons |
| `brand-guide.html` | Print source for the eleven-page PDF |
| `print_pdf.js` | Headless Chromium PDF rendering |
| `finish_pdf.py` | PDF metadata and bookmarks |
| `check_pages.js` | Fixed-page overflow inspection |
| `verify.py` | Contrast, SVG portability, PDF, encoding, dimensions, checksums, and manifest |

Required Python packages: `fonttools`, `uharfbuzz`, `cairosvg`, `pillow`, `pikepdf`, and `numpy`. The PDF step uses Node.js, Playwright, and Chromium.

The verifier fails when an SVG contains live text or a font declaration, a published contrast claim drifts, the PDF loses pages or bookmarks, a page approaches its folio, a required artifact disappears, or a text file contains a BOM or mojibake.
