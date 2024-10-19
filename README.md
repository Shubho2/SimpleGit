# Introduction
A simple Git CLI tool to understand basic Git commands. It supports low-level commands like init, ls-tree, hash-tree, cat-file, commit-tree.

## How to build and run this project

1. Ensure you have `go` installed locally
2. Run `./your_program.sh` to run your Git implementation, which is implemented
   in `cmd/mygit/main.go`.

The `your_program.sh` script is expected to operate on the `.git` folder inside
the current working directory. If you're running this inside the root of this
repository, you might end up accidentally damaging your repository's `.git`
folder.

We suggest executing `your_program.sh` in a different folder when testing. For example:

```sh
mkdir -p /tmp/testing && cd /tmp/testing
/path/to/your/repo/your_program.sh init
```
To make this easier to type out, you could add a
[shell alias](https://shapeshed.com/unix-alias/):

```sh
alias mygit=/path/to/your/repo/your_program.sh

mkdir -p /tmp/testing && cd /tmp/testing
mygit init
```
*Note: write-tree assumes that all files in the working directory is staged. Since, commit-tree doesn't consider the case of merged branch, it only takes one parent commit_sha as argument.*