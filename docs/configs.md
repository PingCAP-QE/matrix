This file describes some important options fuzzing should pay attention to.

# Mysql System Variables
## `sql_mode`
### `HIGH_NOT_PRECEDENCE`
The behavior of `not` seems different with mysql, it (and `MySQL323`, `MYSQL40`) breaks `create table if not exists`.

See https://github.com/pingcap/tidb/issues/18729
