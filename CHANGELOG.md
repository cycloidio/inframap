## [Unreleased]

### Added

- Support for TF State V3 (which is TF 0.11)
  ([Issue #74](https://github.com/cycloidio/inframap/issues/74))

### Changed

- Terraform version from 0.14.6 to 0.14.7
  ([PR #107](https://github.com/cycloidio/inframap/pull/107))

### Fixed

- Improved the mutate logic to consider more data before removing the edges
  ([PR #106](https://github.com/cycloidio/inframap/pull/106))

## [0.5.2] _2021-02-09_

### Changed

- Terraform version from 0.14.5 to 0.14.6
  ([PR #92](https://github.com/cycloidio/inframap/pull/98))

## [0.5.1] _2021-01-22_

### Changed

- Terraform version from 0.14.4 to 0.14.5
  ([PR #92](https://github.com/cycloidio/inframap/pull/92))

## [0.5.0] _2021-01-13_

### Added

- Validation for Terraform State version, we only support 4
  ([Issue #72](https://github.com/cycloidio/inframap/issues/72))

### Changed

- Terraform version from 0.14.3 to 0.14.4
  ([PR #88](https://github.com/cycloidio/inframap/pull/88))

## [0.4.0] _2020-12-09_

### Added

- Support incoming connection without source node for AWS
  ([Issue #5](https://github.com/cycloidio/inframap/issues/5))

### Changed

- `tfdocs` version upgraded
  ([PR #69](https://github.com/cycloidio/inframap/pull/69))
- Terraform version from 0.13.5 to 0.14.2
  ([PR #75](https://github.com/cycloidio/inframap/pull/75))

### Fixed

- Azure not generating a correct tfstate due to renamed method
  ([PR #71](https://github.com/cycloidio/inframap/pull/71))

## [0.3.3] _2020-10-22_

### Changed

- Terraform version from 0.13.4 to 0.13.5
  ([PR #65](https://github.com/cycloidio/inframap/pull/65))

## [0.3.2] _2020-10-01_

### Changed

- Terraform version from 0.13.3 to 0.13.4
  ([PR #61](https://github.com/cycloidio/inframap/pull/61))

## [0.3.1] _2020-09-21_

### Changed

- Terraform version from 0.13.0 to 0.13.3
  ([Issue #58](https://github.com/cycloidio/inframap/issues/58))
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
