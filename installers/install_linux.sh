#!/bin/bash

ARTIFACT_NAME=appstax_linux_386
DOWNLOAD_FILE=appstax_linux_386.tar.gz
DOWNLOAD_DIR=/tmp/appstax
INSTALL_FILE=appstax
INSTALL_DIR=/usr/local/bin

rm -rf $DOWNLOAD_DIR/$DOWNLOAD_FILE
rm -rf $DOWNLOAD_DIR/$ARTIFACT_NAME
mkdir -p $DOWNLOAD_DIR
mkdir -p $INSTALL_DIR

echo "Downloading..."
curl http://appstax.com/download/cli/$DOWNLOAD_FILE -# -o $DOWNLOAD_DIR/$DOWNLOAD_FILE

echo "Installing..."
cd $DOWNLOAD_DIR
tar -zxf $DOWNLOAD_FILE

cp $ARTIFACT_NAME/$INSTALL_FILE $INSTALL_DIR/$INSTALL_FILE
chmod 755 $INSTALL_DIR/$INSTALL_FILE

echo "Done!"
echo "Appstax command line was installed in $INSTALL_DIR"
