---
name: 🪲 Bug Report
about: Report issues with daggerverse modules or tooling
title: "[Module/Tool Name]: Brief description of bug"
labels: "bug"
---

## 🪲 What's Not Working

Provide a clear explanation of what's not working as expected.
Example: "The `aws` module's `deploy` function fails when using a custom VPC configuration"

### Expected Behavior

Provide a clear explanation of what you expected to happen.
Example: "The `aws` module's `deploy` function should create an EC2 instance with the specified VPC configuration"

### Actual Behavior

Provide a clear explanation of what actually happened.
Example: "The `aws` module's `deploy` function fails when using a custom VPC configuration"

## 🎯 Component Type

Please select what you're having issues with:

- [ ] Daggerverse Module (specify which: **\_\_\_**)
- [ ] Daggy CLI Tool
- [ ] Module Template Full
- [ ] Module Template Light
- [ ] Documentation
- [ ] Tooling (E.g.: Justfile)
- [ ] Other (Please describe)

## 📋 Environment Details

- Dagger Version: (e.g., v0.9.3)
- Module/Tool Version: (e.g., aws@v0.2.1)
- OS: (e.g., Ubuntu 22.04)

## 🔬 Reproduction

### Code Sample

```go
// Include your Dagger Function call or relevant code
```

### Dagger Version

```sh
dagger version
```

### Module/Tool Version

```sh
module/tool version
```
