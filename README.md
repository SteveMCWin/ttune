# tTune

![Cool banner :D](screenshots/banner.png?raw=true)

tTune is a stylized, keyboard oriented guitar tuning app right in your terminal!

## Screenshots

![Example Nord](screenshots/nord_example.png?raw=true)
![Example Gruvbox](screenshots/gruvbox_example.png?raw=true)

## Features

- Pitch detection based on the [YIN algorithm](http://audition.ens.fr/adc/pdf/2002_JASA_YIN.pdf)
- Modern terminal UI
- Easily customizable settings
- Batteries included

## Dependencies

- Go `1.24`
- Portaudio `19.7.0`

<details>
<summary><b>Arch Linux</b></summary>
```bash
sudo pacman -S pkgconf portaudio go
```

</details>

<details>
<summary><b>Debian</b></summary>
```bash
sudo apt-get install pkg-config portaudio19-dev golang-go
```

</details>

## Installation

<details>
<summary><b>Arch Linux</b></summary>
```bash
yay -S ttune
```

</details>

<details>
<summary><b>Build From Source</b></summary>
```bash
# get dependencies first
git clone https://github.com/SteveMCWin/ttune.git
cd ttune
go build
```

</details>

## Customization

The existing settings are modifiable via `settings_data.json`, which is created on first launch. You can also add your own options to this file as long as you follow the existing format.

> **Note:** If something breaks, delete the edited file or the entire `~/.config/ttune/` directory — defaults will be regenerated on the next run.

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

To add a new color theme, open `~/.config/ttune/settings_data.json` and find the list of color themes. Add an entry with:

| Field | Description |
|---|---|
| `name` | Name of the color theme |
| `primary` | Hex color for borders |
| `secondary` | Hex color for text |
| `tertiary` | Hex color for selections |

> **Note:** Your terminal's background color will not be affected.

---

### Displayed Tuning

To add a new tuning, open `~/.config/ttune/settings_data.json` and find the list of tunings. Add an entry with:

| Field | Description |
|---|---|
| `name` | Name of the tuning |
| `notes` | List of notes it contains |

> **Note:** Each note must be **exactly 3 characters long**. Pad shorter note names with spaces (e.g., `"E4 "`).

## Planned Features

- [x] Frequency detection
- [x] Customizable settings
- [ ] Separate user settings
- [ ] Improve pitch detection
- [ ] Scrollable settings options
- [ ] Options for frequency detection
