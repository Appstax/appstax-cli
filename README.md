

### Pre-release note

Please note that this repository contains pre-release software, and cannot be used until the Appstax services are officially released. Want beta access? [Send us an email](ea@appstax.com) to request an early access code.


Appstax CLI 
===========

This is the official command-line tool for [Appstax](http://appstax.com).


Installing
----------

**Manual installation:** [Download the latest release](https://github.com/appstax/appstax-cli/releases/latest) and copy the appropriate binary to somewhere on your system `PATH`. 

**On windows**, place `appstax.exe` in `c:\Windows\System32` or any custom `PATH` you have configured in your environment variables.

**Quick install on Mac OSX:**, install the latest version from the terminal:  
`curl -s http://appstax.com/download/cli/install_osx.sh | sudo /bin/bash`

**Quick install on Linux:**  
`curl -s http://appstax.com/download/cli/install_linux.sh | sudo /bin/bash`


Command overview
----------------

```bash
appstax init     # Initialize current directory as an appstax app
appstax deploy   # Deploy local files to <yourapp>.appstax.io
appstax info     # Info about app configured in current directory
appstax open     # Open your browser to the specified destination
appstax login    # Log in (also done automatically from other commands)
appstax logout   # Log out
appstax serve    # Run development http server on localhost:9000
```

You can also type `appstax help` to see a list of available commands.

Building from source
--------------------

You need [go](http://golang.org/) and [goxc](http://github.com/laher/goxc) to build appstax-cli. Go is available from Homebrew: `brew install go`, and goxc is automatically installed when running the build.

Run `./build.sh` to compile for your current platform. You will find the built artifacts in `.godeps/bin`

Run `./build.sh XC` to cross-compile. Built artifacts in `.godeps/bin/appstax-xc/snapshot`

### Running

`cd .godeps/bin`, then `appstax --help` to show available options.


License
-------

[MIT License](LICENSE)


