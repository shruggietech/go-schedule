"""Shared go-schedule logo geometry and palette."""

NIGHT = "#071014"
PANEL = "#0D171C"
RAISED = "#13232A"
LINE = "#28414B"
TEXT = "#F3F7F8"
MUTED = "#9BAEB6"
PAPER = "#F6F8F7"
INK = "#132027"
INTERVAL = "#62D9B7"
ANCHOR = "#58A6FF"
HOLD = "#F2B84B"
STOP = "#E05F5F"
LIGHT_INTERVAL = "#087A62"
LIGHT_ANCHOR = "#0067C5"
LIGHT_HOLD = "#805500"
LIGHT_STOP = "#B4232A"
WHITE = "#FFFFFF"
BLACK = "#071014"

MARK_CANVAS = 512
MARK_BOX = (78, 102, 434, 410)
CLEAR_SPACE = 32


def rounded_frame(fill=LINE):
    return (
        '<path fill="%s" fill-rule="evenodd" d="'
        'M132 102h248c29.82 0 54 24.18 54 54v200c0 29.82-24.18 54-54 54H132'
        'c-29.82 0-54-24.18-54-54V156c0-29.82 24.18-54 54-54z '
        'M132 122c-18.78 0-34 15.22-34 34v200c0 18.78 15.22 34 34 34h248'
        'c18.78 0 34-15.22 34-34V156c0-18.78-15.22-34-34-34H132z"/>' % fill
    )


def mark_group(mono, interval=INTERVAL, anchor=ANCHOR, line=LINE, small=False):
    """Return portable filled artwork on the native 512-unit grid."""
    if small:
        cells = []
        for i, x in enumerate((122, 184, 246, 308, 370)):
            color = anchor if i == 2 else interval
            cells.append('<rect x="%d" y="326" width="20" height="20" rx="5" fill="%s"/>' % (x, color))
        return (
            rounded_frame(line)
            + '<path fill="%s" d="M116 164l18-18 44 40v16l-44 40-18-18 34-30z"/>' % anchor
            + '<rect x="194" y="238" width="172" height="24" rx="12" fill="%s"/>' % interval
            + ''.join(cells)
        )

    stars = []
    for i, x in enumerate((134, 195, 256, 317, 378)):
        color = anchor if i == 2 else interval
        width = mono.width("*", 46)
        d = mono.path_data("*", 46, x - width / 2, 344)
        stars.append('<path fill="%s" d="%s"/>' % (color, d))
    rails = ''.join(
        '<rect x="%d" y="360" width="32" height="10" rx="5" fill="%s"/>' % (x, line)
        for x in (118, 179, 240, 301, 362)
    )
    return (
        rounded_frame(line)
        + '<path fill="%s" d="M116 164l18-18 44 40v16l-44 40-18-18 34-30z"/>' % anchor
        + '<rect x="194" y="239" width="172" height="22" rx="11" fill="%s"/>' % interval
        + ''.join(stars)
        + rails
    )


def svg(width, height, viewbox, body, title="go-schedule"):
    return (
        '<svg xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" '
        'viewBox="%s" role="img" aria-label="%s">\n  <title>%s</title>\n%s\n</svg>\n'
        % (width, height, viewbox, title, title, body)
    )
