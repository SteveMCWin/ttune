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

## Customization

The existing settings are modifiable, you can edit the settings_data.json file that will be created once you start the app for the first time.\
You can also add your own settings options in the settings_data.json file, as long as you follow the format.\
Note that if you break something, you can delete the file you edited or the whole ~/.config/ttune directory, and the defaults will be re-generated the next time you run ttune.\
Here's how to add more options to the settings:

### Ascii Art
		In your ~/.config/ttune/art direcotry you can add your own ascii art as a file. Note that for a tuning to be displayed along your ascii art, place '%%%' where you wish the notes to appear (check out the ~/.config/ttune/art/utf_guitar file for reference). Also I recommend you keep the art no more than 24 characters wide and 20 lines tall. Ascii art that is too big will not be rendered properly.

### Border Style
		Not extensible at the moment :^/. In case you want this to be moddable, open an issue at https://github.com/SteveMCWin/ttune/issues.

### Color Theme
		To add a new color theme, go into the file located at ~/.config/ttune/settings_data.json where you will see a list of color themes. Add an item containing the name of the color theme you want and three colors in hex code. The primary color is used for the borders, secondary for the text and tertiary for selections. Note that the background of your terminal will not change colors.

### Displayed Tuning
		To add a new tuning that will be displayed next to the ascii art, go into the file located at ~/.config/ttune/settings_data.json where you will see a list of tunings. Add an item containing the name of the tuning you want and a list of notes it is comprised of. Note that the notes should be exactly 3 characters long. For notes that require less characters, fill out the rest of the string with spaces.

## Planned features
- ~Frequency detection~
- ~Customizable settings~
- Separate user settings
- Improve  pitch detection
- Scrollable settings options
- Options for frequency detection

<!-- ## Credits -->


