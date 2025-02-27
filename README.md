# hypr-dock
### Интерактовная док-панель для Hyprland

![Screen shot](https://raw.githubusercontent.com/lotos-linux/hypr-dock/refs/heads/main/github/screenshot1.png)

## Установка

### Зависимости

- `go` (make)
- `gtk3`
- `gtk-layer-shell`

### Установка

```bash
git clone https://github.com/lotos-linux/hypr-dock.git
cd hypr-dock
make get
make build
make install
```

## Запуск

### Параметры запуска:

```text
  -config string
    	config file (default "~/.config/hypr-dock/config.jsonc")
  -theme string
    	theme (default "lotos")
  -help
```
#### Все параметры являются необязательными.

Конфигурация и темы по умолчания ставяться в `~/.config/hypr-dock`
### Добавьте запуск в `hyprland.conf`:

```text
exec-once = hypr-dock [arguments]
```
## Настройка

### В `config.jsonc` доступны такие параметры

```jsonc
{
    "CurrentTheme": "lotos",

    // Icon size (px) (default 21)
    "IconSize": 21,

    // Window overlay layer height (auto, background, bottom, top, overlay) (default "auto")
    "Layer": "auto",

    // Window position on screen (top, bottom, left, right) (default "bottom")
    "Position": "bottom",

    // Indent from the edge of the screen (px) (default 8)
    "Margin": 8
}
```
#### Если параметр не указан значение будет выставлено по умолчанию
### Автопереключение слоя дока
#### При `"Layer": "auto"` окно дока находиться под всеми окнами, но если увести курсор мыши к краю экрана - док поднимается над всеми окнами


### Также есть файл `pinned.json` с закрепленными приложениями
#### Например
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
Вы можете менять его в ручную. Но зачем? ¯\_(ツ)_/¯

## Темы

#### Темы находяться в папке `~/.config/hypr-dock/themes/`

### Тема состоит из нескольких файлов
- `[название_темы].jsonc` например `lotos.jsonc`
- `style.css`
- Четыре `svg` файла с иконками индикатора запуска
- Иконка-заглушка: `Необязательный` (Если в теме ее нет - ищет в `~/.config/hypr-dock`)

### В конфиге темы всего два параметра
```jsonc
{
    // Blur window ("on", "off") (default "on")
    "Blur": "on",

    // Distance between elements (px) (default 8)
    "Spacing": 9
}
```
#### Файл `style.css` крутите как хотите 

## Использованые библеотки
- <github.com/akshaybharambe14/go-jsonc> `v1.0.0`
- <github.com/allan-simon/go-singleinstance> `v0.0.0-20210120080615-d0997106ab37`
- <github.com/dlasky/gotk3-layershell> `v0.0.0-20240515133811-5c5115f0d774`
- <github.com/goccy/go-json> `v0.10.3`
- <github.com/gotk3/gotk3> `v0.6.3`
