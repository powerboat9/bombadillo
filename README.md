# Bombadillo

Bombadillo is a modern [Gopher](https://en.wikipedia.org/wiki/Gopher_(protocol)) client for the terminal, and functions as a pager/terminal UI. Bombadillo features vim-like keybindings, configurable settings, and a robust command selection. Bombadillo is under active development.


## Getting Started

These instructions will get a copy of the project up and running on your local machine. 

### Prerequisites

If building from source, you will need to have [Go](https://golang.org/) version >= 1.11. Bombadillo uses the module system, so if using 1.11 you will need to have that feature enabled. If using a version > 1.11, you already have modules enabled.

Bombadillo does not use any outside dependencies beyond the Go standard library.

### Installing

Assuming you have `go` installed, run the following:

```
git clone https://tildegit.org/sloum/bombadillo.git
cd bombadillo
go install
```

Once you have done that you should, assuming `go install` is set up to install to a place on your path, you should be able to run the following from anywhere on your system to use Bombadillo:

```
bombadillo
```

#### Troubleshooting

If you run `bombadillo` and get `bombadillo: command not found`, try running `go build` from within the cloned repo. Then try: `./bombadillo`. If that works it means that Go does not install to your path. `go build` added an executable file to the repo directory. Move that file to somewhere on your path. I suggest `/usr/local/bin` on most systems, but that may be a matter of personal preference.

### Downloading

If you would prefer to download a binary for your system, rather than build from source, please visit the [Bombadillo downloads](https://rawtext.club/~sloum/bombadillo.html#downloads) page. Don't see your OS/architecture? Bombadillo can be built for use with any POSIX compliant system that is supported as a target for the Go compiler (Linux, BSD, OS X, Plan 9). No testing has been done for Windows. The program will build, but will likely not work properly outside of the Linux subsystem. If you are a Windows user and would like to do some testing or get involved in development please reach out or open an issue.

### Documentation

Bombadillo has documentation available in two places currently. The first if the [Bombadillo homepage](https://rawtext.club/~sloum/bombadillo.html#docs), which has lots of information about the program, links to places around Gopher, and documentation of the commands and configuration options. 

Secondly, and possibly more importantly, documentation is available via Gopher from within Bombadillo. When a user launches Bombadillo for the first time, their `homeurl` is set to the help file. As such they will have access to all of the key bindings, commands, and configuration from the first run. A user can also type `:?` or `:help` at any time to return to the documentation. Remember that Bombadillo uses vim-like key bindings, so scroll with `j` and `k` to view the docs file.

## Contributing

Bombadillo development is largely handled by Sloum, with help from jboverf and some community input. If you would like to get involved, please reach out or submit an issue. At present the developers use the tildegit issues system to discuss new features, track bugs, and communicate with users about hopes and/or issues for/with the software.

## License

This project is licensed under the GNU GPL version 3- see the [LICENSE](LICENSE) file for details.

