# How to Publish the CLI as a Homebrew Package

This guide explains how to distribute the `skills-oci` CLI via [Homebrew](https://brew.sh/) so users can install it with `brew install`.

## Overview

There are two approaches:

1. **Homebrew Tap** (recommended) — Your own repository that hosts the formula. Users install with `brew install salaboy/tap/skills-oci`.
2. **Homebrew Core** — The official Homebrew repository. Requires meeting their [acceptance criteria](https://docs.brew.sh/Acceptable-Formulae) (notable project, sufficient GitHub stars, etc.).

This guide covers the Tap approach since it's available immediately for any project.

## Prerequisites

- A published GitHub release with compiled binaries and checksums (see [How to Create a New Release](./create-a-new-release.md))
- A GitHub account with permission to create repositories

## Steps

### 1. Create a Homebrew Tap repository

Create a new GitHub repository named `homebrew-tap` under your account:

```
https://github.com/salaboy/homebrew-tap
```

The `homebrew-` prefix is required — it allows Homebrew to resolve `salaboy/tap` to `salaboy/homebrew-tap`.

### 2. Get the SHA256 checksums for the release

Download the checksums from your latest release, or compute them from the binaries:

```bash
# From the checksums.txt in the release
curl -sL https://github.com/salaboy/skills-oci/releases/download/v0.0.1/checksums.txt
```

Note the SHA256 values for `skills-oci-darwin-amd64` and `skills-oci-darwin-arm64` (and Linux if you want Linux support).

### 3. Create the formula

In your `homebrew-tap` repository, create `Formula/skills-oci.rb`:

```ruby
class SkillsOci < Formula
  desc "CLI tool for packaging and managing AI agent skills as OCI artifacts"
  homepage "https://github.com/salaboy/skills-oci"
  version "1.0.0"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/salaboy/skills-oci/releases/download/v#{version}/skills-oci-darwin-arm64"
      sha256 "REPLACE_WITH_DARWIN_ARM64_SHA256"
    else
      url "https://github.com/salaboy/skills-oci/releases/download/v#{version}/skills-oci-darwin-amd64"
      sha256 "REPLACE_WITH_DARWIN_AMD64_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/salaboy/skills-oci/releases/download/v#{version}/skills-oci-linux-arm64"
      sha256 "REPLACE_WITH_LINUX_ARM64_SHA256"
    else
      url "https://github.com/salaboy/skills-oci/releases/download/v#{version}/skills-oci-linux-amd64"
      sha256 "REPLACE_WITH_LINUX_AMD64_SHA256"
    end
  end

  def install
    binary_name = stable.url.split("/").last
    bin.install binary_name => "skills-oci"
  end

  test do
    assert_match "Manage agent skills as OCI artifacts", shell_output("#{bin}/skills-oci --help")
  end
end
```

Replace each `REPLACE_WITH_..._SHA256` with the actual checksum from step 2.

### 4. Commit and push the formula

```bash
cd homebrew-tap
git add Formula/skills-oci.rb
git commit -m "Add skills-oci formula v1.0.0"
git push origin main
```

### 5. Test the installation

```bash
# Add the tap
brew tap salaboy/tap

# Install
brew install salaboy/tap/skills-oci

# Verify
skills-oci --help
```

### 6. Users can now install with

```bash
brew tap salaboy/tap
brew install skills-oci
```

Or in a single command:

```bash
brew install salaboy/tap/skills-oci
```

## Updating the Formula for New Releases

When you publish a new release, update the formula:

### 1. Download the new checksums

```bash
curl -sL https://github.com/salaboy/skills-oci/releases/download/v1.1.0/checksums.txt
```

### 2. Update the formula

In `homebrew-tap/Formula/skills-oci.rb`:

- Update the `version` field
- Replace all SHA256 checksums with the new values

### 3. Commit and push

```bash
git commit -am "Update skills-oci formula to v1.1.0"
git push origin main
```

Users will get the new version on their next `brew upgrade`.

## Automating Formula Updates

You can automate formula updates by adding a step to your release workflow. Add this job to `.github/workflows/release.yml` after the release job:

```yaml
  update-homebrew:
    needs: release
    runs-on: ubuntu-latest
    steps:
      - name: Update Homebrew formula
        uses: mislav/bump-homebrew-formula-action@v3
        with:
          formula-name: skills-oci
          homebrew-tap: salaboy/homebrew-tap
          download-url: https://github.com/salaboy/skills-oci/releases/download/${{ github.ref_name }}/skills-oci-darwin-arm64
        env:
          COMMITTER_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```

This requires a personal access token (`HOMEBREW_TAP_TOKEN`) with write access to the `homebrew-tap` repository, stored as a repository secret.

> **Note:** The `mislav/bump-homebrew-formula-action` works best with source tarballs or single-URL formulas. For multi-architecture formulas like ours, you may need a custom script that updates all SHA256 values. An alternative is to use [GoReleaser](https://goreleaser.com/) which has built-in Homebrew Tap support and handles multi-arch formulas automatically.

## Alternative: Using GoReleaser

[GoReleaser](https://goreleaser.com/) can replace both the release workflow and Homebrew formula management. It cross-compiles, creates GitHub releases, and updates Homebrew taps in a single tool. If you adopt GoReleaser, create a `.goreleaser.yml`:

```yaml
builds:
  - binary: skills-oci
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w -X main.version={{.Version}}

brews:
  - repository:
      owner: salaboy
      name: homebrew-tap
    homepage: https://github.com/salaboy/skills-oci
    description: CLI tool for packaging and managing AI agent skills as OCI artifacts
    license: Apache-2.0
```

Then replace the release workflow with:

```yaml
- name: Run GoReleaser
  uses: goreleaser/goreleaser-action@v6
  with:
    version: latest
    args: release --clean
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    HOMEBREW_TAP_TOKEN: ${{ secrets.HOMEBREW_TAP_TOKEN }}
```
