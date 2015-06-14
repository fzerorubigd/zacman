# zacman

zacman, a simple zsh package manager in go

## Goal

I need a simple way to use zsh plugins and theme, but I don't want to use a lot of blather just for some simple functions that I need. something like antigen, but with a static load mechanism like antigen-hs and no Haskel or other runtime, so I choose Go.

## How to build?
### gb
This is a [gb](http://getgb.io/) project, so you can simply cline this, and use gb to build projects.

```
git clone --recursive  https://github.com/fzerorubigd/zacman.git

gb build
```

Note : I use git submodules instead of gb vendor.
### go tools
but if you want, you can build it with `go get` (not recomanded)

```
go get github.com/fzerorubigd/zacman/src/zacman
```

## Usage

### bundle
first install some bundles. use `zacman bundle` command :

```
zacman bundle zsh-users/zsh-syntax-highlighting
```

if the plugin files are not in the root folder use the secound parameter for sub path :

```
zacman bundle zsh-users/zsh-completions src
```

see `zacman help bundle` for detail help.

### list
then check for list of installed bundle :

```
zacman list
```

### remove
for removing a bundle simply use
```
zacman remove sharat87/autoenv
```

Note : if bundle contain a sub path, in remove should supply that sub path too :
```
zacman remove zsh-users/zsh-completions src
```

### snapshots

you can save snapshots from your current bundles and restore them when you want to.

```
zacman snapshot A_NAME

zacman restore A_NAME
```

### Final touch : Compile

After any change to bundles, you need to compile things.
```
zacman compile
```

## .zshrc

begin your .zshrc with this :

```
if [ ! -e ~/.zacman/zacman.zsh ];then
	zacman compile
fi;
source ~/.zacman/zacman.zsh

# other configuration ....
```
