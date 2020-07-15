# InfraMap

Read your TF State or HCL to generate a Graph specific for each Provider, only showing the
Resources that are most important/relevant to see.

## Cloud Providers

We support specific implementations of specific Providers that basically allow us to have
a better visual representation on those Providers.

If the State or HCL provided is from a Provider we do not support, then the result
will be all the resources an how they are connected without making any simplification.

Support:

| Provider | State | HCL |
|--|:--:|:--:|
| AWS | Yes | Yes |
| FlexibleEngine | Yes | Yes |
| OpenStack | Yes | Yes |

## Installation

### Development

You can build and install with the latest sources, you will enjoy the new features and bug fixes. It uses Go Modules

```shell
$ git clone https://github.com/cycloidio/inframap
$ cd inframap
$ go mod download
$ make build
```

## Usage

Using the `inframap --help` you will know the basics.

[![asciicast](https://asciinema.org/a/347600.svg)](https://asciinema.org/a/347600)

The important subcommands are:

* `generate`: Which generates the Graph from Stdin or File.
* `prune`: Which removes all the not needed information from the State or HCL (not supported yet) so it can be shared without any security concern

### Example

Visualizing with dot

```shell
$ inframap generate --tfstate state.json | dot -Tsvg > graph.svg
```

or from the terminal itself

```shell
$ inframap generate --tfstate state.json | graph-easy
```

or from HCL

```shell
$ inframap generate --hcl config.tf | graph-easy
```

or HCL module

```shell
$ inframap generate --hcl ./my-module/ | graph-easy
```

## How does it work?

For each of those Providers we support specific types of connections, we have a static list of Resources that can be
Nodes or Edges. Once we identify the edges we try to create an unique edge from the Resources that they connect.

This is based on the `depends_on` on the State and on Interpolation on the HCL to create the base graph in which then
we apply specific Provider logic if supported. If not supported then that basic graph is the one returned

### AWS

* `aws_security_group`
* `aws_security_group_rule`

### FlexibleEngine

* `flexibleengine_compute_interface_attach_v2`
* `flexibleengine_networking_port_v2`
* `flexibleengine_networking_secgroup_rule_v2`
* `flexibleengine_networking_secgroup_v2`
* `flexibleengine_lb_listener_v2`
* `flexibleengine_lb_pool_v2`
* `flexibleengine_lb_member_v2`

### OpenStack

* `openstack_compute_interface_attach_v2`
* `openstack_networking_port_v2`
* `openstack_networking_secgroup_rule_v2`
* `openstack_networking_secgroup_v2`
* `openstack_lb_listener_v2`
* `openstack_lb_pool_v2`
* `openstack_lb_member_v2`

## License

TODO

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
