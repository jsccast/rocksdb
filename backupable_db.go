package rocksdb

/*
#cgo LDFLAGS: -lrocksdb
#include <stdlib.h>
#include "rocksdb/c.h"
*/
import "C"

import (
	"unsafe"
)

type BackupEngine struct {
	Engine *C.rocksdb_backup_engine_t
}

type RestoreOptions struct {
	Opt *C.rocksdb_restore_options_t
}

type BackupEngineInfo struct {
	Info *C.rocksdb_backup_engine_info_t
}

func BackupEngineOpen(o *Options, path string) (*BackupEngine, error) {
	var errStr *C.char
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var backupEngine *C.rocksdb_backup_engine_t
	backupEngine = C.rocksdb_backup_engine_open(o.Opt, cPath, &errStr)
	if errStr != nil {
		gs := C.GoString(errStr)
		C.free(unsafe.Pointer(errStr))
		return nil, DatabaseError(gs)
	}
	return &BackupEngine{backupEngine}, nil
}

func (be *BackupEngine) CreateNewBackup(db *DB) error {
	var errStr *C.char

	C.rocksdb_backup_engine_create_new_backup(be.Engine, db.Ldb, &errStr)
	if errStr != nil {
		gs := C.GoString(errStr)
		C.free(unsafe.Pointer(errStr))
		return DatabaseError(gs)
	}
	return nil
}

func (be *BackupEngine) GetBackupInfo() *BackupEngineInfo {
	return &BackupEngineInfo{C.rocksdb_backup_engine_get_backup_info(be.Engine)}
}

func (bei *BackupEngineInfo) Count() int {
	return int(C.rocksdb_backup_engine_info_count(bei.Info))
}

func (bei *BackupEngineInfo) Timestamp(index int) int {
	return int(C.rocksdb_backup_engine_info_timestamp(bei.Info, C.int(index)))
}

func (bei *BackupEngineInfo) BackupId(index int) int {
	return int(C.rocksdb_backup_engine_info_backup_id(bei.Info, C.int(index)))
}

func (bei *BackupEngineInfo) NumberFiles(index int) int {
	return int(C.rocksdb_backup_engine_info_number_files(bei.Info, C.int(index)))
}

func (bei *BackupEngineInfo) Destroy() {
	C.rocksdb_backup_engine_info_destroy(bei.Info)
}

func (bei *BackupEngineInfo) Size(index int) int {
	return int(C.rocksdb_backup_engine_info_size(bei.Info, C.int(index)))
}

func CreateRestoreOptions() *RestoreOptions {
	return &RestoreOptions{C.rocksdb_restore_options_create()}
}

func (ro *RestoreOptions) Destroy() {
	C.rocksdb_restore_options_destroy(ro.Opt)
}

// If true, restore won't overwrite the existing log files in wal_dir. It will
// also move all log files from archive directory to wal_dir. Use this option
// in combination with BackupableDBOptions::backup_log_files = false for
// persisting in-memory databases.
// Default: false
func (ro *RestoreOptions) SetKeepLogFiles(v int) {
	C.rocksdb_restore_options_set_keep_log_files(ro.Opt, C.int(v))
}

func (be *BackupEngine) RestoreDbFromLatestBackup(dbDir string, walDir string, options *RestoreOptions) error {
	var errStr *C.char
	cDbDir := C.CString(dbDir)
	defer C.free(unsafe.Pointer(cDbDir))
	cWalDir := C.CString(walDir)
	defer C.free(unsafe.Pointer(cWalDir))

	C.rocksdb_backup_engine_restore_db_from_latest_backup(be.Engine, cDbDir, cWalDir, options.Opt, &errStr)
	if errStr != nil {
		gs := C.GoString(errStr)
		C.free(unsafe.Pointer(errStr))
		return DatabaseError(gs)
	}
	return nil
}

func (be *BackupEngine) Close() {
	C.rocksdb_backup_engine_close(be.Engine)
}
