# Dogbot

<img src="https://raw.githubusercontent.com/dogtools/dogbot/master/img/dogbot_screenshot.png">

Dogbot is a _ChatOps_ bot based on the [Dog](https://github.com/dogtools/dog) task runner and the [Dogfile Spec](https://github.com/dogtools/dog/blob/master/DOGFILE_SPEC.md).

Unlike similar bots, Dogbot allows you to directly expose shell scripts in your chat room. It's designed to use `/bin/sh` by default but it supports multiple scripting languages. Check the Dog documentation for more details on the supported runners.

This initial version works only in [Slack](https://slack.com/) but other platforms will be implemented eventually.

## Use in chat

List all tasks

    @dogbot list

Run a task

    @dogbot taskname

Ask for help

    @dogbot help

## Configure

Dogbot requires you to provide a Slack API key and a Slack Bot ID.

You can either use the environment variables `DOGBOT_API_KEY` and `DOGBOT_BOT_ID` or provide the values as command-line arguments (`-key` and `-id`).

By default Dogbot looks for a Dogfile in the current path of execution, but you can specify an alternative directory using `-dogfile`.
