#!/bin/bash

# $1 = folder location (time part)
# $2 = version number

FOLDER_LOCATION=$1
VERSION_NUMBER=$2
CONFIG_FILE="config.json"
DEST_FOLDER="configs"

mkdir -p "$DEST_FOLDER"

if [ ! -d "test_reports/test_report_$FOLDER_LOCATION" ]; then
    echo "Error: folder $FOLDER_LOCATION does not exist!"
    exit 1
fi

if [ ! -f "test_reports/test_report_$FOLDER_LOCATION/$CONFIG_FILE" ]; then
    echo "Error: No $CONFIG_FILE in $FOLDER_LOCATION"
    exit 1
fi

cp "test_reports/test_report_$FOLDER_LOCATION/$CONFIG_FILE" "$DEST_FOLDER/$VERSION_NUMBER.json"

echo "Successfully copied $CONFIG_FILE to $DEST_FOLDER/$VERSION_NUMBER.json"