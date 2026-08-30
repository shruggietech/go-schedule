"""Add metadata and bookmarks to the printed brand guide."""

import os
import pikepdf
from pikepdf import Dictionary, Name, OutlineItem, String

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
PDF = os.path.join(ROOT, "brand-guide.pdf")
BOOKMARKS = [
    ("Cover", 0), ("Brand foundation", 1), ("Logo system - mark and lockups", 2),
    ("Logo system - backgrounds and use", 3), ("Color", 4), ("Typography", 5),
    ("Visual language", 6), ("Voice and writing", 7), ("UI system", 8),
    ("Parent relationship and implementation", 9), ("Asset inventory", 10),
]


def main():
    with pikepdf.open(PDF, allow_overwriting_input=True) as pdf:
        with pdf.open_metadata(update_docinfo=False) as meta:
            meta["dc:title"] = "go-schedule Brand System"
            meta["dc:creator"] = ["ShruggieTech"]
            meta["dc:description"] = "Brand and design system for the go-schedule cross-platform task scheduler."
            meta["pdf:Keywords"] = "go-schedule, ShruggieTech, brand system, task scheduler, visual identity"
        pdf.docinfo[Name.Title] = String("go-schedule Brand System")
        pdf.docinfo[Name.Author] = String("ShruggieTech")
        pdf.docinfo[Name.Subject] = String("Brand and design system, version 1.0.0")
        pdf.docinfo[Name.Keywords] = String("go-schedule, ShruggieTech, brand system, task scheduler")
        with pdf.open_outline() as outline:
            outline.root.clear()
            for title, page in BOOKMARKS:
                outline.root.append(OutlineItem(title, page))
        pdf.Root[Name.PageMode] = Name.UseOutlines
        pdf.Root[Name.PageLayout] = Name.SinglePage
        pdf.Root[Name.ViewerPreferences] = Dictionary(DisplayDocTitle=True)
        pdf.save(PDF, linearize=True)
    with pikepdf.open(PDF) as check:
        with check.open_outline() as outline:
            print("pages=%d bookmarks=%d bytes=%d" % (len(check.pages), len(outline.root), os.path.getsize(PDF)))


if __name__ == "__main__":
    main()
