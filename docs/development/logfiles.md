---
title: Logfiles
parent: Development
order: 40
---

Rolens keeps track of log-worthy events and logs them in its log directory.

* Windows: `~/AppData/Local/Rolens/Logs`
* Mac: `~/Library/Logs/Rolens`
* Linux: `~/.config/rolens/logs`

In those directories, you can find the following files:

* `rolens.log`: This is the main log file which is a chronological stream of events such as frontend errors, connection problems et cetera.
* `environment.json`: This file contains information about your environment that could be useful while debugging.
