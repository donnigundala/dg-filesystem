# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-12-27

### Added
- Initial stable release of the `dg-filesystem` plugin.
- **Multi-Driver Support**: Local storage and AWS S3 with unified API.
- **File Operations**: Put, Get, Delete, Exists, Size, MimeType, LastModified.
- **Directory Operations**: MakeDirectory, DeleteDirectory, Files, AllFiles, Directories.
- **Advanced Features**:
  - Temporary URLs for S3 (pre-signed URLs)
  - Visibility control (public/private)
  - Streaming support for large files
  - Path manipulation utilities
- **Container Integration**: Auto-registration with Injectable pattern and helper functions.
- **Observability**: OpenTelemetry metrics for file operations and performance tracking.

### Features
- Unified Disk interface for all storage backends
- Fluent configuration API
- Thread-safe operations
- Automatic MIME type detection
- Production-ready with comprehensive test coverage

### Documentation
- Complete README with examples
- Driver-specific configuration guides
- S3 integration documentation

### Performance
- Efficient streaming for large files
- Optimized S3 operations with AWS SDK v2
- Minimal memory footprint

---

## Development History

The following versions represent the development journey leading to v1.0.0:

### 2025-11-24
- OpenTelemetry observability integration
- Metrics for file operations

### 2025-11-23
- AWS S3 driver implementation
- Temporary URL support
- Visibility control

### 2025-11-22
- Initial implementation with local storage driver
- Core file operations
- Directory management
