
## 2014 November 12

With a new RocksDB:

```Shell
go test
38: error: 'rocksdb_options_set_block_size' undeclared (first use in this function)
38: error: 'rocksdb_options_set_filter_policy' undeclared (first use in this function)
38: error: 'rocksdb_options_set_block_restart_interval' undeclared (first use in this function)
38: error: 'rocksdb_options_set_cache' undeclared (first use in this function)
FAIL	github.csv.comcast.com/jsteph206/gorocksdb [build failed]
```

Re `rocksdb_options_set_block_size`:

In `rocksdb/include/options.h`:

```C
  // Block-based table related options are moved to BlockBasedTableOptions.
  // Related options that were originally here but now moved include:
  //   no_block_cache
  //   block_cache
  //   block_cache_compressed
  //   block_size
  //   block_size_deviation
  //   block_restart_interval
  //   filter_policy
  //   whole_key_filtering
  // If you'd like to customize some of these options, you will need to
  // use NewBlockBasedTableFactory() to construct a new table factory.
```

For now, just commenting out this option.

For most of the other errors, same deal.



