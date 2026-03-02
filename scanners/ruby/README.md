# Ruby Scanner

This scanner detects Ruby projects and generates Bitrise CI/CD configurations for them.

## Detection Criteria

The scanner looks for projects containing a `Gemfile` in the project directory.

## Detection Logic

For each Ruby project found, the scanner detects:

1. **Bundler**: Checks for `Gemfile.lock` to determine if Bundler is used
2. **Rakefile**: Checks for `Rakefile` to enable rake-based task execution
3. **Test Framework**:
   - RSpec: Detected by presence of `spec/spec_helper.rb` or `.rspec`
   - Minitest: Detected by presence of `test/test_helper.rb`
4. **Ruby Version**: Checks for `.ruby-version` or `.tool-versions` files

## Generated Workflow

The generated workflow includes:

1. **Setup Steps**:
   - Activate SSH key (if configured)
   - Git clone
   - Install Ruby (using asdf if version file is present)

2. **Dependency Management**:
   - Restore gem cache
   - Install dependencies with Bundler (`bundle install`)

3. **Test Execution**:
   - Run tests based on detected framework:
     - RSpec: `bundle exec rspec`
     - Minitest: `bundle exec rake test` (if Rakefile exists)
     - Default: `bundle exec rake test` (if Rakefile exists)

4. **Caching & Deployment**:
   - Save gem cache
   - Deploy to Bitrise

## Configuration Options

### Project Directory
- **Title**: "Project Directory"
- **Summary**: "The directory containing the Gemfile"
- **Environment Key**: `RUBY_PROJECT_DIR`
- **Type**: Selector (for detected projects) or User Input (for manual config)

## Example Gemfile Detection

```
project-root/
├── Gemfile
├── Gemfile.lock
├── Rakefile
├── .ruby-version
└── spec/
    └── spec_helper.rb
```

This structure would be detected as a Ruby project with:
- Bundler enabled
- Rakefile present
- RSpec as test framework
- Ruby version file present

## Config Naming Convention

Generated config names follow the pattern:
- `ruby-[root]-[bundler]-[testframework]-config`
- `default-ruby-config` (for manual configuration)

Examples:
- `ruby-root-bundler-rspec-config`
- `ruby-bundler-minitest-config`
- `ruby-config` (minimal, no bundler, no test framework)
