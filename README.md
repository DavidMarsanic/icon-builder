# Icon Builder

A graphical picker for composing app icons: pick a background icon, type
up to 3 letters, watch it render live. Opens as its own window.

## Features

- Browse (and search) a bundled set of icons, or start from none for a
  plain gradient+letters badge.
- Type an app name to seed a deterministic gradient and auto-suggest
  letters, or override both by hand.
- Live preview, updated on every change.
- Copy the composed result as a self-contained SVG.

## Requirements

**A Chromium-based browser already installed**: Google Chrome, Chromium,
Brave, Microsoft Edge, or Arc — renders the app's own UI window.

## How it works

All the actual composition logic — the gradient palette, the deterministic
hashing, the letter-overlay math — lives in
[icon-composer](https://github.com/DavidMarsanic/icon-composer), a sibling
package this app imports rather than reimplements. This app is just the
interactive frontend on top of it: a local HTTP server plus an embedded
static UI, the same shape as every other applet in this family.

## Notes

- No .icns has been added to `packaging/macos/icon-builder.app` yet — same
  gap every other applet in this family currently has. Once this tool can
  export a rasterized icon (today it only outputs SVG), using it on itself
  to fill that in would be a fitting first real use.

## License

MIT — see [LICENSE](LICENSE).
