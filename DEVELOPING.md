# Developing Bombadillo

## Getting Started

Following the standard install instructions should lead you to have nearly everything you need to commence development. The only additions to this are:

- To be able to submit pull requests, you will need to fork this repository first.
- The build process must be tested with Go 1.11 to ensure backward compatibility. This version can be installed as per the [Go install documentation](https://golang.org/doc/install#extra_versions). Check that changes build with this version using `make test`.
- Linting must be performed on new changes using `gofmt` and [golangci-lint](https://github.com/golangci/golangci-lint)


## How changes are made

A stable version of Bombadillo is kept in the default branch, so that people can easily clone the repo and get a good version of the software.

New changes are introduced to the **develop** branch.

Changes to the default branch occur as part of the software release process. This usually occurs when:

 - there are a set of changes in **develop** that are good enough to be considered stable.
 - an urgent issue is identified in the stable version that requires immediate changes


### Process for introducing a new change

Please refer to our [notes on contributing](README.md#contributing) to get an understanding of how new changes are initiated, the type of changes accepted and the review process.

1. Create a new feature branch based on the **develop** branch.
1. Raise a pull request (PR) targeting the **develop** branch.
1. The PR is reviewed.
1. If the PR is approved, it is merged.
1. The version patch number is incremented.


### Incrementing the version number

This process is handled by maintainers after a change has been merged.

Version numbers are comprised of three digits: major version number, minor version number, and patch number.

The version number is incremented in the following situations:

#### New changes

Each new change added to **develop** should increment the patch number. For example, version 2.0.1 would become 2.0.2. After the change is merged from the feature branch to **develop**:

```shell
# ensure everything is up to date and in the right place
git checkout develop
git pull

# get the commit ID for the recent merge
git log

# get the current version number (the highest number)
git tag

# add the incremented version number to the commit-id, for example:
git tag 2.0.2 abcdef
```

#### Release process

As part of the software release process, any part of the version number may change:

- Urgent changes increment the **patch** number
- A set of small changes increments the **minor** version number
- A significant change to large parts of the application increments the **major** version number

1. The version number in the VERSION file is incremented. This change is committed to the default branch.
1. The **develop** branch is merged to the default branch.
1. The commands from the New Changes section are followed, but the final command includes an annotation: `git tag 2.1.0 abdef -a "This version adds several new features..."`
1. Release information should also be added to the [tildegit releases page](https://tildegit.org/sloum/bombadillo/releases).
