# hypr-dock
### Интерактовная док-панель для Hyprland

![screenshot1](https://github.com/user-attachments/assets/b98cdf7c-83b0-4c12-9da1-ada9e1543178)
![250318_10h19m33s_screenshot](https://github.com/user-attachments/assets/3ef014e4-4613-4e28-b186-71ce262db404)


https://github.com/user-attachments/assets/50d26918-ac74-4b3b-b384-9dd98c62a799



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
sudo make install
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
exec-once = hypr-dock
bind = Super, D, exec, hypr-dock
```

#### Док поддерживает только один запущенный экземпляр, так что повторный запуск закроет предыдующий.

## Настройка

### В `config.jsonc` доступны такие параметры

```jsonc
{
    "CurrentTheme": "lotos",

    // Icon size (px) (default 23)
    "IconSize": 23,

    // Window overlay layer height (auto, background, bottom, top, overlay) (default "auto")
    "Layer": "auto",

    // Window position on screen (top, bottom, left, right) (default "bottom")
    "Position": "bottom",

    // Use system gap (true, false) (default "true")
    "SystemGapUsed": "true",

    // Indent from the edge of the screen (px) (default 8)
    "Margin": 8
}
```
#### Если параметр не указан значение будет выставлено по умолчанию
## Разберем неочевидные параметры
### Layer
#### При `"Layer": "auto"` слой дока находиться под всеми окнами, но если увести курсор мыши к краю экрана - док поднимается над ними
### SystemGapUsed
#### При `"SystemGapUsed": "true"` док будет задавать для себя отступ от края экрана беря значение из конфигурации `hyprland`, а конкретно значения `general:gaps_out`, при этом док динамически будет подхватывать изменение конфигурации `hyprland`
#### При `"SystemGapUsed": "false"` отступ от края экрана будет задаваться параметром `Margin`

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

### Тема состоит из
- `[название_темы].jsonc` например `lotos.jsonc`
- `style.css`
- Папка с `svg` файлами для индикации количества запущенных приложения

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

## Использованные библиотки
- <github.com/akshaybharambe14/go-jsonc> `v1.0.0`
- <github.com/allan-simon/go-singleinstance> `v0.0.0-20210120080615-d0997106ab37`
- <github.com/dlasky/gotk3-layershell> `v0.0.0-20240515133811-5c5115f0d774`
- <github.com/goccy/go-json> `v0.10.3`
- <github.com/gotk3/gotk3> `v0.6.3`
