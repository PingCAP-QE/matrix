# Matrix

Matrix is a config fuzzer, i.e. generate a set of configurations for a program that need to be fuzzed.

## Quickstart
```bash
make
bin/matrix -c examples/matrix-tidb.yaml # this will generate a TiDB config `tidb.yaml` in current folder
```

Output folder can be configured with `-d <output folder>`

## Type of configs
- Domain specific config file in specific format;
   - [x] toml
   - [ ] yaml
- [ ] Command line arguments;
- [ ] Command (Query in Database) need to be run.

Note 3. could also be used for dynamic configuring (change config during runtime),
which will not be considered now.

## DSL

Matrix use DSL in yaml format to describe the configurations need to be fuzzed.

### Configs of Matrix
The config file of Matrix tool is a yaml file, formatted as followed:
```yaml
config:
  - <config_template>
```

Each `config_template` is a set of configs that will be generated as one entity (i.e. file).

| Field | Description | Sample Value |
|:------|:------------------|:--------------|
| **tag** | Tag for this config |
| **serializer** | Serializer used to dump this config | `yaml` / `toml` / `...` |
| **target** | Path of dumped config |
| **value** | Specification of this config |

```
config := { key: value, ... }

value := <literal_type> # random generated value with no constraints, e.g. [u]int, bool, string
       | <literal> # literal value of int, float, string
       | <list_of_choices>
       | <full_value_decl>
       | <struct_map> # struct_map as {...} is a shortcut of {type: sturct, value: ...}

list_of_choices := []
                 | [<value>] :: <list_of_choices>

full_value_decl := {type: <type>[, when: <condition>][, arguments]}
```

## Type of config values
1. Basic types
   1. string
   2. bool
   3. int
   
      3.1. uint is an alias of int with minimal value 0
   4. float
   5. time (`s`/`m`/`h`)
   6. size (`B`/`KB`/`MB`/`...`)
2. Struct

See `examples/matrix-tidb.yaml` for an example of TiDB.

## TODO
### Serializer
- [x] toml serializer
- [ ] yaml serializer
- [ ] SQL serializer
- [ ] command line argument serializer
### Generator
- [x] simple recursive random generator
- [ ] non-recursive generator that supports dependency
### Random
- [x] random generate value
- [ ] specific random seed to have same output with same seed, to support continuous testing
