# git-publish

[![Build Status](https://travis-ci.org/wyattjoh/git-publish.svg?branch=master)](https://travis-ci.org/wyattjoh/git-publish)

This is a small tool to help out with frequent work with extra branches.

It essentially mirrors the `git push --set-upstream origin $current_branch`.

There's not much to it then that.

## Installation

Install via with the Go toolchain to compile from source:

```bash
go get github.com/wyattjoh/git-publish
```

Download pre-compiled binary on the [Releases Page](https://github.com/wyattjoh/git-publish/releases/latest) for your Arch/OS.

### Installation Via Homebrew

brew install wyattjoh/stable/git-publish

## Usage

```bash
git checkout -b some-branch
# do some work... commit it
# instead of doing `git push -u origin some-branch`, do:
git publish
```

You may specify another remote with `git publish <remote name>`, the default
however is `origin`.
