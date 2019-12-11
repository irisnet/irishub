# modules
IRIS Hub Modules

**Note**: This repository is meant to house modules that are created outside of the [IRIS Hub](https://github.com/irisnet/irishub) repository.

**Note**: Requires [Go 1.13+](https://golang.org/dl/)

## Quick Start

This repo organizes modules into 3 subfolders:

- `stable/`: this folder houses modules that are stable, production-ready, and well-maintained.
- `incubator/`: this folder houses modules that are buildable but makes no guarantees on stability or production-readiness. Once a module meets all requirements specified in [contributing guidelines](./CONTRIBUTING.md), the owners can make a PR to move module into `stable/` folder. Must be approved by at least one `modules` maintainer for the module to be moved.
- `inactive/`: Any stale module from the previous 2 folders may be moved to the `inactive` folder if it is no longer being maintained by its owners or the module has been accepted into the [Cosmos modules](https://github.com/cosmos/modules) repository. `modules` maintainers reserve the right to move a module into this folder after public discussion in an issue and a specified grace period for module owners to restart work on module.

Any changes to where modules are located will only happen on major releases of the `modules` repo to ensure we only break import paths on major releases.
