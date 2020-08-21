# InfraMap

[![PkgGoDev](https://pkg.go.dev/badge/github.com/cycloidio/inframap)](https://pkg.go.dev/github.com/cycloidio/inframap)
[![AUR package](https://repology.org/badge/version-for-repo/aur/inframap.svg)](https://repology.org/project/inframap/versions)

Read your tfstate or HCL to generate a graph specific for each provider, showing only the
resources that are most important/relevant.

<p align="center">
  <img src="docs/inframap.png" width="400">
</p>

## Cloud Providers

We support certain providers. This allows us to better represent information that comes from these providers.

If your state file or HCL is from a provider we do not support, the resulting representation will simply be all resources present without any simplification or refinement.

We support Terraform: 0.13 (and previous)

| Provider | State | HCL |  Grouping<sup>1</sup> | External Nodes<sup>2</sup> | IAM<sup>3</sup> |
|:--:|:--:|:--:|:--:|:--:|:--:|
| <img alt="AWS" src="docs/aws.png" width="50"> | :heavy_check_mark: | :heavy_check_mark: | [WIP](https://github.com/cycloidio/inframap/issues/6) | [WIP](https://github.com/cycloidio/inframap/pull/41) | :heavy_multiplication_x: (https://github.com/cycloidio/inframap/issues/11)|
| <img alt="Google" src="docs/google-cloud.svg" width="50"> | :heavy_check_mark: | :heavy_check_mark: | :heavy_multiplication_x: | :heavy_multiplication_x: | :heavy_multiplication_x: |
| <img alt="Azure" src="docs/azure.svg" width="50"> | :heavy_check_mark: | :heavy_check_mark: | :heavy_multiplication_x: | :heavy_multiplication_x: | :heavy_multiplication_x: |
| <img alt="OpenStack" src="docs/Openstack-vertical-small.png" width="50"> | :heavy_check_mark: | :heavy_check_mark: | :heavy_multiplication_x: | :heavy_multiplication_x: | :heavy_multiplication_x: |
| <img alt="FlexibleEngine" src="docs/flexibleengine.png" width="50"> | :heavy_check_mark: | :heavy_check_mark: | :heavy_multiplication_x: | :heavy_multiplication_x: | :heavy_multiplication_x: |

1. **Grouping**: Group elements that belong to the same group like VPCs or regions
2. **External Nodes**: Show the ingress of the Nodes if any
3. **IAM**: Connections based on IAM (Identity Access Management)

## Installation

### Stable

To install the latest release of Inframap, you can pick one of this methods:
  * pull the latest release from the [Releases](https://github.com/cycloidio/inframap/releases/) page
  * pull the latest docker [image](https://hub.docker.com/r/cycloid/inframap) from the Docker hub
  * use your Linux package manager (only [AUR](https://aur.archlinux.org/packages/inframap) at the moment)

### Development

You can build and install with the latest sources, you will enjoy the new features and bug fixes. It uses Go Modules (1.13+)

```shell
$ git clone https://github.com/cycloidio/inframap
$ cd inframap
$ go mod download
$ make build
```

## Usage

The `inframap --help` will show you the basics.

[![asciicast](https://asciinema.org/a/347600.svg)](https://asciinema.org/a/347600)

The most important subcommands are:

* `generate`: generates the graph from STDIN or file.
* `prune`: removes all unnecessary information from the state or HCL (not supported yet) so it can be shared without any security concerns

### Example

Visualizing with dot

```shell
$ inframap generate --tfstate state.tfstate | dot -Tpng > graph.png
```

or from the terminal itself

```shell
$ inframap generate --tfstate state.tfstate | graph-easy
```

or from HCL

```shell
$ inframap generate --hcl config.tf | graph-easy
```

or HCL module

```shell
$ inframap generate --hcl ./my-module/ | graph-easy
```

using docker image (assuming that your Terraform files are in the working directory)

```shell
$ docker run --rm -v ${PWD}:/opt cycloid/inframap generate --tfstate /opt/terraform.tfstate
```

## How is it different to `terraform graph`

[Terraform Graph](https://www.terraform.io/docs/commands/graph.html) outputs a dependency graph of all the resources on the tfstate/HCL. We try to go one step further,
by trying to make it human-readable.

If the provider is not supported, the output will be closer to the Terraform Graph version (without displaying provider / variable nodes)

Taking https://github.com/cycloid-community-catalog/stack-magento/ as a reference this is the difference in output:

With `terraform graph`:

<p align="center">
  <img src="docs/terraformgraph.svg" width="400">
</p>

With `inframap generate --hcl ./terraform/module-magento/ | dot -Tpng > inframap.png`:

<p align="center">
  <img src="docs/inframap.png" width="400">
</p>

With `inframap generate --hcl --connections=false ./terraform/module-magento/ | dot -Tpng > inframapconnections.png`:

<p align="center">
  <img src="docs/inframapconnections.png" width="400">
</p>

With `inframap generate --hcl ./terraform/module-magento/ --raw | dot -Tpng > inframapraw.png`:

<p align="center">
  <img src="docs/inframapraw.png" width="400">
</p>

## How does it work?

For each provider, we support specific types of connections; we have a static list of resources that can be
nodes or edges. Once we identify the edges, we try to create one unique edge from the resources they connect.

For a state file, we rely on the `dependencies` key (for the <0.13 we replace all `depends_on` for `dependencies` so we support them) and, for HCL we rely on interpolation to create the base graph one which we then
apply specific provider logic if supported. If not supported, then basic graph is returned.

## FAQ

### Why is my Graph generated empty?

If a graph is returned empty, it means that we support one of the providers you are using on your HCL/TFState but we do
not recognize any connection or relevant node.

To show the configuration without any InfraMap applied logic you can use the `--raw` flag logic and print everything that we read.
If it works, it would be good to try to know why it was empty before so we can take a look
at it as it could potentially be an issue on InfraMap (open an issue if you want us to take a look).

By default unconnected nodes are removed, you can use `--clean=false` to prevent that.

### Does InfraMap support Terraform backends ?

Terraform allows users to use `backends` (S3, Google Cloud Storage, Swift, etc.) in order to store the `terraform.state`. We currently do not support graph generation from `tfstate` stored in such backends. As mentioned in this [issue](https://github.com/cycloidio/inframap/issues/44), it is possible to play around `stdin/out` to generate graph from Terraform backends.

| backend | command                                                                          |
|---------|----------------------------------------------------------------------------------|
| S3      | `aws s3 cp s3://bucket/path/to/your/file.tfstate - \| inframap generate --tfstate` |
| GCS     | `gsutil cat gs://bucket/path/to/your/file.tfstate \| inframap generate --tfstate`  |

### What does the error `Error: error while reading TFState: state snapshot was created by Terraform v0.13.0, which is newer than current v0.12.28; upgrade to Terraform v0.13.0 or greater to work with this state` mean?

We use Terraform internally to read the TFStates and it has a built-in validation that checks that the version of Terraform used (specified as `"version": "X.Y.Z",` on the state) is the same or older than the one defined.

In this case it means that Terraform released a new version and InfraMap did not yet update to that version so we have to update to the latest version of Terraform. It'll not take more than a few days but if you update on the same day as the Terraform release then InfraMap will fail due to that.

You can open an issue to notify us just in case we missed it, similar to https://github.com/cycloidio/inframap/issues/47

## License

Please see the [MIT LICENSE](https://github.com/cycloidio/inframap/blob/master/LICENSE) file.

## Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## About Cycloid

<p align="center">
  <img src="https://user-images.githubusercontent.com/393324/65147266-0b010100-da1e-11e9-9a49-d27e5035c4c4.png">
</p>

[Cycloid](https://www.cycloid.io/our-culture) is a European fully-remote company, building a product to **simplify**, **accelerate** and **optimize your DevOps and Cloud adoption**.

We built [Cycloid, your DevOps framework](https://www.cycloid.io/devops-framework) to encourage Developers and Ops to work together with the respect of best practices. We want to provide a tool that eliminates the silo effect in a company and allows to share the same level of informations within all professions.

[Cycloid](https://www.cycloid.io/devops-framework) supports you to factorize your application in a reproducible way, to deploy a new environment in one click. This is what we call a stack.

A stack is composed of 3 pillars:

1. the pipeline ([Concourse](https://concourse-ci.org/))
2. infrastructure layer ([Terraform](https://www.terraform.io/))
3. applicative layer ([Ansible](https://www.ansible.com/))

Thanks to the flexible pipeline, all the steps and technologies are configurable.

To make it easier to create a stack, we build an Infrastructure designer named **StackCraft** that allows you to drag & drop Terraform resources and generate your Terraform files for you.

InfraMap is a brick that will help us to visualize running infrastructures.

The product comes also with an Open Source service catalog ([all our public stacks are on Github](https://github.com/cycloid-community-catalog)) to deploy applications seamlessly.
To manage the whole life cycle of an application, it also integrates the diagram of the infrastructure and the application, a cost management control to centralize Cloud billing, the monitoring, logs and events centralized with Prometheus, Grafana, ELK.

[Don't hesitate to contact us, we'll be happy to meet you !](https://www.cycloid.io/contact-us)
