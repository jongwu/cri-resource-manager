package sysfs

import (
	"math/rand"
//	"path/filepath"
	"strconv"
	"runtime"
	"testing"
	"time"

	logger "github.com/intel/cri-resource-manager/pkg/log"
	idset "github.com/intel/goresctrl/pkg/utils"
)

func TestCcxIDs(t *testing.T) {
	cpus0 := idset.NewIDSetFromIntSlice([]int{0, 1, 2, 3}...)
	ccx0 := ccx {
		id:		idset.ID(0),
		pkg:		idset.ID(0),
		die:		idset.ID(0),
		node:		idset.ID(0),
		cpus:		cpus0,
		cpulist:	[]int{0, 1, 2, 3},
	}
	
	cpus1 := idset.NewIDSetFromIntSlice([]int{4, 5, 6, 7}...)
	ccx1 := ccx {
		id:		idset.ID(1),
		pkg:		idset.ID(1),
		die:		idset.ID(1),
		node:		idset.ID(1),
		cpus:		cpus1,
		cpulist:	[]int{4, 5, 6, 7},
	}
	stm := system {
		ccxs: make(map[idset.ID]*ccx),
	}
	stm.ccxs[ccx0.id] = &ccx0
	stm.ccxs[ccx1.id] = &ccx1
	expected := []idset.ID{0, 1}
	got := stm.CcxIDs()
	if got == nil {
		t.Errorf("Test CcxIDs: got nil\n")
		return
	}
	for id, v := range expected {
		if v != got[id] {
			t.Errorf("expected: %v, got: %v\n", expected, got)
			return
		}
	}
}

func TestDiscoverCCXs(t *testing.T) {
	sys := &system {
		Logger:  logger.NewLogger("sysfs"),
		path: "/sys",
		offline: idset.NewIDSet(),
	}
	if err := sys.discoverCPUs(); err != nil {
		t.Errorf("discoverCPUs fail: %v\n", err)
		return
	}
	if err := sys.discoverCCXs(); err != nil {
		t.Errorf("discoverCCXs fail: %v\n", err)
		return
	}
	cpuNum := runtime.NumCPU()
	rand.Seed(time.Now().UnixNano())
	checkCpu := rand.Intn(cpuNum)
	cpustr := strconv.Itoa(checkCpu)
	path := "/sys/devices/system/cpu/cpu" + cpustr
	var cacheId string
	if _, err := readSysfsEntry(path, "/cache/index3/id", &cacheId); err != nil {
		t.Errorf("read sysfs entry fail: %v\n", err)
		return
	}
	var expected int
	if id, err := strconv.Atoi(cacheId); err != nil {
		t.Errorf("convert string to int fail: %v\n", err)
		return
	} else {
		expected = id
	}
	got := sys.cpus[idset.ID(checkCpu)].ccx
	if got != expected {
		t.Errorf("expected: %d, got: %d\n", expected, got)
	}
}
