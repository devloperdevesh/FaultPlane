# Versioning Policy

FaultPlane follows **Semantic Versioning (SemVer)** to provide predictable release management and compatibility expectations.

Version format:

```
MAJOR.MINOR.PATCH
```

Example:

```
v1.4.2
```

Each version component represents a different type of change.

---

# Version Components

## MAJOR Version

A MAJOR version indicates breaking changes that may require user action.

Format:

```
v1.x.x → v2.x.x
```

Examples:

- breaking API changes
- incompatible configuration changes
- major architecture changes
- removal of previously supported features

Migration guidance will be provided for major releases whenever possible.

---

## MINOR Version

A MINOR version introduces new functionality while maintaining backward compatibility.

Format:

```
v1.1.x → v1.2.x
```

Examples:

- new runtime capabilities
- additional configuration options
- new integrations
- improved developer experience

Existing users should be able to upgrade without major changes.

---

## PATCH Version

A PATCH version contains backward-compatible fixes and improvements.

Format:

```
v1.2.1 → v1.2.2
```

Examples:

- bug fixes
- security patches
- documentation improvements
- performance improvements
- internal fixes

---

# Pre-Release Versions

Before reaching a stable release, FaultPlane may publish pre-release versions.

Examples:

```
v0.1.0-alpha
v0.1.0-beta
v1.0.0-rc.1
```

Pre-release versions may include:

- experimental features
- unstable APIs
- incomplete functionality

Production usage may require additional validation.

---

# Development Versions

During active development:

```
v0.x.x
```

indicates that the project is still evolving.

During this phase:

- APIs may change
- internal architecture may evolve
- migration paths may be introduced gradually

---

# Release Process

Each release follows a structured workflow.

```
Development

      ↓

Testing & Validation

      ↓

Update CHANGELOG.md

      ↓

Create Git Tag

      ↓

Publish GitHub Release

      ↓

Release Announcement
```

---

# Release Checklist

Before publishing a release:

## Engineering

- [ ] Tests passing
- [ ] Linting completed
- [ ] Documentation updated
- [ ] Breaking changes reviewed

## Release Management

- [ ] CHANGELOG.md updated
- [ ] Version tag created
- [ ] GitHub Release published
- [ ] Release notes prepared

## Communication

- [ ] Contributors notified
- [ ] Migration notes added if required

---

# Compatibility Policy

FaultPlane aims to maintain predictable compatibility.

Stable releases should provide:

- documented APIs
- upgrade guidance
- migration notes for breaking changes

Experimental features may evolve faster during pre-release development.

---

# Release Philosophy

FaultPlane follows incremental releases focused on:

- reliability improvements
- operational stability
- developer experience
- maintainable architecture

Version numbers communicate the maturity and compatibility expectations of each release.
