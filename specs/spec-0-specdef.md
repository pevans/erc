---
Specification: 0
Drafted At: 2026-03-03
Authors:
  - Peter Evans
---

# 1. Overview

A specification (or "spec") is a structured markdown document that describes
some functionality that can be implemented in the system. Unlike pure vibe
coding, a spec is a long-lived written artifact that can be referred to after
the the initial implementation has been written.

# 2. Sections

A spec is broken down into sections separated by markdown headings. Each
section is numbered.

## 2.1. Subsections

Each section can be further broken down into subsections. A subsection is
related to the topic of the parent section, but can be used to describe a
distinct detail or component of the overall topic covered with the section.

# 3. Implementation

Specs (besides this one!) are generally meant to be implementable with code,
but should not be written with a specific language in mind. It should be
possible for someone to implement the concepts in a spec with any language.

## 3.1. Pseudocode

Should the need arise to describe a specific algorithm to be used for some
purpose, pseudocode should be used. Any manner of pseudocode is fine --
C-like; Python-like; etc. -- as long as a reader can correctly interpret the
intention of the algorithm.
