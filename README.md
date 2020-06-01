# Matrix

Matrix is a config fuzzer, i.e. generate a set of configurations for a program that need to be fuzzed.

## DSL

Matrix use yaml to describe the configurations need to be fuzzed.

## Type of configs
1. Domain specific config file, e.g. toml, yaml etc.;
2. Command line arguments;
3. Command (Query in Database) need to be run.

Note 3. could also be used for dynamic configuring (change config during runtime),
which will not be considered now.

## Config of Matrix
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

The semantic of field `value` is serializer-dependent
 to support non-trivial config like SQL-based ones.

```
config := <key>: <value>
        | <key>: {value: <value>[, when: <condition>]}

condition := <key> <op> <cond_operand>
op := ==
    | !=
    | in
    | not in

value := <literal>
       | {type: <literal_type>[, value: <literal>]}
       | {type: <numeric_and_time_type>[, range: [start, end]]}
       | {type: <struct_type>, value: {[<key>: <value>]*} | [<value>*]}
       | <list_of_choices>

list_of_choices := []
                 | [<value>] :: <list_of_choices>
                 | [{value: <value>, when: <condition>}] :: <list_of_choices>
```

## Type of config values
1. Basic types
   1. string
   2. bool
   3. int
   4. float
   5. time
2. Struct
todo: List for `repair-table-list`

See `test.yaml` for a not-yet-complete example.