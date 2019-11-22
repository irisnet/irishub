# Contributing

Thank you for considering making contributions to irishub!

Contributing to this repo can mean many things such as participated in
discussion or proposing code changes. To ensure a smooth workflow for all
contributors, the general procedure for contributing has been established:

  1. either [open](https://github.com/irisnet/irishub/issues/new) or
     [find](https://github.com/irisnet/irishub/issues) an issue you'd like to help with,
  2. participate in thoughtful discussion on that issue,
  3. if you would then like to contribute code:
     1. if the issue is a proposal, ensure that the proposal has been accepted,
     2. ensure that nobody else has already begun working on this issue, if they have
       make sure to contact them to collaborate,
     3. if nobody has been assigned the issue and you would like to work on it
       make a comment on the issue to inform the community of your intentions
       to begin work,
     4. follow standard github best practices: fork the repo,
       if the issue if a bug fix, branch from the
       tip of `develop`, make some commits, and submit a PR to `develop`; if the issue is a new feature,  branch from the tip of `feature/XXX`, make some commits, and submit a PR to `feature/XXX`
     5. include `WIP:` in the PR-title to and submit your PR early, even if it's
       incomplete, this indicates to the community you're working on something and
       allows them to provide comments early in the development process. When the code
       is complete it can be marked as ready-for-review by replacing `WIP:` with
       `R4R:` in the PR-title.

Note that for very small or blatantly obvious problems (such as typos) it is 
not required to an open issue to submit a PR, but be aware that for more complex
problems/features, if a PR is opened before an adequate design discussion has
taken place in a github issue, that PR runs a high likelihood of being rejected. 

Note, we use `make
get_dev_tools` and `make update_dev_tools` for installing the linting tools.Please make sure to use `gofmt` before every commit - the easiest way to do this is have your editor run it for you upon saving a file.

## Pull Requests

To accommodate review process we suggest that PRs are categorically broken up.
Ideally each PR addresses only a single issue. Additionally, as much as possible
code refactoring and cleanup should be submitted as a separate PRs. And the feature branch `feature/XXX` should be synced with `develop` regularly.

## Dependencies

We use [dep](https://github.com/golang/dep) to manage dependencies.

Since some dependencies are not under our control, a third party may break our
build, in which case we can fall back on `dep ensure` (or `make
get_vendor_deps`).

## Testing

The `Makefile` defines `make test` and includes its continuous integration. For any new comming feature, the `test_unit` / `test_cli` and `test_lcd` must be provided.

We expect tests to use `require` or `assert` rather than `t.Skip` or `t.Fail`,unless there is a reason to do otherwise.

### PR Targeting

Ensure that you base and target your PR on the correct branch:

- `release/vxx.yy` for a merge into a release candidate
- `master` for a merge of a release
- `develop` in the usual case

All feature additions should be targeted against `feature/XXX`. Bug fixes for an outstanding release candidate
should be targeted against the release candidate branch. Release candidate branches themselves should be the
only pull requests targeted directly against master.

### Development Procedure

- the latest state of development is on `develop`
- `develop` must never fail `make test`
- no --force onto `develop` (except when reverting a broken commit, which should seldom happen)

### Pull Merge Procedure

- ensure `feature/XXX` is rebased on `develop`
- ensure pull branch is rebased on `feature/XXX`
- run `make test` to ensure that all tests pass
- merge pull request
- push `feature/XXX` into `develop` regularly

### Release Procedure

- start on `develop`
- prepare changelog/release issue
- bump versions
- push to `release-vX.X.X` to run CI
- merge to master
- merge master back to develop

### Hotfix Procedure

- start on `release-vX.X.X`
- make the required changes
  - these changes should be small and an absolute necessity
  - add a note to CHANGELOG.md
- bump versions
- merge `release-vX.X.X` to master if necessary
- merge `release-vX.X.X` to develop if necessary
