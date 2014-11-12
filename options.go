package rocksdb

// #cgo LDFLAGS: -lrocksdb
// #include "rocksdb/c.h"
import "C"

import "log"

// CompressionOpt is a value for Options.SetCompression.
type CompressionOpt int

// Known compression arguments for Options.SetCompression.
const (
	NoCompression     = CompressionOpt(0)
	SnappyCompression = CompressionOpt(1)
)

// Options represent all of the available options when opening a database with
// Open. Options should be created with NewOptions.
//
// It is usually with to call SetCache with a cache object. Otherwise, all
// data will be read off disk.
//
// To prevent memory leaks, Close must be called on an Options when the
// program no longer needs it.
type Options struct {
	Opt      *C.rocksdb_options_t
	readOnly bool
}

// ReadOptions represent all of the available options when reading from a
// database.
//
// To prevent memory leaks, Close must called on a ReadOptions when the
// program no longer needs it.
type ReadOptions struct {
	Opt *C.rocksdb_readoptions_t
}

// WriteOptions represent all of the available options when writeing from a
// database.
//
// To prevent memory leaks, Close must called on a WriteOptions when the
// program no longer needs it.
type WriteOptions struct {
	Opt *C.rocksdb_writeoptions_t
}

// NewOptions allocates a new Options object.
func NewOptions() *Options {
	opt := C.rocksdb_options_create()
	return &Options{opt, false}
}

// NewReadOptions allocates a new ReadOptions object.
func NewReadOptions() *ReadOptions {
	opt := C.rocksdb_readoptions_create()
	return &ReadOptions{opt}
}

// NewWriteOptions allocates a new WriteOptions object.
func NewWriteOptions() *WriteOptions {
	opt := C.rocksdb_writeoptions_create()
	return &WriteOptions{opt}
}

// Close deallocates the Options, freeing its underlying C struct.
func (o *Options) Close() {
	C.rocksdb_options_destroy(o.Opt)
}

func (o *Options) SetLogLevel(n int) {
	C.rocksdb_options_set_info_log_level(o.Opt, C.int(n))
}

// Tthe info LOG dir.
// If it is empty, the log files will be in the same dir as data.
// If it is non empty, the log files will be in the specified dir,
// and the db data dir's absolute path will be used as the log file
// name's prefix.
func (o *Options) SetLogDir(dir string) {
	C.rocksdb_options_set_db_log_dir(o.Opt, C.CString(dir))
}

// The absolute dir path for write-ahead logs (WAL).
// If it is empty, the log files will be in the same dir as data,
//   dbname is used as the data dir by default
// If it is non empty, the log files will be in kept the specified dir.
// When destroying the db,
//   all log files in wal_dir and the dir itself is deleted
func (o *Options) SetWalDir(dir string) {
	C.rocksdb_options_set_wal_dir(o.Opt, C.CString(dir))
}

// If not zero, dump rocksdb.stats to LOG every stats_dump_period_sec
// Default: 3600 (1 hour)
func (o *Options) SetStatsDumpPeriod(secs uint) {
	C.rocksdb_options_set_stats_dump_period_sec(o.Opt, C.uint(secs))
}

// Target file size for compaction.
// target_file_size_base is per-file size for level-1.
// Target file size for level L can be calculated by
// target_file_size_base * (target_file_size_multiplier ^ (L-1))
// For example, if target_file_size_base is 2MB and
// target_file_size_multiplier is 10, then each file on level-1 will
// be 2MB, and each file on level 2 will be 20MB,
// and each file on level-3 will be 200MB.
// by default target_file_size_base is 2MB.
func (o *Options) SetTargetFileSizeBase(n uint64) {
	C.rocksdb_options_set_target_file_size_base(o.Opt, C.uint64_t(n))
}

