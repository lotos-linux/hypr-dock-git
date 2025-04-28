## Open Window Indicators  
#### Place indicator images in `~/.config/hypr-dock/themes/[theme_name]/point`.  
**Default set:** `0.svg` `1.svg` `2.svg` `3.svg`  

![Preview](https://github.com/user-attachments/assets/9f9cb607-c0c7-48ef-9379-266f2b253246)  
![Preview](https://github.com/user-attachments/assets/46b0e75f-2212-4e54-a2b3-02edf3965142)  

### Customization  
- Delete the `point` folder/files to **disable indicators entirely**  
- Supported formats: `.svg`, `.png`, `.jpg`, `.webp`  
- **Minimum requirement:** Two files (`0.*` and any other)  
- Naming must match window counts (e.g. `5.png` for 5 windows)  
- Add unlimited files in any order  

### Logic Example  
#### For set: `0.svg` `3.png` `10.jpg`  
- **0-2 windows:** Uses `0.svg`  
- **3-9 windows:** Uses `3.png`  
- **10+ windows:** Uses `10.jpg`  