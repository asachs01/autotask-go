# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]


## [1.2.1] - 2025-04-14

### Fixed
- Fixed logger method calls in queryWithEmptyFilter to use correct method signature instead of chaining syntax
- Improved logging implementation to properly handle structured logging fields 

### Changed
- Removed unused duplicate queryWithEmptyFilter function in favor of QueryWithEmptyFilter in query_helpers.go

## [1.2.0] - 2025-03-19

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
- Comprehensive test suite for entity operations including:
  - Entity retrieval (Get)
  - Entity querying with filters
  - Entity creation
  - Entity updates
  - Entity deletion
  - Entity counting

### Fixed
- Fixed pagination parameter handling where parameters were incorrectly included in filter string
- Updated FetchPage function to set page parameter in query parameters
- Fixed client.Do method to properly handle 204 No Content responses
- Fixed entity tests to match the actual API implementation
- Fixed unchecked error returns in response body closing operations
- Fixed unchecked error returns in test write operations
- Fixed unchecked error returns in file operations and logging
- Fixed error handling in HTTP response body management
- Improved error handling in test helpers and mock servers

### Changed
- Optimized code structure by converting if-else chains to switch statements
- Removed redundant type declarations for better code clarity
- Cleaned up unused code and functions
- Enhanced test reliability with proper error checking

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

## [1.0.0] - 2025-03-05

### Added
- Initial release
- Basic Autotask API client implementation
- Support for Companies, Tickets, and Contacts entities
- Basic error handling and logging
