#!/bin/bash

ARTIFACT_NAME=appstax_linux_386
DOWNLOAD_FILE=appstax_linux_386.tar.gz
DOWNLOAD_DIR=/tmp/appstax
INSTALL_FILE=appstax
INSTALL_DIR=/usr/local/bin

echo "The appstax command-line tool will be installed in $INSTALL_DIR"
read -r -p "Press enter to continue or Ctr-C to stop." confirm

rm -rf $DOWNLOAD_DIR/$DOWNLOAD_FILE
rm -rf $DOWNLOAD_DIR/$ARTIFACT_NAME
mkdir -p $DOWNLOAD_DIR
mkdir -p $INSTALL_DIR

DOWNLOAD_URL=$(curl -s https://api.github.com/repos/appstax/appstax-cli/releases | grep "browser_download_url.*$ARTIFACT_NAME" | head -n 1 | cut -d '"' -f 4)

echo "Downloading..."
curl -L $DOWNLOAD_URL -# -o $DOWNLOAD_DIR/$DOWNLOAD_FILE || exit 0

echo "Installing..."
cd $DOWNLOAD_DIR
tar -zxf $DOWNLOAD_FILE || exit 0

cp $ARTIFACT_NAME/$INSTALL_FILE $INSTALL_DIR/$INSTALL_FILE || exit 0
chmod 755 $INSTALL_DIR/$INSTALL_FILE || exit 0

echo "Done!"
echo "Appstax command line was installed in $INSTALL_DIR"
