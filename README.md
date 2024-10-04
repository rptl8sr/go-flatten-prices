# Flatten prices

Util to insert price data from SQLite over the background image and save it as a flat image.

### Hot to use

1. Check that fonts/fontName.ttf, config/config.ini and db.sqlite (with prices) exist in the root project dir
2. Put paired files (jpg|jpeg|png and csv) with the same names into the 'input' dir
3. Run app
4. Check logs in the console and in the dir logs/YYYY-MM-DD_log.txt
5. Flattened image will be saved in the 'output' dir with the same names /path
6. Change config/config.ini if you need
7. Enjoy a result (or not)

### CSV-file structure
```text
code;font_name;font_size;color;left_upper_corner_x;left_upper_corner_y,align
```

example
```text
1234;Arial;42;#123456;100;120;left
```
