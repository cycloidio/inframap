## [Unreleased]

### Changed

- Terraform version from 0.13.0 to 0.13.3
  ([PR #55](https://github.com/cycloidio/inframap/pull/55))
- CONTRIBUTING to include Architecture description
  ([Issue #52](https://github.com/cycloidio/inframap/issues/52))


## [0.3.0] _2020-08-21_

### Added

- Capability to have 2 Nodes connected with 2 edges of different directions
  ([PR #38](https://github.com/cycloidio/inframap/pull/38))
- Azure support
  ([Issue #8](https://github.com/cycloidio/inframap/issues/8))
- Flexible Engine icons
  ([PR #45](https://github.com/cycloidio/inframap/pull/45))

### Changed

- Terraform version from 0.12.28 to 0.13
  ([Issue #47](https://github.com/cycloidio/inframap/issues/47))

### Fixed

- Google graph generation from HCL
  ([PR #34](https://github.com/cycloidio/inframap/pull/34))
- Generation error when multiple Edges hanging (not merged)
  ([PR #33](https://github.com/cycloidio/inframap/pull/33))
- Padding between the image and the label for the `dot` printer
  ([PR #42](https://github.com/cycloidio/inframap/pull/42))

## [0.2.0] _2020-07-27_

### Added

- New flag to `generate`, `--connections` to apply or not the Provider logic of merging Edges between Nodes
  ([PR #23](https://github.com/cycloidio/inframap/pull/23))
- Graph generation with Icons
  ([Issue #13](https://github.com/cycloidio/inframap/issues/13))
- Google graph generation from TFState
  ([Issue #7](https://github.com/cycloidio/inframap/issues/7))
- Google graph generation from HCL
  ([Issue #27](https://github.com/cycloidio/inframap/issues/27))

### Fixed

- HCL generation errors
  ([Issue #29](https://github.com/cycloidio/inframap/issues/29))

## [0.1.1] _2020-07-16_

### Added

- Difference between `terraform graph` and InfraMap to the README
  ([PR #14](https://github.com/cycloidio/inframap/pull/14))
- CI/CD configuration and Dockerfile
  ([PR #15](https://github.com/cycloidio/inframap/pull/15))

## [0.1.0] _2020-07-16_

First version and first implementation
