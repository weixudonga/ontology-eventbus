/****************************************************
Copyright 2018 The ont-eventbus Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*****************************************************/

/***************************************************
Copyright 2016 https://github.com/AsynkronIT/protoactor-go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*****************************************************/
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/ontio/ontology-eventbus/actor"
	"github.com/ontio/ontology-eventbus/eventhub"
	"github.com/ontio/ontology-eventbus/example/testRemoteCrypto/commons"
	"github.com/ontio/ontology-eventbus/mailbox"
	"github.com/ontio/ontology-eventbus/remote"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	var wg sync.WaitGroup
	var vrftimeSum int64
	var latencySum int64
	remote.Start("127.0.0.1:9082")
	props := actor.FromProducer(func() actor.Actor {
		return &commons.BusynessActor{Datas: make(map[string][]byte), WgStop: &wg, VrftimeSum: &vrftimeSum, LatencySum: &latencySum}
	}).WithMailbox(mailbox.Bounded(100))

	bActor, _ := actor.SpawnNamed(props, "busi")

	signActor := actor.NewPID("127.0.0.1:9080", "sign")
	vfActor1 := actor.NewPID("127.0.0.1:9081", "verify1")
	//vfActor2 := actor.NewPID("172.26.127.136:9081", "verify2")
	//vfActor3 := actor.NewPID("172.26.127.138:9081", "verify3")

	eventhub.GlobalEventHub.Subscribe(commons.SetTOPIC, signActor)
	eventhub.GlobalEventHub.Subscribe(commons.SigTOPIC, signActor)
	eventhub.GlobalEventHub.Subscribe(commons.VerifyTOPIC, vfActor1)
	//eventhub.GlobalEventHub.Subscribe(commons.VerifyTOPIC,vfActor2)
	//eventhub.GlobalEventHub.Subscribe(commons.VerifyTOPIC,vfActor3)

	wg.Add(1)
	start := time.Now()

	bActor.Tell(&commons.RunMsg{})
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Elapsed %s\n", elapsed)
	x := int(float32(commons.Loop2) / (float32(elapsed) / float32(time.Second)))
	fmt.Printf("Msg per sec %v\n", x)
	vrftime := float64(vrftimeSum) / float64(commons.Loop2) / float64(1000000)
	latency := int(float64(latencySum) / float64(commons.Loop2) / float64(1000000))
	fmt.Printf("vrftime %f\n", vrftime)
	fmt.Printf("latency %d\n", latency)

	for {
		time.Sleep(1 * time.Second)
	}
}
