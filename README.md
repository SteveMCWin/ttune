![Cool banner :D](screenshots/banner.png?raw=true)

# tTune
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
- Go 1.24
- Portaudio 19.7.0

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

<!-- ## Credits -->


