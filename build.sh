
# Set new GOPATH
export GOPATH=$(pwd)/.godeps

# Get dependencies
mkdir -p .godeps/{src,pkg,bin}
if [ "$1" == "--nodeps" ] 
then
    sleep 0
else    
    curl -s https://raw.githubusercontent.com/pote/gpm/v1.2.3/bin/gpm | bash
fi

# Configure local depenceny name
ln -nsf ../.. .godeps/src/appstax-cli

# Build application
export GOBIN=$GOPATH/bin
echo "Building appstax-cli"
go install appstax-cli/appstax || exit 1 

# Cross-compilation
if [ "$1" = "XC" ]; then
	echo "Cross-compiling"
	go get github.com/laher/goxc
    $GOPATH/bin/goxc -wd=appstax

    rm -rf .godeps/bin/appstax-xc/snapshot/.goxc-temp
    rm -rf .godeps/bin/appstax-xc/snapshot/downloads.md
    cp installers/install_linux.sh .godeps/bin/appstax-xc/snapshot/install_linux.sh
    cp installers/install_osx.sh .godeps/bin/appstax-xc/snapshot/install_osx.sh
fi

