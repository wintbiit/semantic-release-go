# Semantic Release Go
> [semantic-release](https://semantic-release.gitbook.io/) like versioning, commit analyzer, and changelog generator

## Features
- [x] Versioning
- [x] Commit Analyzer
- [x] Changelog Generator
- [x] Branch hassle free
- [x] CI/CD Friendly

## Usage
### Command:
```bash
semantic-release -path <path> -branch <branch>
```

### Arguments:
- `-path`: Path to the git repository
- `-branch`: Branch to release (not necessarily the git branch, you can name it anything)
- `-channel`: Channel to release (values: `insider`(default), `alpha`, `beta`, `release`)
- `-dry`: Dry run (values: `true`, `false`(default))
- `-debug`: Debug mode (values: `true`, `false`(default))
- `-analyzer`: Commit analyzer (values: `angular`(default))
- `-repo`: Repository URL
- `-tag`: If create tag (values: `true`(default), `false`)
- `-push`: If push to remote (values: `true`(default), `false`)
- `-changelog`: Changelog file path(default: `Changelog.md`, empty for skip)

## Inspired by
- [semantic-release](https://semantic-release.gitbook.io/)

## Note
Before v1.4.5, this project uses `go-git` as git parser, but it's incomplete and buggy. So, I switched to `os/exec` for git commands.