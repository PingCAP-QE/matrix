# Configs of Matrix
The config file of Matrix tool is a yaml file, formatted as followed:
```yaml
config:
  - <config_template>
```

Each `config_template` is a set of configs that will be generated as one entity (i.e. file).

| Field | Description | Sample Value |
|:------|:------------------|:--------------|
| **tag** | Tag for this config |
| **serializer** | Serializer used to dump this config | `yaml` / `toml` / `stmt` / `...` |
| **target** | Path of dumped config |
| **value** | Specification of this config |

Each `value` is a map with its type and some other constraints,
some grammar sugars can be used for simplification.

```
value := { type: <type>[, when: <condition>][, arguments] }

config := { key: value, ... }

value := <literal_type> # sugar for random generated value with no constraints, e.g. [u]int, bool, string
       | <literal> # sugar for literal value (int, float, string)
       | <list_of_choices> # sugar for type `choice`
       | <full_value_decl>
       | <struct_map> # sugar for map

list_of_choices := []
                 | [<value>] :: <list_of_choices>

full_value_decl := {type: <type>[, when: <condition>][, arguments]}
```

## Type of config values
```
type := bool
      | int
      | string
      | float
      | time
      | size
      | time
      | map | struct
      | list
      | choice | choice_n
size := <num>{kb,mb,gb,...}
time := [<num>h][<num>m][<num>s] # at least one part should exist
```

## Arguments
| Type | Field | Description |
|:------|:------------------|:--------------|
| **bool** | `value` | Optional, value to use instead of randomly generate. |
| **string** | `value` | Optional, value to use instead of randomly generate. |
| **list** | `value` | Value of list |
| **map** | `value` | Value of map |
| **int** | `range` | List, [start, end], end is optional |
| **float** | `range` | List, [start, end], end is optional |
| **size** | `range` | List, [start, end], end is optional |
| **time** | `range` | List, [start, end], end is optional |
| **choice** | `value` | List, one random element will be selected. |
| **choice_n** | `value` | List, `n` random element will be selected. |
| **choice_n** | `n` | int, `n` of `choice_n` |
| **choice_n** | `sep` | Optional, string, if given and all elements in list is of type string, this will be used to join the selected elements. |

## Examples
See `examples/matrix-tidb-tikv-pd.yaml` for an example of TiDB.
