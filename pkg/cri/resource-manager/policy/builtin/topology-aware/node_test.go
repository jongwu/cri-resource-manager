package topologyaware

import (
	"testing"
	resapi "k8s.io/apimachinery/pkg/api/resource"
	policyapi "github.com/intel/cri-resource-manager/pkg/cri/resource-manager/policy"
	system "github.com/intel/cri-resource-manager/pkg/sysfs"
	idset "github.com/intel/goresctrl/pkg/utils"
)

func TestNewCcxNode(t *testing.T) {
	var numanode0 Node = &node{
		name:	"numa0",
		id:	0,
		kind:	NumaNode,
		depth:	2,
	}

	sys, err := system.DiscoverSystemAt("/sys")
	if err != nil {
		panic(err)
	}
	reserved, _ := resapi.ParseQuantity("750m")
	policyOptions := &policyapi.BackendOptions{
		Cache:  &mockCache{},
		System: sys,
		Reserved: policyapi.ConstraintSet{
			policyapi.DomainCPU: reserved,
		},
	}
	policy := CreateTopologyAwarePolicy(policyOptions).(*policy)
	ccxNode := policy.NewCcxNode(idset.ID(0), numanode0)
	if ccxNode.Name() != "CCX node #0" || ccxNode.Kind() != "ccx" {
		t.Errorf("expected Name: CCX node #0, got: %s, expected kind: ccx, got: %v", ccxNode.Name(), ccxNode.Kind())
	}
}
