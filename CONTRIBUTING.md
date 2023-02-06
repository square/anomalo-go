# Contributing

We welcome contributions. For small changes, open a pull request. If you'd like to discuss an improvement or fix before writing the code, please open a Github issue.

Before code can be accepted all contributors must complete our [Individual Contributor License Agreement (CLA)](https://spreadsheets.google.com/spreadsheet/viewform?formkey=dDViT2xzUHAwRkI3X3k5Z0lQM091OGc6MQ&ndplr=1)

## Reporting Issues

Please report issues in Github issues.

# Table of Contents

- [Contributing](#contributing)
  - [Reporting Issues](#reporting-issues)
- [Table of Contents](#table-of-contents)
  - [Technical Details](#technical-details)
    - [Local Development](#local-development)


## Technical Details
### Local Development
1. Clone the repo
2. Run `go get && go mod tidy`

### Releasing a new version
Tag the branch with the appropriate version number (ex v1.2.0). 

Run `git push origin v1.2.0`

Docs: https://go.dev/blog/publishing-go-modules