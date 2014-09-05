+++
title = "Go development environment for Vim"
date = 2014-03-29T06:40:42Z
author = ["Fatih Arslan"]
+++

# Go development environment for Vim

## The reason for creating vim-go

Go has a very versatile toolchain and commands that makes Go programming fun.
One of the famous tools is gofmt, which automatically reformats the code
according to some predefined rules. However there are many other tools like
goimports, oracle, godef, etc.. which help to provide a more productive workflow. 

There are many independent vim plugins that integrate these tools into Vim. We
also have the offical Vim plugins that provides some basic Go support for Vim.
However there are many flaws with these plugins:

- Gofmt destroy the undo history and it's impossible to do undo on your code
- Building, running or testing are not integrated, all of these needs additional settings.
- Binary paths are hard coded
- Plugins are not using the same format. Each plugin is built differently, some support lookup under the cursor, some do not.
- Syntax highlighting can't highlight functions
- Auto completion needs a package and additional settings.
- Many other small flaws..

The main reason to create vim-go was to integrate all these plugins, fix the
flaws and provide a common and seamless experience. 


## Vim-go features

[vim-go](https://github.com/fatih/vim-go) improves and supports the following
features:

- Improved Syntax highlighting
- Auto completion support
- Integrated and improved snippets support
- Gofmt on save, keeps cursor position and doesn't break your undo history
- Go to symbol/declaration
- Automatically import packages
- Compile and build package
- Run quickly your snippet
- Run tests and see any errors in quick window
- Lint your code
- Advanced source analysis tool with oracle
- Checking for unchecked errors.

vim-go automatically installs all necessary binaries if they are not found in
the pre-defined path `~/.vim-go` (can be disabled if not desired). It comes with
pre-defined sensible settings.  Under the hood it uses `goimports`, `oracle`,
`godef`, `errcheck`, `gocode` and `golint`. Let's see how we can make use of
those tools.

## Installation

vim-go is a pathogen compatible bundle. To use it just clone it into your
bundle directory:

    $ cd ~/.vim/bundle
    $ git clone https://github.com/fatih/vim-go.git

Auto completion is enabled by default via `<C-x><C-o>`, to get real-time
completion (completion by type) install YCM:

    $ cd ~/.vim/bundle
    $ git clone https://github.com/Valloric/YouCompleteMe.git
    $ cd YouCompleteMe
    $ ./install.sh

## Development with vim-go

On first usage it tries to download and install all necessary binaries.
Let's see what commands we now can use. Below is a list of commands vim-go
supports, most of the commands are improved and new ones were introduced (some
of them are just from the official Go plugin).

vim-go uses goimports and reformats whenever you save your file. However you
have the option to disable goimports and explicitly import/drop packages:

    :Import <path>
    :ImportAs <localname> <path>
    :Drop <path>

    :DisableGoimports
    :EnableGoimports

Godoc can be called on any identifier. Just put your cursor under the
identifier and call `:Godoc`, it will open a new view and show the necessary documentation.

    :Godoc
    :Godoc <identifier>

Godef is one of my favorites tools. It's find the source definition and jumps
to that file (go to definition). To use it just put your cursor under a
identifier and hit `:Godef`. It opens a new buffer (just like ctags). You might
add the following settings to your `vimrc`, to open the definitions in vertical,
horizontal or new tab with :

    au Filetype go nnoremap <leader>v :vsp <CR>:exe "GoDef" <CR>
    au Filetype go nnoremap <leader>s :sp <CR>:exe "GoDef"<CR>
    au Filetype go nnoremap <leader>t :tab split <CR>:exe "GoDef"<CR>

Building, testing, running are all important steps in development workflow and
should be seamless integrated. vim-go has several features that you can use. First check out the build commands:

    :make
    :GoBuild

`:make` is the default Vim command to build a project. vim-go integrates it
in way that it ** doesn't produce ** any binary. That it is really useful
because it doesn't pollute your work environment. Any errors are listed in a
quickfix window and can be jumped easily with default `:cnext` and
`:cprevious`.  `:GoBuild` is similar to `:make`, but it creates a binary for
the given main package.

Sometimes we only have small main package that we can want to run and see the
output. For that we have:

    :GoRun 
    :GoRun <expand>

Just calling `:GoRun` is going to include all files that belongs to the main
package (useful for multi file programs). To run a single file just run
`:Gorun %`. You can map this to a key, like `<leader>r`:

    au Filetype go nnoremap <leader>r :GoRun %<CR>

To call `go test` just run:

    :GoTest

Another tool we have is `errcheck`, which checks unchecked errors:

    :GoErrCheck

Linting is useful to print out mistakes or tips about coding style. For example
if you don't provide any documentation comment for a function `golint` will
warn you. To call it just execute:

    :Lint

To see the dependencies of your current package run `:GoDeps`. If you have
multiple files you can easily see all source files (test files are excluded)
via `:Gofiles`.

And then we have the still experimental but powerful "oracle" tool. See the
extensive official documentation for more info: [Oracle docs](https://docs.google.com/document/d/1SLk36YRjjMgKqe490mSRzOPYEDe0Y_WQNRv-EiFYUyw/view) vim-go implements and includes the following commands (which are part of the
offical oracle vim plugin):

    :GoOracleDescribe
    :GoOracleCallees
    :GoOracleCallers
    :GoOracleCallgraph
    :GoOracleImplements
    :GoOracleChannelPeers

These are useful especially if you want to find out how your Code is
structured, how your channels are interacting with each other, which struct is
implementing which interface, etc...


## Summary and Thanks!

Thanks for the following users and projects to make this project happen:

- Go Authors for offical vim plugins
- Gocode, Godef, Golint, Oracle, Goimports, Errcheck projects and authors of those projects.
- Other vim-plugins, thanks for inspiration (vim-golang, go.vim, vim-gocode, vim-godef)

Check out the github page for fare more information (snippets, settings, etc..):

[github.com/fatih/vim-go](https://github.com/fatih/vim-go)

There is also a Youtube vide that shows vim-go in action:

[Youtube: Go development in Vim](https://www.youtube.com/watch?v=rD11pEx5h8c)

There are still tons of modifications and improvements one can make to this
setup. vim-go is a new project. Check it out and try it to see how it fits your
needs. Any improvements and feedback are welcome.












