# Bitrise Init

This repository hosts the `bitrise-init` which contains all the shared project detection and config generation logic. 

This package ia consumed by the following tools:
- [project-scanner step](https://github.com/bitrise-steplib/steps-project-scanner)
- [bitrise-init plugin](https://github.com/bitrise-io/bitrise-plugins-init)
- [bitrise-add-new-project](https://github.com/bitrise-io/bitrise-add-new-project)

## How to release new bitrise-init version

- update the step versions in steps/const.go
    - `go get -u github.com/godrei/stepper`
    - `stepper stepLatests --steps-const-file="steps/const.go"`
    - copy the output after “Generated” to the const.go file
- bump `version` in version/version.go
- commit these changes & open PR
- merge to master
- create tag with the new version
- test the generated release and its binaries

### Update manual config on website

- Use the included go app to generate the manual configuration:

```
~/path/to/bitrise-init ❯❯❯ cd _manual-config
~/p/t/b/_manual-config ❯❯❯ go run main.go
Generating manual config
Config saved to generated/result.yml
~/p/t/b/_manual-config ❯❯❯ 
```

This will generate the manual configuration yaml file to `_manual-config/generated/result.yml`.

- Update the file https://github.com/bitrise-io/bitrise-website/blob/master/config/bitrise_ymls/custom_config.yml, with the contents of `results.yml`.

### Update the [project-scanner step](https://github.com/bitrise-steplib/steps-project-scanner)

- Update the bitrise-init dependency
- Share a new version into steplib

### Update the [bitrise init plugin](https://github.com/bitrise-io/bitrise-plugins-init)

- Update the bitrise-init dependency
- Release a new version.

### Update the [bitrise-add-new-project](https://github.com/bitrise-io/bitrise-add-new-project)

- Update the bitrise-init dependency
- Release a new version.