package main

type HelpItem struct {
	Name string
	Contents string
}

func InitHelpItems() []HelpItem {
	about := HelpItem{
		Name: "ℹ About",
		Contents:
`
This is a small personal project (that took about 10x longer to complete than expected hihi).
I once complained to my cousin about burnout, so he gave me a guitar, saying a creative hobby is a really good way to combat fatigue. After a couple of months of actually trying to learn how to play the thing, I got the idea to make this little program :^). I haven't played the guitar since, except for plucking a few strings to check if the tuning is accurate, but that counts too, right?
Why a TUI? I was learning golang for backend web development and somehow stumbled accross the bubbletea TUI framework, and was impressed with how good terminal apps can look, and wanted to make something. The first project idea I started working on had too big of a scope, so I decided to go and make something in a weekend or two, and this is where that led me.
`,
	}

	modding := HelpItem{
		Name: "⚙ Modding",
		Contents:
`
The existing settings are modifiable and extensible, mostly through JSON. Here's how to add more options to the settings:

- Ascii Art
		In your ~/.config/ttune/art direcotry you can add your own ascii art as a file. Note that for a tuning to be displayed along your ascii art, place '%%%' where you wish the notes to appear (check out the ~/.config/ttune/art/utf_guitar file for reference). Also I recommend you keep the art no more than 24 characters wide and 20 lines tall.

- Border Style
		Not extensible at the moment :^/. In case you want this to be moddable, open an issue at https://github.com/SteveMCWin/ttune/issues.

- Color Theme
		To add a new color theme, go into the file located at ~/.config/ttune/settings_data.json where you will see a list of color themes. Add an item containing the name of the color theme you want and three colors in hex code. The primary color is used for the borders, secondary for the text and tertiary for selections. Note that the background of your terminal will not change colors.

- Displayed Tuning
		To add a new tuning that will be displayed next to the ascii art, go into the file located at ~/.config/ttune/settings_data.json where you will see a list of tunings. Add an item containing the name of the tuning you want and a list of notes it is comprised of. Note that the notes should be exactly 3 characters long. For notes that require less characters, fill out the rest of the string with spaces.
`,
	}

	support := HelpItem{
		Name: "★ Support",
		Contents:
`
I have a goal of one of my repos reaching at least 10 stars by the end of 2026, so if you have one to spare, you can do your magic at https://github.com/SteveMCWin/ttune

If that's too much to ask for, you can always send me a tip at https://ko-fi.com/stevemcwin instead :^]
`,
	}

	return []HelpItem{about, modding, support}
}
