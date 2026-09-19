# Contributing to timeutil

Thank you for considering contributing to **timeutil**! This document outlines the guidelines and workflow for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Pull Request Process](#pull-request-process)
- [Coding Guidelines](#coding-guidelines)
- [Testing](#testing)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

Please be respectful and considerate in all interactions. We aim to maintain a welcoming and inclusive environment for everyone.

## Getting Started

1. **Fork** the repository on GitHub.
2. **Clone** your fork locally:

   ```bash
   git clone https://github.com/<your-username>/timeutil.git
   cd timeutil
   ```

3. **Add the upstream remote**:

   ```bash
   git remote add upstream https://github.com/ghettovoice/timeutil.git
   ```

## Development Setup

### Prerequisites

- Go 1.25 or later
- [Task](https://taskfile.dev/) (used for task automation)

### Setup

```bash
task test
```

This will run the test suite and confirm that your environment is ready.

### Build

The project is a library, so there is no standalone build step. Make sure the package compiles:

```bash
go build ./...
```

## Making Changes

1. **Create a new branch** from `main`:

   ```bash
   git checkout -b feature/my-feature
   # or
   git checkout -b fix/my-bugfix
   ```

2. **Make your changes** following the [coding guidelines](#coding-guidelines).

3. **Run tests and linters**:

   ```bash
   task test
   task lint
   task vuln

   # or run all at once
   task check
   ```

4. **Commit your changes** with a clear and descriptive commit message:

   ```bash
   git commit -m "feat: add feature X"
   # or
   git commit -m "fix: issue with Y"
   ```

5. **Push to your fork**:

   ```bash
   git push origin feature/my-feature
   ```

6. **Open a Pull Request** against the `main` branch.

## Pull Request Process

1. Ensure all tests pass and there are no linter warnings.
2. Update documentation if your changes affect the public API.
3. Fill out the pull request template completely.
4. Link any related issues using `Fixes #issue-number`.
5. Wait for review and address any feedback.

### PR Types

- **Bug fix** — non-breaking change that fixes an issue
- **New feature** — non-breaking change that adds functionality
- **Breaking change** — fix or feature that would cause existing functionality to not work as expected
- **Documentation update** — changes to docs only
- **Refactoring** — no functional changes

## Coding Guidelines

- Follow idiomatic Go conventions and the [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).
- Use `gofumpt` for formatting (enforced by the linter).
- Keep functions small and focused on a single responsibility.
- Write descriptive variable and function names.
- Add comments to explain non-obvious logic and to document all exported symbols.
- Avoid unnecessary dependencies.

## Testing

### Run all tests

```bash
task test
```

### View coverage report

```bash
task check
task cov
```

### Guidelines

- Write tests for all new functionality.
- Ensure existing tests pass before submitting a PR.
- Include both positive and negative test cases.
- Use table-driven tests where appropriate.

## Reporting Issues

When reporting issues, please include:

- A clear and descriptive title.
- Steps to reproduce the problem.
- Expected vs. actual behavior.
- Go version and OS information.
- Relevant code snippets or error messages.

Use the appropriate issue template when creating a new issue.

## License

By contributing to this project, you agree that your contributions will be licensed under the [MIT License](./LICENSE).
