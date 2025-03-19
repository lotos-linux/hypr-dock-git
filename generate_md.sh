#!/bin/bash

output_file="project_structure.md"

> "$output_file"

echo "# Структура проекта" >> "$output_file"
echo "" >> "$output_file"

echo "## Структура проекта" >> "$output_file"
echo '```' >> "$output_file"
tree >> "$output_file"
echo '```' >> "$output_file"
echo "" >> "$output_file"

add_file_content() {
    local file="$1"
    local filename=$(basename "$file")
    local extension="${filename##*.}"

    if [[ "$filename" == "$extension" ]]; then
        extension="$filename"
    fi

    echo "## $file" >> "$output_file"
    echo '```'"$extension" >> "$output_file"
    cat "$file" >> "$output_file"
    echo '```' >> "$output_file"
    echo "" >> "$output_file"
}

find . -type f ! -path "./bin/*" ! -path "./.git/*" ! -path "./$output_file" ! -path "./generate_md.sh" | while read -r file; do
    add_file_content "$file"
done

echo "$output_file"
