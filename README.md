# antigo

Antigo, a simple zsh package manager in go

## Goal

I need a simple way to use zsh plugins and theme, but I don't want to use a lot of blather just for some simple functions that I need. something like antigen, but with a static load mechanism like antigen-hs and no Haskel or other runtime, so I choose Go.

## How to build?
### gb
This is a [gb](http://getgb.io/) project, so you can simply cline this, and use gb to build projects.

```
git clone --recursive  https://github.com/fzerorubigd/antigo.git

gb build
```

Note : I use git submodules instead of gb vendor.
### go tools
but if you want, you can build it with `go get` (not recomanded)

```
go get github.com/fzerorubigd/antigo/src/antigo
```

## Usage

### bundle
first install some bundles. use `antigo bundle` command :

```
antigo bundle zsh-users/zsh-syntax-highlighting
```

if the plugin files are not in the root folder use the secound parameter for sub path :

```
antigo bundle zsh-users/zsh-completions src
```

see `antigo help bundle` for detail help.

### list
then check for list of installed bundle :

```
antigo list
```

### remove
for removing a bundle simply use
```
antigo remove sharat87/autoenv
```

Note : if bundle contain a sub path, in remove should supply that sub path too :
```
antigo remove zsh-users/zsh-completions src
```

### snapshots

you can save snapshots from your current bundles and restore them when you want to.

```
antigo snapshot A_NAME

antigo restore A_NAME
```

### Final touch : Compile

After any change to bundles, you need to compile things.
```
antigo compile
```

## .zshrc

begin your .zshrc with this :

```
if [ ! -e ~/.antigo/antigo.zsh ];then
	antigo compile
fi;
source ~/.antigo/antigo.zsh

# other configuration ....
```
