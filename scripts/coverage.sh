#!/bin/bash

# 出力のヘッダー部分
echo "| Package | Func | Coverage |"
echo "|--------------|--------|------------|"

# 標準入力を読み取ってテーブル行に変換
while IFS= read -r line
do
    # ファイルパス、関数名、カバレッジを抽出
    file_path=$(echo "$line" | awk -F: '{print $1}')
    line_number=$(echo "$line" | awk -F: '{print $2}')
    function_name=$(echo "$line" | awk -F: '{print $3}' | awk '{$1=$1;print}')
    coverage=$(echo "$line" | awk -F: '{print $4}' | awk '{$1=$1;print}')
    
    # Markdownのテーブル行として出力
    echo "| $file_path:$line_number | $function_name | $coverage |"
done

