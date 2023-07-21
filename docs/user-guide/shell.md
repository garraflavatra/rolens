---
title: Shell
summary: "Write and execute MongoDB shell scripts within Rolens"
parent: User guide
order: 60
stub: true
---

Rolens has a shell feature: it provides an editor for writing shell scripts and executing them against a local or external host, database, or even a single collection. You can find it under the _Shell_ tab.

_Note: this feature is currently in development and will be shipped with version 0.3.0. You can [follow the development on GitHub](https://github.com/garraflavatra/rolens)._

![The shell tab](/images/shell.png)

## Requirements

To use the script editor, you need to install the official [`mongosh` tool](https://www.mongodb.com/docs/mongodb-shell/) from MongoDB. You can [download it](https://www.mongodb.com/try/download/shell) from the MongoDB web site, or run `brew install mongosh` if you use a Mac.
