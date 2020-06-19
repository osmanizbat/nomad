package connect

import (
	"os"

	"github.com/hashicorp/nomad/e2e/e2eutil"
	"github.com/hashicorp/nomad/e2e/framework"
	"github.com/hashicorp/nomad/helper/uuid"
)

type ConnectE2ETest struct {
	framework.TC
	jobIds []string
}

func init() {
	// connect tests without Consul ACLs enabled
	framework.AddSuites(&framework.TestSuite{
		Component:   "Connect",
		CanRunLocal: true,
		Consul:      true,
		Cases: []framework.TestCase{
			new(ConnectE2ETest),
			new(ConnectClientStateE2ETest),
		},
	})

	// connect tests with Consul ACLs enabled
	framework.AddSuites(&framework.TestSuite{
		Component:   "ConnectACLs",
		CanRunLocal: false,
		Consul:      true,
		Parallel:    false,
		Cases: []framework.TestCase{
			new(ConnectACLsE2ETest),
		},
	})
}

func (tc *ConnectE2ETest) BeforeAll(f *framework.F) {
	e2eutil.WaitForLeader(f.T(), tc.Nomad())
	e2eutil.WaitForNodesReady(f.T(), tc.Nomad(), 2)
}

func (tc *ConnectE2ETest) AfterEach(f *framework.F) {
	if os.Getenv("NOMAD_TEST_SKIPCLEANUP") == "1" {
		return
	}

	for _, id := range tc.jobIds {
		tc.Nomad().Jobs().Deregister(id, true, nil)
	}
	tc.jobIds = []string{}
	tc.Nomad().System().GarbageCollect()
}

func connectJobID() string {
	id := uuid.Generate()
	jobID := "connect" + id[0:8]
	return jobID
}

// TestConnectDemo tests the demo job file from the Consul Connect Technology
// Preview.
//
// https://github.com/hashicorp/nomad/blob/v0.9.5/website/source/guides/integrations/consul-connect/index.html.md#run-the-connect-enabled-services
//
func (tc *ConnectE2ETest) TestConnectDemo(f *framework.F) {
	t := f.T()

	jobID := connectJobID()
	tc.jobIds = append(tc.jobIds, jobID)

	allocs := e2eutil.RegisterAndWaitForAllocs(t, tc.Nomad(), "connect/input/demo.nomad", jobID, "")
	allocIDs := e2eutil.AllocIDsFromAllocationListStubs(allocs)
	e2eutil.CheckServicesPassing(t, tc.Consul().Agent(), allocIDs)
}

func (tc *ConnectE2ETest) TestConnectNativeDemo(f *framework.F) {
	t := f.T()

	jobID := connectJobID()
	tc.jobIds = append(tc.jobIds, jobID)

	allocs := e2eutil.RegisterAndWaitForAllocs(t, tc.Nomad(), "connect/input/native-demo.nomad", jobID, "")
	allocIDs := e2eutil.AllocIDsFromAllocationListStubs(allocs)
	e2eutil.CheckServicesPassing(t, tc.Consul().Agent(), allocIDs)
}
