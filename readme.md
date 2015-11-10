# [Sentinel][releases]

[![Circle CI](https://circleci.com/gh/hackliff/sentinel.svg?style=svg)](https://circleci.com/gh/hackliff/sentinel)

<h1 align="center">
  <br>
  <img width="400" src="sentinels.jpg">
  <br>
  <br>
</h1>

> Only the paranoid will survice

Modern applications are mostly a composition of services [that will eventually fail][chaos].

__Sentinel__ is a framework to monitor and handle distributed infrastructures,
like modern microservice topologies or complex applications.

__Sentinels__ are composable bots you can launch to monitor your things, and forget
until something goes wrong.

It is built out of of three components:

- Configuration: declarative definition of stuff to monitor
- Plugins:
  - Sensors: pluggable drivers to perform checks
  - Radios: pluggable services for alerting
  - Triggers: orchestrate sensors' measures
- Core library glueing bullet points above in a neat lightweight cli package

It's written in _Go_ to keep it super easy to deploy and elegant to implement
multiple backends through interfaces.

The project is just getting started, don't expect anything to work at this
point but contribution is welcome !


## Deployment

```Sh
local version=0.2.0
local platform=darwin-amd64

curl \
  -ksL \
  -o /usr/local/bin/sentinel \
  https://github.com/hackliff/sentinel/releases/download/${version}/sentinel-${platform}
chmod +x /usr/local/bin/sentinel
```


## [Documentation][doc]

Doc is [available online][doc] or you can build it locally :

```Sh
make doc
open site/index.html
```

### Go Package

Check it out on [gowalker][walker], [godoc][godoc], or browse it locally:

```console
$ make godoc
$ $BROWSER docker-dev:6060/pkg/github.com/hackliff/sentinel
```


## Contributing

> Fork, implement, add tests, pull request, get my everlasting thanks and a
> respectable place here [=)][jondotquote]

```console
make build
make tests TESTARGS=-v
```


## Conventions

__sentinel__ follows some wide-accepted guidelines

* [Semantic Versioning known as SemVer][semver]
* [Git commit messages][commit]


## Authors

| Selfie               | Name            | Twitter                     |
|----------------------|-----------------|-----------------------------|
| <img src="https://avatars.githubusercontent.com/u/1517057" alt="text" width="40px"/> | Xavier Bruhiere | [@XavierBruhiere][xbtwitter] |


## Licence

Copyright 2015 Xavier Bruhiere.

__Sentinel__ is available under the MIT Licence.


---------------------------------------------------------------


<p align="center">
  <img src="https://raw.github.com/hivetech/hivetech.github.io/master/images/pilotgopher.jpg" alt="gopher" width="200px"/>
</p>


[releases]: https://github.com/hackliff/sentinel/releases
[semver]: http://semver.org
[commit]: https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit#
[xbtwitter]: https://twitter.com/XavierBruhiere
[jondotquote]: https://github.com/jondot/groundcontrol
[walker]: http://gowalker.org/github.com/hackliff/sentinel
[godoc]: http://godoc.org/github.com/hackliff/sentinel
[doc]: http://hackliff.github.io/sentinel/
[chaos]: http://techblog.netflix.com/2012/07/chaos-monkey-released-into-wild.html
