## github-commit-status

This is a simple utility to update the status of a commit on github. The
primary use case is to update the status of a commit in a build environment.

## Install

Download the latest binary or
`brew tap thbishop/github-commit-status && brew install github-commit-status`
if you're on OSX.

## Usage

Create a [github token]() with TODO scope and export it as an env var:
```sh
export GITHUB_TOKEN=1234
```

```sh
github-commit-status --user foo --repo bar --commit $SHA --status success
```

If you're using github enterprise, you can set the API endpoint like so:
```sh
export GITHUB_API=https://github.example.com/api/v3
```

## Contribute
* Fork the project
* Make your feature addition or bug fix (with tests and docs) in a topic branch
* Make sure tests pass
* Send a pull request and I'll get it integrated

## LICENSE
See [LICENSE](LICENSE)
