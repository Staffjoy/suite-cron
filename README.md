# Suite Cron

[![Moonlight contractors](https://img.shields.io/badge/contractors-1147-brightgreen.svg)](https://moonlightwork.com/for/staffjoy)

[Staffjoy is shutting down](https://blog.staffjoy.com/staffjoy-is-shutting-down-39f7b5d66ef6#.ldsdqb1kp), so we are open-sourcing our code. This repository is a microservice for [Staffjoy V1, aka "Suite"](https://github.com/staffjoy/suite). It performs an authenticated poll of the cron endpoint, and it reports the results to [Papertrail](https://papertrailapp.com).

## Credit

The author of this repo in its entirety is [@philipithomas](https://github.com/philipithomas). This is a fork of the internal repository. For security purposes, the Git history has been squashed. The `vendor` folder includes a [Papertrail library by @zemirco](https://github.com/zemirco/papertrail) under an MIT license. 
