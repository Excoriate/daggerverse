# Contribution Guidelines

First and foremost, we'd like to express our gratitude to you for taking the time to contribute.
We welcome and appreciate any and all contributions via
[Pull Requests] along the [GitHub Flow].

1. [Open a GitHub issue](#open-a-github-issue)
2. [Fork the repository on GitHub](#fork-the-repository-on-github)
3. [Install the pre-commit hooks](#install-the-pre-commit-hooks)
4. [Update the documentation](#update-the-documentation)
5. [Update the tests](#update-the-tests)
6. [Update the code](#update-the-code)
7. [Create a pull request](#create-a-pull-request)
8. [Merge and release](#merge-and-release)

## Open a GitHub issue

For bug reports or requests, please submit your issue in the appropriate repository.

We advise that you open an issue and ask the
[CODEOWNERS](.github/CODEOWNERS) and community prior to starting a contribution.
This is your chance to ask questions and receive feedback before
writing (potentially wrong) code. We value the direct contact with our community
a lot, so don't hesitate to ask any questions.

## Fork the repository on GitHub

[Fork] the repository into your own GitHub account and [create a new branch] as
described in the [GitHub Flow].

## Install the pre-commit hooks

If the repository you're working on ships with a
[`.pre-commit-config.yaml`][pre-commit-file],
make sure the necessary hooks have been installed before you begin working
(e.g. a `pre-commit install`).

## Update the documentation

We encourage you to update the documentation before writing any code (please see
[Readme Driven Development](https://tom.preston-werner.com/2010/08/23/readme-driven-development.html). This ensures the
documentation stays up to date and allows you to think through the problem fully before you begin implementing any
changes.

## Update the tests

We also recommend updating the automated tests before updating any code
(see [Test Driven Development](https://en.wikipedia.org/wiki/Test-driven_development).

That means that you should add or update a test case, run all tests and verify
that the new test fails with a clear error message and then start implementing
the code changes to get that test to pass.

## Update the code

At this point, make your code changes and constantly test again your new test case to make sure that everything working
properly. Do [commit] early and often and make useful commit messages.

If a backwards incompatible change cannot be avoided, please make sure to call that out when you submit a pull request,
explaining why the change is absolutely necessary.

## Create a pull request

[Create a pull request] with your changes. For it, this repository includes a [pull request template](.github/PULL_REQUEST_TEMPLATE.md) that you can use to help you write a good description of your changes.

## Merge and release

The [CODEOWNERS](.github/CODEOWNERS) of the repository will review your code and provide feedback.
If everything looks good, they will merge the code and release a new version while following the principles of [Semantic Versioning (SemVer)].

<!-- References -->

<!-- markdown-link-check-disable -->

[Pull Requests]: https://github.com/mineiros-io/terraform-aws-dynamodb/pulls
[pre-commit-file]: https://github.com/mineiros-io/terraform-aws-dynamodb/blob/master/.pre-commit-config.yaml

<!-- markdown-link-check-enable -->