// by default target_file_size_multiplier is 1, which means
// by default files in different levels will have similar size.
func (o *Options) SetTargetFileSizeMultiplier(n int) {
	C.rocksdb_options_set_target_file_size_multiplier(o.Opt, C.int(n))
}

// Allows OS to incrementally sync files to disk while they are being
// written, asynchronously, in the background.
// Issue one request for every bytes_per_sync written. 0 turns it off.
// Default: 0
func (o *Options) SetBytesPerSync(n uint64) {
	C.rocksdb_options_set_bytes_per_sync(o.Opt, C.uint64_t(n))
}

// Number of levels for this database
func (o *Options) SetNumLevels(n int) {
	C.rocksdb_options_set_num_levels(o.Opt, C.int(n))
}

// Number of files to trigger level-0 compaction. A value <0 means that
// level-0 compaction will not be triggered by number of files at all.
//
// Default: 4
func (o *Options) SetLevel0FileNumCompactionTrigger(n int) {
	C.rocksdb_options_set_level0_file_num_compaction_trigger(o.Opt, C.int(n))
}

// Maximum number of concurrent background compaction jobs, submitted to
// the default LOW priority thread pool.
// If you're increasing this, also consider increasing number of threads in
// LOW priority thread pool. For more information, see
// Env::SetBackgroundThreads
// Default: 1
func (o *Options) SetMaxBackgroundCompactions(n int) {
	C.rocksdb_options_set_max_background_compactions(o.Opt, C.int(n))
}

// the HIGH priority thread pool.
//
// By default, all background jobs (major compaction and memtable flush) go
// to the LOW priority pool. If this option is set to a positive number,
// memtable flush jobs will be submitted to the HIGH priority pool.
// It is important when the same Env is shared by multiple db instances.
// Without a separate pool, long running major compaction jobs could
// potentially block memtable flush jobs of other db instances, leading to
// unnecessary Put stalls.
//
// If you're increasing this, also consider increasing number of threads in
// HIGH priority thread pool. For more information, see
// Env::SetBackgroundThreads
// Default: 1
func (o *Options) SetMaxBackgroundFlushes(n int) {
	C.rocksdb_options_set_max_background_flushes(o.Opt, C.int(n))
}

// Allow the OS to mmap file for reading sst tables. Default: false
func (o *Options) SetAllowMMapReads(b bool) {
	C.rocksdb_options_set_allow_mmap_reads(o.Opt, boolToUchar(b))
}

// Allow the OS to mmap file for writing. Default: false
func (o *Options) SetAllowMMapWrites(b bool) {
	C.rocksdb_options_set_allow_mmap_writes(o.Opt, boolToUchar(b))
}

// Data being read from file storage may be buffered in the OS
// Default: true
func (o *Options) SetAllowOSBuffer(b bool) {
	C.rocksdb_options_set_allow_os_buffer(o.Opt, boolToUchar(b))
}

// The maximum number of write buffers that are built up in memory.
// The default and the minimum number is 2, so that when 1 write buffer
// is being flushed to storage, new writes can continue to the other
// write buffer.
// Default: 2
func (o *Options) SetMaxWriteBufferNumber(n int) {
	C.rocksdb_options_set_max_write_buffer_number(o.Opt, C.int(n))
}

// The minimum number of write buffers that will be merged together
// before writing to storage.  If set to 1, then
// all write buffers are fushed to L0 as individual files and this increases
// read amplification because a get request has to check in all of these
// files. Also, an in-memory merge may result in writing lesser
// data to storage if there are duplicate records in each of these
// individual write buffers.  Default: 1
func (o *Options) SetMinWriteBufferNumberToMerge(n int) {
	C.rocksdb_options_set_min_write_buffer_number_to_merge(o.Opt, C.int(n))
}

func (o *Options) SetReadOnly(b bool) {
	o.readOnly = b
}

