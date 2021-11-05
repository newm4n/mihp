package handlers

import (
	"github.com/newm4n/mihp/central/server/mux"
)

var (
	PrefixPath = ""
)

func Routing(mux *mux.MyMux) {
	mux.AddRoute(PrefixPath+"/login", "POST", HandleLogin)
	mux.AddRoute(PrefixPath+"/refresh", "POST", HandleRefresh)
	//
	//mux.AddRoute(PrefixPath+"/probe", "POST", HandleProbeRegister)
	//mux.AddRoute(PrefixPath+"/probe/{probeid}", "GET", HandleProbePing)
	//mux.AddRoute(PrefixPath+"/probe/{probeid}", "POST", HandleProbeReport)
	//
	//mux.AddRoute(PrefixPath+"/users", "GET", HandleListUsers)
	//mux.AddRoute(PrefixPath+"/users", "POST", HandleCreateUser)
	//mux.AddRoute(PrefixPath+"/users/add", "POST", HandleCreateUser)
	//mux.AddRoute(PrefixPath+"/users/{useremail}", "GET", HandleGetUser)
	//mux.AddRoute(PrefixPath+"/users/{useremail}", "PUT", HandleUpdateUser)
	//mux.AddRoute(PrefixPath+"/users/{useremail}", "DELETE", HandleDeleteUser)
	//mux.AddRoute(PrefixPath+"/users/{useremail}/organizations", "DELETE", HandleListOrgByUser)
	//
	//mux.AddRoute(PrefixPath+"/organizations", "GET", HandleListOrganizations)
	//mux.AddRoute(PrefixPath+"/organizations", "POST", HandleCreateOrganization)
	//mux.AddRoute(PrefixPath+"/organizations/{orgId}", "GET", HandleGetOrganization)
	//mux.AddRoute(PrefixPath+"/organizations/{orgId}", "PUT", HandleUpdateOrganization)
	//mux.AddRoute(PrefixPath+"/organizations/{orgId}", "DELETE", HandleDeleteOrganization)
	//mux.AddRoute(PrefixPath+"/organizations/{orgId}/probes", "GET", HandleListProbesByOrganization)
	//
	//mux.AddRoute(PrefixPath+"/probedata", "GET", HandleListProbes)
	//mux.AddRoute(PrefixPath+"/probedata", "POST", HandleCreateProbe)
	//
	//mux.AddRoute(PrefixPath+"/probedata/{probeid}", "GET", HandleGetProbe)
	//mux.AddRoute(PrefixPath+"/probedata/{probeid}", "PUT", HandleUpdateProbe)
	//mux.AddRoute(PrefixPath+"/probedata/{probeid}", "DELETE", HandleDeleteProbe)
	//mux.AddRoute(PrefixPath+"/probedata/{probeid}/reqs", "GET", HandleListProbeRequestByProbe)
	//mux.AddRoute(PrefixPath+"/probedata/{probeid}/reqs", "POST", HandleCreateNewProbeRequest)
	//mux.AddRoute(PrefixPath+"/probedata/{probeid}/reqs/{requestid}", "GET", HandleGetProbeRequest)
	//mux.AddRoute(PrefixPath+"/probedata/{probeid}/reqs/{requestid}", "PUT", HandleUpdateProbeRequest)
	//mux.AddRoute(PrefixPath+"/probedata/{probeid}/reqs/{requestid}", "DELETE", HandleDeleteProbeRequest)

}
