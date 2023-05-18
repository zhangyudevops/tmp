#/bin/bash
yaml_file=$1
job_name=$2
new_image=$3

if [ -z "$yaml_file" ] || [ -z "$new_image" ]; then
    echo "Usage: $0 <yaml_file> <image_name>"
    exit 1
fi

if [ ! -f "$yaml_file" ]; then
    echo "File $yaml_file does not exist"
    exit 1
fi

old_image="$(grep -E "image:.*" $yaml_file | grep $job_name | cut -d':' -f 2- | sed -e 's/^[ ]*//g'

sed -i "s|$old_image|$new_image|g" $yaml_file
echo "Image $new_image updated in $yaml_file"
