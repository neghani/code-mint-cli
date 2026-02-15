# Release steps (new build with fixes)

## Done for v0.1.4

- [x] All changes staged and committed (installer paths, remove empty dir, FEATURES_ROADMAP, coordination docs).
- [x] Annotated tag **v0.1.4** created locally.

## You do: push to trigger the release

From your machine (with GitHub auth set up):

```bash
cd /Users/ganeshkumar/Projects/code-mint-ali-cli

# Push main (includes the release commit)
git push origin main

# Push the tag — this triggers the GitHub Action "Release" and creates the release + artifacts
git push origin v0.1.4
```

If you use SSH: `git push origin main && git push origin v0.1.4`

## 3. What the workflow does

On **push of tag `v*`**:

1. Checkout repo, setup Go from `go.mod`
2. Derive `VERSION` from tag (e.g. `v0.1.4` → `0.1.4`)
3. Run `make release-dry-run VERSION="0.1.4"` to build:
   - `codemint_0.1.4_darwin_arm64.tar.gz`
   - `codemint_0.1.4_darwin_amd64.tar.gz`
   - `codemint_0.1.4_linux_amd64.tar.gz`
   - `codemint_0.1.4_linux_arm64.tar.gz`
   - `codemint_0.1.4_windows_amd64.zip`
   - `SHA256SUMS`
   - `install.sh` (copy of `scripts/install.sh`)
4. Create GitHub Release with tag and upload all artifacts; release notes are auto-generated.

## 4. Verify after release

- **Releases:** https://github.com/neghani/code-mint-cli/releases
- **One-liner (root URL):**  
  `curl -fsSL https://raw.githubusercontent.com/neghani/code-mint-cli/main/install.sh | sh`
- **One-liner (releases):**  
  `curl -fsSL https://github.com/neghani/code-mint-cli/releases/latest/download/install.sh | sh`

## Version summary (v0.1.4)

- **Windsurf / Continue / Claude:** skills install to `.windsurf/skills`, `.continue/skills`, `.claude/skills` (was rules dir).
- **Remove:** after removing a file, empty parent dir is removed (e.g. `.cursor/skills/<slug>/`).
- **Docs:** FEATURES_ROADMAP.md, execution plan coordination with code-mint-ai, REPO_COORDINATION in app repo.
