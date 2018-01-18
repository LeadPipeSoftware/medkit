<!--![MEDKIT](medkit.png "MEDKIT")-->
# MEDKIT
> MEDKIT (Multi-Environment Dotfiles Kit) is the dotfile management solution for the obsessive compulsive.

[![GitHub issues][github-issues-image]][github-issues-url]

MEDKIT is a tool to help you take control of your local environment configuration.

### Manage Your Dotfiles
The trouble with dotfiles is that, unless you have just one login on one computer, you'll have lots of them. If you're
like most people, you want your preferences and configuration to be consistent no matter where you log in.

MEDKIT follows a "Write Once, Use Anywhere" philosophy. With a little help from your favorite version control tool, you
can keep a "master" copy of your dotfiles and apply them anywhere you go.

- Apply your dotfiles
- Scan for new dotfiles

### Install Your Software
Getting all your tools installed is often one of the biggest challenges when you start using a new computer for the
first time. With MEDKIT, you can script your installation scripts and run them as-needed.

- Run `*.installer` scripts
- Remember what has already been run
- Reset if needed

### Organize Your Startup
Startup files often start simple, but turn into a nightmare to manage. It gets worse when tools and installers try to be
helpful and make changes to them without asking. You can take advantage of MEDKIT's version control integration and
templates to get your startup scripts under control once and for all.

- Start with the MEDKIT template or roll your own
- Run `*.source` scripts at startup
- Determine the order each source script runs in

### Organize Your PATH
Much like startup scripts, keeping your environment's `PATH` in order can be a real challenge. MEDKIT's convention-based
system can help!

- Run `*.path` scripts
- Determine the order of each PATH element

<!--
## Philosophy
Bringing forth the notion of topics is one of the simple yet genius things Zach brought to his dotfiles. This makes it
super easy to organize things logically and keep things tidy. For example, if you start using the greatest editor of all
time, vim, then you create a `vim` folder, drop your `.vimrc` file (named as `.vimrc.symlink`) in there, and run `redot`
to handle your vim settings. Everything vim related will live in that folder. If you lose your freaking mind and decide
to use a different editor then cleaning up is as simple as removing the `vim` folder and re-running `redot`.

## Conventions Rule!
I didn't want to have to edit the "base" files any time I made a change or added something new. To accomplish that, I
adopted and tweaked Zach's convention-based setup:

### Global Conventions
- `bin/`: This directory is added to $PATH and is where to put useful scripts
- `homebrew/Brewfile`: This is a [Homebrew bundle file][homebrew-bundle] that gets executed if you elect to do so when you run `redot`

### Topic Conventions
- `<topic>/path.sh`: Any file with this name will be sourced during `$PATH` setup
- `<topic>/install.sh`: Any file with this name will be executed if you elect to do so when you run `redot`
- `<topic>/*.symlink`: Any file with the `.symlink` extension will be symlinked into your `$HOME` directory
- `<topic>/*.source`: Any file with the `.source` extension will be sourced when you run `redot`
-->

## Installation
macOS:
```
TBD
```
Linux:
```
TBD
```
Windows WSL:
```
TBD
```

## Usage
Run this:

```sh
medkit
```

The `medkit` command will walk you through everything and won't make any change to your computer without asking permission first.

## How to Contribute
Contributions are welcome! Check out [this link][contributing] on how you can help!

## Release History
MEDKIT's full release history can be found [here][changelog].

## Credits
* [Zach Holman][zach-holman-github-url] and his excellent ["dotfiles are meant to be forked" blog post][zach-holman-blog-url]

[zach-holman-github-url]: zach@zachholman.com
[zach-holman-blog-url]: https://zachholman.com/2010/08/dotfiles-are-meant-to-be-forked/

## License

Distributed under the MIT license. See the [LICENSE][license] file for more information.

<!-- Markdown link & img definitions -->
[homebrew-bundle]: https://coderwall.com/p/afmnbq/homebrew-s-new-feature-brewfiles
[changelog]: https://github.com/LeadPipeSoftware/medkit/blob/master/CHANGELOG.md
[authors]: https://github.com/LeadPipeSoftware/medkit/blob/master/AUTHORS.md
[contributing]: https://github.com/LeadPipeSoftware/medkit/blob/master/CONTRIBUTING.md
[security]: https://github.com/LeadPipeSoftware/medkit/blob/master/SECURITY.md
[license]: https://github.com/LeadPipeSoftware/medkit/blob/master/LICENSE
[github-issues-image]: https://img.shields.io/github/issues/badges/shields.svg
[github-issues-url]: https://github.com/LeadPipeSoftware/medkit/issues
[wiki]: https://github.com/LeadPipeSoftware/medkit/wiki