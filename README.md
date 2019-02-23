Miroir CLI 
==========

![](https://img.shields.io/github/release/tadashi-aikawa/miroir-cli.svg)

CLI for Miroir. Not support for Windows.


Install
-------

Download from release pages and unarchive.

* ex: `https://github.com/tadashi-aikawa/miroir-cli/releases/download/v0.4.0/miroir-0.4.0-x86_64-linux.tar.gz`


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

You have to create PR and merge.


#### deploy

After you merged PR, then

```
$ make deploy version=x.y.z
```
