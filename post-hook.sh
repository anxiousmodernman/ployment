#/bin/sh


SOURCE_DIR=$1
TARGET_DIR=$2

echo "SOURCE_DIR is $1 and TARGET_DIR is $2"

cp -rf "$SOURCE_DIR"/* "$TARGET_DIR"
