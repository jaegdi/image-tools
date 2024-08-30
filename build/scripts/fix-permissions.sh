#!/bin/bash
files=$1
user=${2:-root}
find $(echo $files|tr ',' ' ') -print0 |
    while IFS= read -r -d '' filename; do
    chown $user:root "$filename"
    chmod g+u "$filename"
    done
