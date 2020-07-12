# Contributing Guidelines

Cycloid Team is glad to see you contributing to this project ! In this document, we will provide you some guidelines in order to help get your contribution accepted.

## Reporting an issue

### Issues

When you find a bug in InfraMap, it should be reported using [GitHub issues](https://github.com/cycloidio/inframap/issues). Please define key information like your Operating System (OS), InfraMap origin (docker or from source) and finally the version you are using.

### Issue Types

There are 6 types of labels, they can be used for issues or PRs:

- `enhancement`: These track specific feature requests and ideas until they are completed. They can evolve from a `specification` or they can be submitted individually depending on the size.
- `specification`: These track issues with a detailed description, this is like a proposal.
- `bug`: These track bugs with the code
- `docs`: These track problems with the documentation (i.e. missing or incomplete)
- `maintenance`: These tracks problems, update and migration for dependencies / third-party tools
- `refactoring`: These tracks internal improvement with no direct impact on the product
- `need review`: this status must be set when you feel confident with your submission
- `in progress`: some important change has been requested on your submission, so you can toggle from `need review` to `in progress`
- `under discussion`: it's time to take a break, think about this submission and try to figure out how we can implement this or this

## Submit a contribution

### Setup your git repository

If you want to contribute to an existing issue, you can start by _forking_ this repository, then clone your fork on your machine.

```shell
$ git clone https://github.com/<your-username>/inframap.git
$ cd inframap
```

In order to stay updated with the upstream, it's highly recommended to add `cycloidio/inframap` as a remote upstream.

```shell
$ git remote add upstream https://github.com/cycloidio/inframap.git
```

Do not forget to frequently update your fork with the upstream.

```shell
$ git fetch upstream --prune
$ git rebase upstream/master
```

### Play with the codebase

#### Build from sources

Since InfraMap is a Go project, Go must be installed and configured on your machine (really ?). We currently support Go1.14 and go `modules` as dependency manager. You can simply pull all necessaries dependencies by running an initial.

```shell
$ make build
```

This basically builds `inframap` with the current sources.

You also need to install other code dependencies not mandatory in the runtime environment:
  * [enumer](https://github.com/dmarkham/enumer) is used to generate some code.
  * [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) is used to format / organize code and imports. CI will perform a check on this, we highly recommend to run `$ goimports -w <your-modified-files>`.


#### Add a new Provider

To add a new Provider you need to follow the `provider.Provider` interface and create it on the `provider/{provider_name}`, also add it to the `provider.Type` list and then run `make generate`. Basically you can also check any of the other providers that we already have.

To add a new test we recommend to run the `inframap prune --canonicals --tfstate [FILE]` and add the test to the `infraview/testdata/` and add a new test to the `infraview/graph_test.go`

#### Add a new Printer

To add a new Printer just follow the `printer.Printer` interface and add it to `printer/{printer_name}` and `printer/type.go` and run `make generate`. After that it'll be directly added to the list of available `--printer`.
