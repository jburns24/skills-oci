# How to Create a New Release

This guide walks through the steps to create a new release of the `skills-oci` CLI. The release process is automated via GitHub Actions — you only need to push a version tag.

## Prerequisites

- Push access to the `salaboy/skills-oci` repository
- All changes for the release are merged into `main`
- CI is passing on `main`

## Steps

### 1. Decide on a version number

Follow [Semantic Versioning](https://semver.org/):

- **Patch** (`v1.0.1`) — bug fixes, no new features
- **Minor** (`v1.1.0`) — new features, backwards compatible
- **Major** (`v2.0.0`) — breaking changes

### 2. Make sure main is up to date

```bash
git checkout main
git pull origin main
```

### 3. Verify the build and tests pass locally

```bash
go build ./...
go test ./...
```

### 4. Create and push a version tag

```bash
git tag v1.0.0
git push origin v1.0.0
```

Replace `v1.0.0` with your chosen version number. The tag **must** start with `v` to trigger the release workflow.

### 5. Monitor the release workflow

Open the **Actions** tab in GitHub to watch the workflow run:

```
https://github.com/salaboy/skills-oci/actions
```

The workflow will:

1. Check out the code at the tagged commit
2. Run all tests
3. Cross-compile binaries for five targets:
   - `skills-oci-linux-amd64`
   - `skills-oci-linux-arm64`
   - `skills-oci-darwin-amd64`
   - `skills-oci-darwin-arm64`
   - `skills-oci-windows-amd64.exe`
4. Generate a `checksums.txt` file with SHA256 hashes
5. Create a GitHub release with auto-generated release notes and all binaries attached

### 6. Verify the release

Once the workflow completes, check the release page:

```
https://github.com/salaboy/skills-oci/releases
```

Confirm that:

- All five binaries and the checksums file are attached
- The release notes accurately describe the changes since the last release
- You can download and run one of the binaries:

```bash
# Example: download and test the macOS arm64 binary
curl -L -o skills-oci https://github.com/salaboy/skills-oci/releases/download/v1.0.0/skills-oci-darwin-arm64
chmod +x skills-oci
./skills-oci --help
```

### 7. (Optional) Edit the release notes

If the auto-generated notes need adjustments, edit them directly on the GitHub release page. You can add a summary section at the top highlighting the key changes.

## Troubleshooting

### The workflow didn't trigger

- Make sure the tag starts with `v` (e.g., `v1.0.0`, not `1.0.0`)
- Check that the tag was pushed to the remote: `git ls-remote --tags origin`

### Tests fail during the release

The release workflow runs tests before building. If tests fail:

1. Delete the tag locally and remotely:

   ```bash
   git tag -d v1.0.0
   git push origin --delete v1.0.0
   ```

2. Fix the failing tests on `main`
3. Re-tag and push once CI is green

### Releasing a hotfix

If you need to patch an older release:

1. Create a branch from the release tag: `git checkout -b release/v1.0.x v1.0.0`
2. Cherry-pick or apply the fix
3. Tag the branch: `git tag v1.0.1`
4. Push the tag: `git push origin v1.0.1`
