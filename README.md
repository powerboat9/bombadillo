# Bombadillo

Bombadillo is a non-web client for the terminal, and functions as a pager/terminal UI. Bombadillo features vim-like keybindings, configurable settings, and a robust command selection. Currently, Bombadillo supports the following protocols:
- gopher
- gemini
- finger
- local (a user's filesystem)
- telnet (by opening links in a subprocess w/ a telnet application)
- http/https links can be opened in a user's default web browser as an opt-in behavior


## Getting Started

These instructions will get a copy of the project up and running on your local machine. 

### Prerequisites

If building from source, you will need to have [Go](https://golang.org/) version >= 1.12.

### Building, Installing, Uninstalling

Bombadillo installation uses `make`. It is possible to just use the go compiler directly (`go install`) if you do not wish to have a man page installed and do not want a program to manage uninstalling and cleaning up.

By default Bombadillo will try to install to `$GOBIN`. If it is not set it will try `$GOPATH/bin` (if `$GOPATH` is set), otherwise `~/go/bin`.

#### Basic Installation

Once you have `go` installed you can build a few different ways. Most users will want the following:

```
git clone https://tildegit.org/sloum/bombadillo.git
cd bombadillo
make install
```

Once that is done you should be able to run `bombadillo` (assuming that one of the default install locations exists and is on your path) or view the manual with `man bombadillo`.

#### Custom Installation

There are a number of default configuration options in the file `defaults.go`. These can all be set prior to building in order to have these defaults apply to all users of Bombadillo on a given system. That said, the basic configuration already present should be suitable for most users (and all settings but one can be changed from within a Bombadillo session).

The installation location can be overridden at compile time, which can be very useful for administrators wanting to set up Bombadillo on a multi-user machine. If you wanted to install to, for example, `/usr/local/bin` you would do the following:

```
git clone https://tildegit.org/sloum/bombadillo.git
cd bombadillo
make install BUILD_PATH=/usr/local/bin
```

#### Uninstall

If you used the makefile to install Bombadillo then uninstalling is very simple. From the Bombadillo source folder run:

```
make uninstall
```

#### Troubleshooting

If you run `bombadillo` and get `bombadillo: command not found`, try running `make build` from within the cloned repo. Then try: `./bombadillo`. If that works it means that Go does not install to your path, or the custom path you selected is not on your path. Try the custom install from above to a location you know to be on your path.

### Downloading

If you would prefer to download a binary for your system, rather than build from source, please visit the [Bombadillo downloads](https://rawtext.club/~sloum/bombadillo.html#downloads) page. Don't see your OS/architecture? Bombadillo can be built for use with any POSIX compliant system that is supported as a target for the Go compiler (Linux, BSD, OS X, Plan 9). No testing has been done for Windows. The program will build, but will likely not work properly outside of the Linux subsystem. If you are a Windows user and would like to do some testing or get involved in development please reach out or open an issue.

### Documentation

Bombadillo's primary documentation can be found in the man entry that installs with Bombadillo. To access it run `man bombadillo` after first installing Bombadillo. If for some reason that does not work, the document can be accessed directly from the source folder with `man ./bombadillo.1`.

In addition to the man page, users can get information on Bombadillo on the web @ [http://bombadillo.colorfield.space](http://bombadillo.colorfield.space). Running the command `help` inside Bombadillo will navigate a user to the gopher server hosted at [bombadillo.colorfield.space](gopher://bombadillo.colorfield.space), specifically the user guide.

## Contributing

Bombadillo development is largely handled by Sloum, with help from asdf, jboverf, and some community input. If you would like to get involved, please reach out or submit an issue. At present the developers use the tildegit issues system to discuss new features, track bugs, and communicate with users about hopes and/or issues for/with the software. If you have forked and would like to make a pull request, please make the pull request into `develop` where it will be reviewed by one of the maintainers. That said, a heads up or comment/issue somewhere is advised. While input is always welcome, not all requests will be granted. That said, we do our best to make Bombadillo a useful piece of software for its users and in general want to help you out.

## License

This project is licensed under the GNU GPL version 3. See the [LICENSE](LICENSE) file for details.

## Releases

Starting with v2.0.0 releases into `master` will be version-tagged. Work done toward the next release will be created on work branches named for what they are doing and then merged into `develop` to be combined with other ongoing efforts before a release is merged into `master`. At present there is no specific release schedule. It will depend on the urgency of the work that makes its way into develop and will but up to the project maintainers' judgement when to release from `develop`.

