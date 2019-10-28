# Bombadillo - a non-web client

Bombadillo is a modern [Gopher](https://en.wikipedia.org/wiki/Gopher_(protocol)) client for the terminal, and functions as a pager/terminal UI. Bombadillo features vim-like keybindings, configurable settings, and a robust command selection. Bombadillo is under active development.


## Getting Started

These instructions will get a copy of the project up and running on your local machine. 

### Prerequisites

If building from source, you will need to have [Go](https://golang.org/) version >= 1.11.

#### Optional

[Lynx](https://lynx.invisible-island.net/), the text based web browser, can be used as a parsing engine for http/https. This is a totally optional item and Lynx is in no way required in order to compile or run Bombadillo. Having it available on the system can help create a richer experience by allowing users to navigate directly to web content in Bombadillo. Many users may wish to avoid this entirely, and the default configuration does not have this behavior turned on. To turn it on from within Bombadillo enter the command `set lynxmode true`.

### Installing

Assuming you have all prerequisites installed, Bombadillo can be installed on your system using the following commands:

```
git clone https://tildegit.org/sloum/bombadillo.git
cd bombadillo
sudo make install
```

Once this is done, you should be able to start Bombadillo using the following command:

```
bombadillo
``` 

#### Other installation options

If you only want to install Bombadillo for your own user account, you could try the following in the cloned repository:

```
make PREFIX=~ install
```

You can then add `~/bin` to your PATH environment variable, and `~/share/man` to your manpath.

The `PREFIX` option can be used to install Bombadillo to any location different to `/usr/local`.

#### Troubleshooting

If you run `bombadillo` and get `bombadillo: command not found`, try running `make` from within the cloned repository. Next, try: `./bombadillo`. If this works, it means that the installation was not completed to an area in your `PATH`.

Please feel free to [open an issue](https://tildegit.org/sloum/bombadillo/issues) if you experience any other issues with the installation.

### Uninstalling

To uninstall Bombadillo, simply run the following command from the cloned repository:

```
sudo make uninstall
```

Please note that directories created during the installation will not be removed.

### Downloading

If you would prefer to download a binary for your system, rather than build from source, please visit the [Bombadillo downloads](https://rawtext.club/~sloum/bombadillo.html#downloads) page. Don't see your OS/architecture? Bombadillo can be built for use with any POSIX compliant system that is supported as a target for the Go compiler (Linux, BSD, OS X, Plan 9). No testing has been done for Windows. The program will build, but will likely not work properly outside of the Linux subsystem. If you are a Windows user and would like to do some testing or get involved in development please reach out or [open an issue](https://tildegit.org/sloum/bombadillo/issues).

### Documentation

Bombadillo has documentation available in four places currently. The first is the [Bombadillo homepage](https://rawtext.club/~sloum/bombadillo.html#docs), which has lots of information about the program, links to places around Gopher, and documentation of the commands and configuration options.

Secondly, and possibly more importantly, documentation is available via Gopher from within Bombadillo. When a user launches Bombadillo for the first time, their `homeurl` is set to the help file. As such they will have access to all of the key bindings, commands, and configuration from the first run. A user can also type `:?` or `:help` at any time to return to the documentation. Remember that Bombadillo uses vim-like key bindings, so scroll with `j` and `k` to view the docs file.

Thirdly, this repo contains a file `bombadillo-info`. This is a duplicate of the help file that is hosted over gopher mentioned above. Per user request it has been added to the repo so that pull requests can be created with updates to the documentation.

Lastly, but perhaps most importantly, a manpage is now included in the repo as `bombadillo.1`. Current efforts are underway to automate the install of both bombadillo and this manpgage.

## Contributing

Bombadillo development is largely handled by Sloum, with help from jboverf, asdf, and some community input. If you would like to get involved, please reach out or submit an issue. At present the developers use the tildegit issues system to discuss new features, track bugs, and communicate with users about hopes and/or issues for/with the software. If you have forked and would like to make a pull request, please make the pull request into `develop` where it will be reviewed by one of the maintainers. That said, a heads up or comment/issue somewhere is advised.

## License

This project is licensed under the GNU GPL version 3- see the [LICENSE](LICENSE) file for details.

## Releases

Starting with v2.0.0 releases into `master` will be version-tagged. Work done toward the next release will be created on work branches named for what they are doing and then merged into `develop` to be combined with other ongoing efforts before a release is merged into `master`. At present there is no specific release schedule. It will depend on the urgency of the work that makes its way into develop and will but up to the project maintainers' judgement when to release from `develop`.

