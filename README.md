# InfraMap

Read your tfstate or HCL to generate a graph specific for each provider, showing only the
resources that are most important/relevant.

## Cloud Providers

We support certain providers. This allows us to better represent information that comes from these providers.

If your state file or HCL is from a provider we do not support, the resulting representation will simply be all resources present without any simplification or refinement.

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

The `inframap --help` will show you the basics.

[![asciicast](https://asciinema.org/a/347600.svg)](https://asciinema.org/a/347600)

The most important subcommands are:

* `generate`: generates the graph from STDIN or file.
* `prune`: removes all unecessary information from the state or HCL (not supported yet) so it can be shared without any security concerns

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

For each provider, we support specific types of connections; we have a static list of resources that can be
nodes or edges. Once we identify the edges, we try to create one unique edge from the resources they connect.

For a state file, we rely on the `depends_on` key and, for HCL we rely on interpolation to create the base graph one which we then
apply specific provider logic if supported. If not supported, then basic graph is returned.

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
