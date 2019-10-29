# Bombadillo

Bombadillo is a non-web client for the terminal, and functions as a pager/terminal UI.

Bombadillo features vim-like keybindings, configurable settings, and a robust command selection. Currently, Bombadillo supports the following protocols as first class citizens:
* gopher
* gemini
* finger
* local (a user's file system)

Support for the following protocols is also available via integration with 3rd party applications:
* telnet
    * Links are opened in a telnet application run as a subprocess
* http/https
    * Web support is opt-in (turned off by default)
    * Links can be opened in a user's default web browser when in a graphical environment
    * Web pages can be rendered directly in Bombadillo if [Lynx](https://lynx.invisible-island.net/) is installed on the system to handle the document parsing.


## Getting Started

These instructions will get a copy of the project up and running on your local machine. The following only applies if you are building from source (rather than using a precompiled binary).

### Prerequisites

You will need to have [Go](https://golang.org/) version >= 1.11.

### Building, Installing, Uninstalling

Bombadillo installation uses `make`. It is also possible to use Go to build and install (i.e `go build`, `go install`), but this is not the recommended approach.

Running `make` from the source code directory will build Bombadillo in the local directory. This is fine for testing or trying things out. For usage system-wide, and easy access to documentation, follow the installation instructions below. 

#### Basic Installation

Most users will want to install using the following commands:

```
git clone https://tildegit.org/sloum/bombadillo.git
cd bombadillo
sudo make install
```
*Note: the usage of `sudo` here will be system dependent. Most systems will require it for installation to `/usr/local/bin`.*

You can then start Bombadillo by running the command:
```
bombadillo
```
To familiarize yourself with the application, documentation is available by running the command:
```
man bombadillo
```

#### Custom Installation

There are a number of default configuration options in the file `defaults.go`. These can all be set prior to building in order to have these defaults apply to all users of Bombadillo on a given system. That said, the basic configuration already present should be suitable for most users (and all settings but one can be changed from within a Bombadillo session).

The installation location can be overridden at compile time, which can be very useful for administrators wanting to set up Bombadillo on a multi-user machine. 

```
git clone https://tildegit.org/sloum/bombadillo.git
cd bombadillo
sudo make DESTDIR=/some/directory install
```

There are two things to know about when using the above format:
1. The above would install Bombadillo to `/some/directory/usr/local/bin`, _not_ to `/some/directory`. So you will want to make sure your `$PATH` is set accordingly.
2. Using the above will install the man page to `/some/directory/usr/local/share/man`, rather than its usual location. You will want to update your `manpath` accordingly.

#### Uninstall

If you used the makefile to install Bombadillo then uninstalling is very simple. From the Bombadillo source folder run:

```
sudo make uninstall
```

If you used a custom `DESTDIR` value during install, you will need to supply it when uninstalling:
```
sudo make DESTDIR=/some/directory uninstall
```

Uninstall will clean up any build files, remove the installed binary, and remove the man page from the system. If will _not_ remove any directories created as a part of the installation, nor will it remove any Bombadillo user configuration files.

#### Troubleshooting

If you run `bombadillo` and get `bombadillo: command not found`, try running `make` from within the cloned repo. Then try: `./bombadillo`. If that works it means  that the application is getting built correctly and the issue is likely in your path settings. Any errors during `make install` should be pretty visible, as you will be able to see what command it failed on.

### Downloading

If you would prefer to download a binary for your system, rather than build from source, please visit the [Bombadillo releases](http://bombadillo.colorfield.space/releases) page. Don't see your OS/architecture? Bombadillo can be built for use with any system that is supported as a target for the Go compiler (Linux, BSD, OS X, Plan 9). There is no explicit support for, or testing done for, Windows or Plan 9. The program should build on those systems, but you may encounter unexpected behaviors or incompatibilities.

### Documentation

Bombadillo's primary documentation can be found in the man entry that installs with Bombadillo. To access it run `man bombadillo` after first installing Bombadillo. If for some reason that does not work, the document can be accessed directly from the source folder with `man ./bombadillo.1`.

In addition to the man page, users can get information on Bombadillo on the web @ [http://bombadillo.colorfield.space](http://bombadillo.colorfield.space). Running the command `help` inside Bombadillo will navigate a user to the gopher server hosted at [bombadillo.colorfield.space](gopher://bombadillo.colorfield.space); specifically the user guide.

## Contributing

Bombadillo development is largely handled by Sloum, with help from asdf, jboverf, and some community input. If you would like to get involved, please reach out or submit an issue.

At present the developers use the tildegit issues system to discuss new features, track bugs, and communicate with users about hopes and/or issues for/with the software.

If you have forked and would like to make a pull request, please make the pull request into develop where it will be reviewed by one of the maintainers. That said, a heads up or comment/issue somewhere is advised. While input is always welcome, not all requests will be granted. That said, we do our best to make Bombadillo a useful piece of software for its users and in general want to help you out.

## License

This project is licensed under the GNU GPL version 3. See the [LICENSE](LICENSE) file for details.

## Releases

Starting with version 2.0.0 releases into `master` will be version-tagged. Work done toward the next release will be created on work branches named for what they are doing and then merged into `develop` to be combined with other ongoing efforts before a release is merged into `master`. At present there is no specific release schedule. It will depend on the urgency of the work that makes its way into develop and will be up to the project maintainers' judgement when to release from `develop`.

