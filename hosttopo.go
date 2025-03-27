package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/EwenLan/vanadium-schedule/slog"
)

type VirtualMachineInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type VMTypeInfo struct {
	Name        string `json:"name"`
	Moveable    bool   `json:"moveable"`
	Core        int32  `json:"core"`
	LargeMemory int32  `json:"largeMemory"`
}

type NumaInfo struct {
	Name            string               `json:"name"`
	Core            int32                `json:"core"`
	LargeMemory     int32                `json:"largeMemory"`
	VitrualMachines []VirtualMachineInfo `json:"vms"`
}

type HostInfo struct {
	Name  string     `json:"name"`
	Numas []NumaInfo `json:"numas"`
}

type HostTopo struct {
	Hosts   []HostInfo   `json:"hosts"`
	VmTypes []VMTypeInfo `json:"vmTypes"`
}

var (
	vmTypesDict map[string]VMTypeInfo
	once        sync.Once
)

func getVmTypesDict() map[string]VMTypeInfo {
	once.Do(func() {
		vmTypesDict = make(map[string]VMTypeInfo)
	})
	return vmTypesDict
}

func (h *HostTopo) Load(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		slog.Errorf("fail to read file %s: %s", filepath, err)
		return err
	}
	if len(data) == 0 {
		slog.Errorf("file %s is empty", filepath)
		return fmt.Errorf("file %s is empty", filepath)
	}
	err = json.Unmarshal(data, h)
	if err != nil {
		slog.Errorf("fail to unmarshal file %s: %s", filepath, err)
		return err
	}
	for i := range h.VmTypes {
		vmType := &h.VmTypes[i]
		getVmTypesDict()[vmType.Name] = *vmType
	}
	return nil
}

func (n *NumaInfo) getRestCore() int32 {
	var core int32 = 0
	for _, vm := range n.VitrualMachines {
		vmType, ok := getVmTypesDict()[vm.Type]
		if !ok {
			slog.Errorf("fail to get vm type %s", vm.Type)
			continue
		}
		core += vmType.Core
	}
	return n.Core - core
}

func (n *NumaInfo) getRestLargeMemory() int32 {
	var largeMemory int32 = 0
	for _, vm := range n.VitrualMachines {
		vmType, ok := getVmTypesDict()[vm.Type]
		if !ok {
			slog.Errorf("fail to get vm type %s", vm.Type)
			continue
		}
		largeMemory += vmType.LargeMemory
	}
	return n.LargeMemory - largeMemory
}
