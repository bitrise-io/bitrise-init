---
name: database-version-maintenance
description: Instructions for upgrading the pinned Docker image major versions in the Ruby scanner's database configuration to the latest available versions.
disable-model-invocation: true
allowed-tools: WebFetch, WebSearch
---

### Context

The Ruby scanner detects database gems (pg, mysql2, redis, mongoid, mongo) in a project's Gemfile and generates Bitrise workflow configuration that spins up the corresponding service containers. The Docker image versions for those containers are hardcoded in `scanners/ruby/databases.go` in the `knownDatabaseGems` slice.

Currently pinned images:
- `postgres:18`
- `mysql:9`
- `redis:8`
- `mongo:8`

These version numbers can become outdated over time. This skill guides you through checking for newer major versions and updating the file accordingly.

### Instructions

Prerequisites:

- Internet access to reach Docker Hub (via WebFetch or WebSearch)
- A proper Go environment according to @.tool-versions

If the above are not met, do not proceed, just flag the issue to the user.

1. Read `scanners/ruby/databases.go` and collect all Docker image references from the `knownDatabaseGems` slice (the `image` field).

2. For each image, determine the latest stable major version using WebSearch (e.g. "latest stable postgres docker image major version"). Look for the highest integer-only major version tag (e.g. `18`, `9`, `8`) that is a stable/general-availability release. Ignore tags like `latest`, `alpine`, `beta`, `rc`, or any tag containing a `-`.

   Alternatively, you can try querying the Docker Hub v2 API and filtering the `name` fields in the JSON response for integer-only tags. Note: sorting by `last_updated` surfaces patch tags (e.g. `17.5`) rather than major-version tags, so filter carefully or use a larger `page_size`.
   - `https://hub.docker.com/v2/repositories/library/postgres/tags?page_size=100`
   - `https://hub.docker.com/v2/repositories/library/mysql/tags?page_size=100`
   - `https://hub.docker.com/v2/repositories/library/redis/tags?page_size=100`
   - `https://hub.docker.com/v2/repositories/library/mongo/tags?page_size=100`

3. Update the `image` field values in `scanners/ruby/databases.go` if a newer major version is available.

4. If the major version of an image changed, check the official Docker Hub page or release notes for that image for breaking changes (e.g. renamed environment variables, changed default behaviour, dropped support for older authentication methods). If breaking changes affect the `healthCheck` commands or `containerEnvKey` values in `knownDatabaseGems`, update those as well.

5. Run Go unit tests to verify your changes:
   ```
   go test ./scanners/ruby/...
   ```
   There might be failing tests unrelated to the changes (mostly tooling issues). In this case, go ahead and let CI be the judge.

6. Create a new branch, commit your changes, and open a PR.
