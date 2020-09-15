# Matrix
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fchaos-mesh%2Fmatrix.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fchaos-mesh%2Fmatrix?ref=badge_shield)


Matrix is a config fuzzer, i.e. generate a set of configurations for a program that need to be fuzzed.

## Quickstart
```bash
$ make
$ bin/matrix -c examples/matrix-tidb-tikv-pd.yaml
# this will generate three TiDB-related configs in current folder
# including `tidb.yaml`, `tikv.yaml`, `pd.yaml`, `mysql-system-vars.sql` and `tidb-system-vars.sql`
```

Output folder can be configured with `-d <output folder>`

## Usage
```
$ ./bin/matrix -h          
Usage of ./bin/matrix:
  -c string
        config file
  -d string
        output folder (default ".")
  -s int
        seed of rand (default UTC nanoseconds of now)
```

Matrix reads one config file contains generation rules,
generate values based on those rules,
then dump values with serializer.

## Matrix Configuration
Matrix use DSL in yaml format to describe the configurations need to be fuzzed.
See [DSL](./docs/DSL.md) for detailed information.

## Type of configs Matrix support to generate
With different serializer, Matrix supports to generate various kinds of configurations,
refer to [serializer](./docs/serializer.md) for detailed usage.

Currently supported format:
- toml
- yaml
- line-based generation

## TODO
### Parser
- [ ] Better unit conversion
### Serializer
- [x] toml serializer
- [x] yaml serializer
- [x] SQL serializer
- [ ] Command line argument serializer
### Generator
- [ ] List generation
- [ ] Instructive information (e.g. default value)
- [x] Simple recursive random generator
- [ ] Non-recursive generator that supports dependency
### Random
- [x] Random generate value
- [x] Specific random seed to have same output with same seed, to support continuous testing


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fchaos-mesh%2Fmatrix.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fchaos-mesh%2Fmatrix?ref=badge_large)