// If true, then the contents of data files are not synced
// to stable storage. Their contents remain in the OS buffers till the
// OS decides to flush them. This option is good for bulk-loading
// of data. Once the bulk-loading is complete, please issue a
// sync to the OS to flush all dirty buffesrs to stable storage.
// Default: false
func (o *Options) SetDisableDataSync(b bool) {
	C.rocksdb_options_set_disable_data_sync(o.Opt, C.int(boolToUchar(b)))
}

// SetComparator sets the comparator to be used for all read and write
// operations.
//
// The comparator that created a database must be the same one (technically,
// one with the same name string) that is used to perform read and write
// operations.
//
// The default comparator is usually sufficient.
func (o *Options) SetComparator(cmp *C.rocksdb_comparator_t) {
	C.rocksdb_options_set_comparator(o.Opt, cmp)
}

// SetErrorIfExists, if passed true, will cause the opening of a database that

// already exists to throw an error.
func (o *Options) SetErrorIfExists(error_if_exists bool) {
	eie := boolToUchar(error_if_exists)
	C.rocksdb_options_set_error_if_exists(o.Opt, eie)
}

// SetCache places a cache object in the database when a database is opened.
//
// This is usually wise to use. See also ReadOptions.SetFillCache.
func (o *Options) SetCache(cache *Cache) {
	// ToDo: SetCache
	// C.rocksdb_options_set_cache(o.Opt, cache.Cache)
	log.Println("Warning: SetCache not currently implemented")
}

// SetEnv sets the Env object for the new database handle.
func (o *Options) SetEnv(env *Env) {
	C.rocksdb_options_set_env(o.Opt, env.Env)
}

func (e *Env) SetBackgroundThreads(n int) {
	C.rocksdb_env_set_background_threads(e.Env, C.int(n))
}

func (e *Env) SetHighPriorityBackgroundThreads(n int) {
	C.rocksdb_env_set_high_priority_background_threads(e.Env, C.int(n))
}

// SetInfoLog sets a *C.rocksdb_logger_t object as the informational logger
// for the database.
func (o *Options) SetInfoLog(log *C.rocksdb_logger_t) {
	C.rocksdb_options_set_info_log(o.Opt, log)
}

// SetWriteBufferSize sets the number of bytes the database will build up in
// memory (backed by an unsorted log on disk) before converting to a sorted
// on-disk file.
func (o *Options) SetWriteBufferSize(s int) {
	C.rocksdb_options_set_write_buffer_size(o.Opt, C.size_t(s))
}

// SetParanoidChecks, when called with true, will cause the database to do
// aggressive checking of the data it is processing and will stop early if it
// detects errors.
//
// See the LevelDB documentation docs for details.
func (o *Options) SetParanoidChecks(pc bool) {
	C.rocksdb_options_set_paranoid_checks(o.Opt, boolToUchar(pc))
}

// SetMaxOpenFiles sets the number of files than can be used at once by the
// database.
//
// See the LevelDB documentation for details.
func (o *Options) SetMaxOpenFiles(n int) {
	C.rocksdb_options_set_max_open_files(o.Opt, C.int(n))
}

// SetBlockSize sets the approximate size of user data packed per block.
//
// The default is roughly 4096 uncompressed bytes. A better setting depends on
// your use case. See the LevelDB documentation for details.
func (o *Options) SetBlockSize(s int) {
	// ToDo: SetBlockSize
	// C.rocksdb_options_set_block_size(o.Opt, C.size_t(s))
	log.Println("Warning: SetBlockSize not currently implemented")
}

// SetBlockRestartInterval is the number of keys between restarts points for
// delta encoding keys.
//
// Most clients should leave this parameter alone. See the LevelDB
// documentation for details.
func (o *Options) SetBlockRestartInterval(n int) {
	// ToDo: SetBlockRestartInterval
	// C.rocksdb_options_set_block_restart_interval(o.Opt, C.int(n))
	log.Println("Warning: SetBlockRestartInterval not currently implemented")
}

