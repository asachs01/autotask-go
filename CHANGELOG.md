# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complete webhook implementation with proper event handling
- Webhook signature verification for security
- Webhook handler registration system
- Example code and documentation for webhooks
- Comprehensive test suite for client functionality
- Comprehensive test suite for webhook functionality
- Comprehensive test suite for query builder functionality
- Comprehensive test suite for pagination functionality including:
  - Pagination iterator creation and usage
  - Page fetching with proper parameter handling
  - Multi-page result fetching
  - Callback-based pagination processing

### Fixed
- Fixed pagination parameter handling where parameters were incorrectly included in filter string
- Updated FetchPage function to set page parameter in query parameters

## [1.1.0] - 2025-03-06

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
- Enhanced pagination support with:
  - Iterator pattern for efficient traversal of paginated results
  - Convenience methods for common pagination scenarios
  - Generic type support for strongly-typed results
  - Callback-based processing for large result sets
  - Proper handling of Autotask API pagination requirements
  - Improved pagination using Autotask's nextPageUrl/prevPageUrl mechanism

### Changed
- Reorganized project structure according to Go standards
- Improved error handling and reporting
- Enhanced documentation and examples

### Fixed
- Fixed context propagation in API calls
- Fixed telemetry span handling
- Fixed rate limit handling
- Fixed linter errors related to unchecked error returns
- Fixed page numbering in pagination results to ensure consistent and accurate page numbers

### Removed
- Removed zone metrics (to be reimplemented with better design)

## [1.0.0] - 2024-03-05

### Added
- Initial release
- Basic Autotask API client implementation
- Support for Companies, Tickets, and Contacts entities
- Basic error handling and logging 