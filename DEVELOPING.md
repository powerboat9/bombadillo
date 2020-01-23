# Developing Bombadillo

## Getting Started

Following the standard install instructions should lead you to have nearly everything you need to commence development. The only additions to this are:

- To be able to submit pull requests, you will need to fork this repository first.
- The build process must be tested with Go 1.11 to ensure backward compatibility. This version can be installed as per the [Go install documentation](https://golang.org/doc/install#extra_versions). Check that changes build with this version using `make test`.
- Linting must be performed on new changes using `gofmt` and [golangci-lint](https://github.com/golangci/golangci-lint)


## How changes are made

A stable version of Bombadillo is kept in the default branch, so that people can easily clone the repo and get a good version of the software.

New changes are implemented to the **develop** branch as **development releases**.

Changes are implemented to the default branch when:

 - There are a set of changes in **develop** that are good enough to be considered stable.
   - This may be a **minor** set of changes for a **minor release**, or
   - a large **major** change for **major release**.
 - An urgent issue is identified in the stable version that requires an immediate **patch release**.


### Process for introducing a new change

Please refer to our [notes on contributing](README.md#contributing) to get an understanding of how new changes are initiated, the type of changes accepted and the review process.

1. Create a new feature branch based on the **develop** branch.
1. Raise a pull request (PR) targeting the **develop** branch.
1. The PR is reviewed.
1. If the PR is approved, it is merged.
1. The version number is incremented, along with any other release activity.


### Process for incrementing the version number

The version number is incremented during a **development release**, **patch release**, and **minor** and **major releases**. This is primarily managed through git tags in the following way: 

```shell
# switch to the branch the release is being performed for
git checkout branch

# ensure everything is up to date
git pull

# get the commit ID for the recent merge
git log

# get the current version number (the highest number)
git tag

# for a development release, add the incremented version number to the commit-id, for example:
git tag 2.0.2 abcdef

# for releases to the default branch, this tag can also be added with annotations
git tag 2.1.0 abdef -a "This version adds several new features..."
```

Releases to the default branch also include the following tasks:

1. The version number in the VERSION file is incremented and committed.
1. Release information should also be verified on the [tildegit releases page](https://tildegit.org/sloum/bombadillo/releases).