// SetCompression sets whether to compress blocks using the specified
// compresssion algorithm.
//
// The default value is SnappyCompression and it is fast enough that it is
// unlikely you want to turn it off. The other option is NoCompression.
//
// If the LevelDB library was built without Snappy compression enabled, the
// SnappyCompression setting will be ignored.
func (o *Options) SetCompression(t CompressionOpt) {
	C.rocksdb_options_set_compression(o.Opt, C.int(t))
}

// SetCreateIfMissing causes Open to create a new database on disk if it does
// not already exist.
func (o *Options) SetCreateIfMissing(b bool) {
	C.rocksdb_options_set_create_if_missing(o.Opt, boolToUchar(b))
}

// By default, RocksDB uses only one background thread for flush and
// compaction. Calling this function will set it up such that total of
// `total_threads` is used. Good value for `total_threads` is the number of
// cores. You almost definitely want to call this function if your system is
// bottlenecked by RocksDB.
func (o *Options) IncreaseParallelism(n int) {
	C.rocksdb_options_increase_parallelism(o.Opt, C.int(n))
}

// SetFilterPolicy causes Open to create a new database that will uses filter
// created from the filter policy passed in.
func (o *Options) SetFilterPolicy(fp *FilterPolicy) {
	// var policy *C.rocksdb_filterpolicy_t
	// if fp != nil {
	// 	policy = fp.Policy
	// }
	// ToDo: SetFilterPolicy
	// C.rocksdb_options_set_filter_policy(o.Opt, policy)
	log.Println("Warning: SetFilterPolicy not currently implemented")
}

// Close deallocates the ReadOptions, freeing its underlying C struct.
func (ro *ReadOptions) Close() {
	C.rocksdb_readoptions_destroy(ro.Opt)
}

// SetVerifyChecksums controls whether all data read with this ReadOptions
// will be verified against corresponding checksums.
//
// It defaults to false. See the LevelDB documentation for details.
func (ro *ReadOptions) SetVerifyChecksums(b bool) {
	C.rocksdb_readoptions_set_verify_checksums(ro.Opt, boolToUchar(b))
}

// SetFillCache controls whether reads performed with this ReadOptions will
// fill the Cache of the server. It defaults to true.
//
// It is useful to turn this off on ReadOptions for DB.Iterator (and DB.Get)
// calls used in offline threads to prevent bulk scans from flushing out live
// user data in the cache.
//
// See also Options.SetCache
func (ro *ReadOptions) SetFillCache(b bool) {
	C.rocksdb_readoptions_set_fill_cache(ro.Opt, boolToUchar(b))
}

// SetSnapshot causes reads to provided as they were when the passed in
// Snapshot was created by DB.NewSnapshot. This is useful for getting
// consistent reads during a bulk operation.
//
// See the LevelDB documentation for details.
func (ro *ReadOptions) SetSnapshot(snap *Snapshot) {
	var s *C.rocksdb_snapshot_t
	if snap != nil {
		s = snap.snap
	}
	C.rocksdb_readoptions_set_snapshot(ro.Opt, s)
}

// Close deallocates the WriteOptions, freeing its underlying C struct.
func (wo *WriteOptions) Close() {
	C.rocksdb_writeoptions_destroy(wo.Opt)
}

// SetSync controls whether each write performed with this WriteOptions will
// be flushed from the operating system buffer cache before the write is
// considered complete.
//
// If called with true, this will signficantly slow down writes. If called
// with false, and the host machine crashes, some recent writes may be
// lost. The default is false.
//
// See the LevelDB documentation for details.
func (wo *WriteOptions) SetSync(b bool) {
	C.rocksdb_writeoptions_set_sync(wo.Opt, boolToUchar(b))
}

func (wo *WriteOptions) DisableWAL(b bool) {
	n := 0
	if b {
		n = 1
	}
	C.rocksdb_writeoptions_disable_WAL(wo.Opt, C.int(n))
}
