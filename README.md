# hypr-dock <p style="display: inline; float: right; font-size: 20px; margin-top: 5px">[`ru`](https://github.com/lotos-linux/hypr-dock/blob/main/README_RU.md)</p>
### Interactive Dock Panel for Hyprland

![screenshot1](https://github.com/user-attachments/assets/b98cdf7c-83b0-4c12-9da1-ada9e1543178)
![250318_10h19m33s_screenshot](https://github.com/user-attachments/assets/3ef014e4-4613-4e28-b186-71ce262db404)

https://github.com/user-attachments/assets/50d26918-ac74-4b3b-b384-9dd98c62a799

## Installation

### Dependencies

- go (make)
- gtk3
- gtk-layer-shell

### Install
```bash
git clone https://github.com/lotos-linux/hypr-dock.git
cd hypr-dock
make get
make build
make install
```

### Uninstall
```bash
make uninstall
```

### Update
```bash
make update
```

### Local run (dev mode)
```bash
make exec
```

## Launching

### Launch Parameters:
```text
-config string
    config file (default "~/.config/hypr-dock")
-dev
    enable developer mode
-theme string
    theme dir (default "lotos")
```
#### All parameters are optional.

The default configuration and themes are installed in `~/.config/hypr-dock`

### Add the following to hyprland.conf for autostart:
```text
exec-once = hypr-dock
bind = Super, D, exec, hypr-dock
```

#### The dock supports only one running instance, so launching it again will close the previous instance.

## Configuration

### The following parameters are available in config.jsonc:
```jsonc
{
    "CurrentTheme": "lotos",

    // Icon size (px) (default 23)
    "IconSize": 23,

    // Window overlay layer height (auto, exclusive-top, exclusive-bottom, background, bottom, top, overlay) (default "auto")
    "Layer": "auto",

    // Window position on screen (top, bottom, left, right) (default "bottom")
    "Position": "bottom",

    // Use system gap (true, false) (default "true")
    "SystemGapUsed": "true",

    // Indent from the edge of the screen (px) (default 8)
    "Margin": 8
}
```
#### If a parameter is not specified, the default value will be used.

## Explanation of Non-Obvious Parameters
### Layer
- With `"Layer": "auto"` the dock layer is below all windows, but if you move the mouse cursor to the edge of the screen, the dock rises above them.
- With `"Layer": "exclusive-top"` - exclusive mode is enabled on the top layer. Neither tiled nor floating windows will overlap the dock.
- With `"Layer": "exclusive-bottom"` - exclusive mode is enabled on the bottom layer. Tiled windows won't overlap the dock. Floating windows will appear above the dock.
### SystemGapUsed
- With `"SystemGapUsed": "true"` the dock will set its margin from the edge of the screen based on the hyprland configuration, specifically the `general:gaps_out` value. The dock will dynamically adapt to changes in the hyprland configuration.
- With `"SystemGapUsed": "false"` the margin from the edge of the screen will be set by the `Margin` parameter.

### There is also a pinned.json file for pinned applications
#### Example:
```json
{
  "Pinned": [
    "firefox",
    "org.telegram.desktop",
    "code-oss",
    "kitty"
  ]
}
```
You can edit it manually. But why? ¯\_(ツ)_/¯

## Themes

#### Themes are located in the `~/.config/hypr-dock/themes/` folder

### A theme consists of:
- `[theme_name].jsonc`, for example `lotos.jsonc`
- `style.css`
- A folder with `svg` files for indicating the number of running applications

### The theme configuration has only two parameters:
```jsonc
{
    // Blur window ("on", "off") (default "on")
    "Blur": "on",

    // Distance between elements (px) (default 8)
    "Spacing": 9
}
```
#### Feel free to customize the style.css file as you like.

## Libraries Used
- [github.com/akshaybharambe14/go-jsonc](https://github.com/akshaybharambe14/go-jsonc) v1.0.0
- [github.com/allan-simon/go-singleinstance](https://github.com/allan-simon/go-singleinstance) v0.0.0-20210120080615-d0997106ab37
- [github.com/dlasky/gotk3-layershell](https://github.com/dlasky/gotk3-layershell) v0.0.0-20240515133811-5c5115f0d774
- [github.com/goccy/go-json](https://github.com/goccy/go-json) v0.10.3
- [github.com/gotk3/gotk3](https://github.com/gotk3/gotk3) v0.6.3