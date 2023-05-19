#/bin/bash
#
# list directory under $1


dst_dir=$1

if [ -d "$dst_dir" ] && [ `find  $1/* -maxdepth 0 -type d|wc -l` -gt 0 ]; then
    find  $1/* -maxdepth 0 -type d  | xargs ls  -rtd
else
    echo "ERROR: Directory not found: $dst_dir"
    exit 1
fi

# output array of directory names
# find  $1/* -maxdepth 0 -type d  | xargs ls  -rtd | awk '{print $NF}'