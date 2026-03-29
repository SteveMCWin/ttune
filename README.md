# tTune

![Cool banner :D](.github/extra/screenshots/banner.png?raw=true)

tTune is a stylized, extendable, keyboard oriented guitar tuning app right in your terminal!

> **Note:** tTune is still actively developed. Future releases may be incompatible with previous ones due to my lack of foresight. If you cannot run the program after updating, please back up your configuration, delete your ~/.config/ttune directory and try running tTune again.
---

## Showcase

<div align="center">
  <img src=".github/extra/screenshots/ttune_demo.gif?raw=true" /><br/>
  <em>Tuning G string demo</em>
</div>

<br/>

<div align="center">
  <img src=".github/extra/screenshots/nord_example.png?raw=true" /><br/>
  <em>Tuning view</em>
</div>

<br/>

<div align="center">
  <img src=".github/extra/screenshots/gruvbox_example.png?raw=true" /><br/>
  <em>Settings view</em>
</div>

## Features

- Pitch detection based on the [YIN algorithm](http://audition.ens.fr/adc/pdf/2002_JASA_YIN.pdf)
- Modern terminal UI
- Easily customizable settings
- Batteries included

## Dependencies

- Go `1.24`
- Portaudio `19.7.0`

### Arch Linux

```bash
sudo pacman -S pkgconf portaudio go
```

### Debian

```bash
sudo apt-get install pkg-config portaudio19-dev golang-go
```

## Installation

### Arch Linux

```bash
yay -S ttune
```

### Build From Source

```bash
# get dependencies first
git clone https://github.com/SteveMCWin/ttune.git
cd ttune
go build
```

## Customization

You can add your own options to `custom_options.json` as long as you follow the existing format.

> **Note:** If something breaks, back up and delete the edited files or the entire `~/.config/ttune/` directory - the defaults will be regenerated on the next run.

---

### Ascii Art

Place your own ASCII art files in `~/.config/ttune/art/`. To display a tuning alongside your art, insert `%%%` where you'd like the notes to appear (see `~/.config/ttune/art/utf_guitar` for reference).

> **Recommended size:** no wider than **24 characters** and no taller than **20 lines**. Oversized art will not render properly.

---

### Border Style

Not extensible at the moment `:^/`

Want this to be moddable? [Open an issue](https://github.com/SteveMCWin/ttune/issues).

---

### Color Theme

To add a new color theme, open `~/.config/ttune/custom_options.json` and find the list of color themes. Add an entry with:

| Field       | Description              |
| ----------- | ------------------------ |
| `name`      | Name of the color theme  |
| `primary`   | Hex color for borders    |
| `secondary` | Hex color for text       |
| `tertiary`  | Hex color for selections |

> **Note:** Your terminal's background color will not be affected.

---

### Displayed Tuning

To add a new tuning, open `~/.config/ttune/custom_options.json` and find the list of tunings. Add an entry with:

| Field   | Description               |
| ------- | ------------------------- |
| `name`  | Name of the tuning        |
| `notes` | List of notes it contains |

> **Note:** Each note must be **exactly 3 characters long**. Pad shorter note names with spaces (e.g., `"E4 "`).

## Planned Features

- [x] Frequency detection
- [x] Customizable settings
- [x] Separate user settings
- [x] Scrollable settings options
- [ ] Options for frequency detection
- [ ] Improve pitch detection
