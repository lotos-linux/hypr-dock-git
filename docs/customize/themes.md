## Open Window Indicators  
#### The `~/.config/hypr-dock/themes/[theme_name]/point` folder contains images used to indicate open windows.  
By default, the set includes: `0.svg`, `1.svg`, `2.svg`, `3.svg`  

![250428_13h00m17s_screenshot](https://github.com/user-attachments/assets/9f9cb607-c0c7-48ef-9379-266f2b253246)  
![250428_13h32m50s_screenshot](https://github.com/user-attachments/assets/46b0e75f-2212-4e54-a2b3-02edf3965142)  

### Customization  
- The `point` folder or its contents can be deleted entirely, which will disable indicators.  
- Supported image formats: `.svg`, `.png`, `.jpg`, `.webp`.  
- At least two files are required for indicators to work (`0.*` and any other).  
- Each file must be named according to the number of windows it represents (e.g., `3.png` for 3 windows).  
- You can add unlimited images in any order.  

### Logic  
#### If you add this set: `0.svg`, `3.png`, `10.jpg`  
The indicators will behave as follows:  
- **0–2 windows**: Always loads `0.svg`.  
- **3–9 windows**: Loads `3.png`.  
- **10+ windows**: Loads `10.jpg`.  