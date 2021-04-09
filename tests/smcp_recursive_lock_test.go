// Copyright 2020 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"maistra/util"
	"testing"
)

func testPilotLock() {
	if _ , err := util.Shell("oc rsh -n istio-system -c discovery $(oc get pods -n istio-system -l app=pilot --no-headers | awk '{print $1}') curl -v http://localhost:8080/debug/cdsz") ; err != nil {
		t.Errorf("Pilot is Locked")
	}
}


func queryPilot() {
	util.Shell("PILOT=$(oc get -n istio-system pods -l app=pilot --no-headers | awk '{print $1}') ; count=1 ; while : ; do oc rsh -n istio-system -c discovery ${PILOT} curl -o /dev/null -s -w \"%{http_code}\n\" http://localhost:8080/debug/edsz ; echo $(date): loop $count ; count=$(($count + 1)); done")
}

func rolloutIngressGateway() {
	util.Shell("count=1 ; while : ; do oc -n istio-system rollout restart deployment istio-ingressgateway >/dev/null ; echo $(date): Restarted count $count ; count=$(($count + 1)); sleep 5 ; done")
}


func TestRecursiveLock(t *testing.T) {


	testPilotLock() 
	go queryPilot()

	go rolloutIngressGateway()
	
	testPilotLock()
}