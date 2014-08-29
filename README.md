complete
========

complete - a programatic shell completion library for go

You know shell completion? When you hit TAB and it completes the next part of the command for you?

Yeah, well some shells make that difficult to do access unless you're actually in the shell, especially bash. This library is meant to make this accessible in your go code.

Currently, there's only one function: Bash(). It completes commands using your systems bash configuration.

Currently supports:
* Linux (if your bash completion is at /etc/bash_completion)
* Max OS X (if you installed bash completion support with Homebrew)
