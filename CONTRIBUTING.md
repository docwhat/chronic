# Contributing

I love pull-requests from everyone!

## Getting Started

### Install Go

First you'll need to make sure you have the latest version. <golang.org> has some [good instructions on installing Go](https://golang.org/doc/install).

### Getting the source

If you will be contributing, then you'll want to [fork the repository](https://help.github.com/articles/fork-a-repo/).

Once you've forked it, then you can clone the source:

```console
git clone git@github.com:<your-username>/<repository-name>.git
```

Fetch the required dependencies:

```console
script/bootstrap
```

Before you do any changes, make sure the tests pass:

```console
script/test
script/lint
```

Make your change. Add tests for your change. Make the tests pass:

```console
script/test
script/lint
```

Push to your fork and [submit a pull request](https://help.github.com/articles/creating-a-pull-request/).

At this point you're waiting on me. I try to be responsive to pull requests, but you know life can get in the way. I may suggest some changes or improvements or alternatives.

Some things will increase the chance that your pull request is accepted:

-   Write tests.
-   Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).

## Why is this repo over-engineered?

Firstly, it is because this is a tool that is likely to be downloaded once and just used forever.  It must not fail unless it is for a darn good reason. It

Secondly, it is because I use this repo to test various GitHub bots, techniques, and CI features.
