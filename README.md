Miroir CLI 
==========

![](https://img.shields.io/github/release/tadashi-aikawa/miroir-cli.svg)

CLI for Miroir. Not support for Windows.


Install
-------

```
$ wget "https://github.com/tadashi-aikawa/miroir-cli/releases/download/v0.4.0/miroir-0.4.0-x86_64-linux.tar.gz" \
    && tar zvfx *.tar.gz --remove-file \
    && chmod +x miroir
```


For developer
-------------

### Release

#### Requirements

* make
* bash
* dep
* ghr

#### Packaging and deploy

Confirm that your branch name equals release version, then...

```
$ make release
```
