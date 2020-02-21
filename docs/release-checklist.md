1) Verify that all tests are passing.
1) Update `CHANGELOG.md` and commit.
1) Create release branch:
    ```
    git checkout -b release/v0.2
    ```
1) Update `version.go`.
    * `0.2.0-dev` -> `0.2.0`
1) Commit and ag release.
    ```
    git add -p
    git commit -m "Release 0.2.0"
    git tag -a v0.2.0 -m "Release 0.2.0"
    git push --set-upstream origin release/v0.2
    ```
1) Build and upload release binaries.
    ```
    goreleaser
    ```
1) Update `version.go` and commit.
    * `0.2.0` -> `0.2.1-dev`
1) Push commits and tags.
    ```
    git add -p
    git commit -m "bump patch version"
    git push origin release/v0.2
    ```
1) Review release on GitHub.
    * https://github.com/sjansen/carpenter/releases
