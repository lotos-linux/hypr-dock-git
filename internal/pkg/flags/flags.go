package flags

import "flag"

type Flags struct {
	Config string
	Theme  string
}

func Get(defaultTheme string) Flags {
	list := Flags{
		Config: *flag.String("config", "~/.config/hypr-dock", "config file"),
		Theme:  *flag.String("theme", defaultTheme, "theme dir"),
	}
	flag.Parse()
	return list
}
