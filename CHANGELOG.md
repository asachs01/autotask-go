# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Context support for all API operations
- OpenTelemetry integration for metrics and tracing
- Rate limiting with exponential backoff
- Structured logging
- Retry mechanism for transient errors
- Added support for additional entity services:
  - Projects
  - Tasks
  - Time Entries
  - Contracts
  - Configuration Items
- Enhanced query builder with support for:
  - Complex filtering with logical operators (AND, OR)
  - Nested conditions with parentheses
  - Automatic type conversion for filter values
  - Natural language filter syntax

### Changed
- Reorganized project structure according to Go standards
- Improved error handling and reporting
- Enhanced documentation and examples

### Fixed
- Fixed context propagation in API calls
- Fixed telemetry span handling
- Fixed rate limit handling

### Removed
- Removed zone metrics (to be reimplemented with better design)

## [0.1.0] - 2024-03-05

### Added
- Initial release
- Basic Autotask API client implementation
- Support for Companies, Tickets, and Contacts entities
- Basic error handling and logging 