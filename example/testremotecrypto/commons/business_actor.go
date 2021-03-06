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
package commons

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Ontology/crypto"
	"github.com/ontio/ontology-eventbus/actor"
	"github.com/ontio/ontology-eventbus/eventhub"
)

const (
	SigTOPIC    string = "SIGTOPIC"
	VerifyTOPIC string = "VERIFYTOPIC"
	SetTOPIC    string = "SETTOPIC"
)

const loop1 = 1
const Loop2 = 250000

type BusynessActor struct {
	Datas      map[string][]byte
	privatekey []byte
	pubkey     crypto.PubKey
	start      int64
	respCount  int
	WgStop     *sync.WaitGroup
	LatencySum *int64
	VrftimeSum *int64
}

func (s *BusynessActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *RunMsg:
		fmt.Println("Recieve runMsg")
		crypto.SetAlg("")
		bb := bytes.NewBuffer([]byte(""))
		for i := 0; i < 100; i++ {
			bb.WriteString("1234567890")
		}

		privKey, pubkey, _ := crypto.GenKeyPair()

		s.privatekey = privKey
		s.pubkey = pubkey

		setPrivMsg := &SetPrivKey{PrivKey: privKey}

		setEvent := &eventhub.Event{Topic: "SETTOPIC", Publisher: context.Self(), Message: setPrivMsg, Policy: eventhub.PublishPolicyAll}

		eventhub.GlobalEventHub.Publish(setEvent)

		s.start = time.Now().UnixNano()
		for i := 0; i < loop1; i++ {
			idx := strconv.Itoa(i)
			bb.WriteString(idx)
			data := bb.Bytes()
			sigMsg := &SignRequest{Seq: idx, Data: data}
			s.Datas[idx] = data
			sigEvent := &eventhub.Event{Topic: "SIGTOPIC", Publisher: context.Self(), Message: sigMsg, Policy: eventhub.PublishPolicyRoundRobin}
			eventhub.GlobalEventHub.Publish(sigEvent)
		}

	case *SignResponse:
		seq := msg.Seq
		sig := msg.Signature

		buf := bytes.NewBuffer([]byte(""))
		err := s.pubkey.Serialize(buf)
		if err != nil {
			fmt.Println("ERROR Serialize publickey: ", err)
		}
		pubKeyBytes := buf.Bytes()

		start := time.Now()
		for i := 0; i < Loop2; i++ {
			timeStamp := time.Now().UnixNano()
			vfr := &VerifyRequest{Signature: sig, Data: s.Datas[seq], PublicKey: pubKeyBytes, Seq: seq, Timestamp: timeStamp}
			vrfEvent := &eventhub.Event{Topic: "VERIFYTOPIC", Publisher: context.Self(), Message: vfr, Policy: eventhub.PublishPolicyRoundRobin}
			eventhub.GlobalEventHub.Publish(vrfEvent)
		}
		elapsed := time.Since(start)
		fmt.Printf("Elapsed %s\n", elapsed)

	case *VerifyResponse:
		s.respCount++
		*s.VrftimeSum = *s.VrftimeSum + msg.Vrftime
		*s.LatencySum = *s.LatencySum + (time.Now().UnixNano() - msg.Timestamp)

		if s.respCount%10000 == 0 {
			fmt.Printf("%d ", (time.Now().UnixNano()-msg.Timestamp)/1000000)
			fmt.Println(s.respCount)
		}
		if s.respCount == Loop2 {
			s.WgStop.Done()
		}
	default:
		fmt.Printf("unknown msg %v\n", msg)
	}
}
