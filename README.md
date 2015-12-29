mackerel-plugin-mogilefs [![Build Status](https://travis-ci.org/hfm/mackerel-plugin-mogilefs.svg?branch=master)](https://travis-ci.org/hfm/mackerel-plugin-mogilefs)
===

MogileFS custom metrics plugin for mackerel.io agent.

Inspired by [`mogilefsd_activity` and `mogilefsd_queries`](http://munin-monitoring.org/browser/munin-contrib/plugins/mogilefs)

Synopsis
---

```sh
mackerel-plugin-mogilefs [-host=<host>] [-port=<port>] [-tempfile=<tempfile>] [-version]
```

```console
$ ./mackerel-plugin-mogilefs -h
Options:

  -H, -host=127.0.0.1                             Host of mogilefsd

  -p, -port=7001                                  Port of mogilefsd

  -t, -tempfile=/tmp/mackerel-plugin-mogilefs     Temp file name

  -v, -version=false                              Print version information and quit.
```

Example of mackerel-agent.conf

```toml
[plugin.metrics.mogilefs]
command = "/path/to/mackerel-plugin-mogilefs"
```
