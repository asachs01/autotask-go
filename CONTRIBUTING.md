# Contributing to Autotask Go

Thank you for your interest in contributing to the Autotask Go client library! This document provides guidelines and instructions for contributing.

## Development Setup

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/autotask-go.git
   cd autotask-go
   ```
3. Create a new branch for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Run `go vet` to check for potential issues
- Follow the existing code style in the project

## Testing

- Write unit tests for new functionality
- Ensure all tests pass:
  ```bash
  go test ./...
  ```
- Add integration tests for API interactions
- Add benchmarks for performance-critical code

## Documentation

- Add godoc comments for all exported types and functions
- Update README.md if adding new features
- Add examples in the `examples` directory
- Document any breaking changes

## Pull Request Process

1. Update your branch with the latest changes from main:
   ```bash
   git fetch origin
   git rebase origin/main
   ```
2. Push your changes:
   ```bash
   git push origin feature/your-feature-name
   ```
3. Create a Pull Request
4. Provide a clear description of your changes
5. Reference any related issues

## Code Review

- Respond to review comments promptly
- Make requested changes
- Keep commits focused and atomic
- Squash commits when appropriate

## Release Process

1. Update version in go.mod
2. Update CHANGELOG.md
3. Create a new release tag
4. Update documentation

## Questions?

Feel free to open an issue for any questions or concerns about contributing. 