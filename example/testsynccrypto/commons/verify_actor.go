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
	"github.com/Ontology/crypto"
	"github.com/ontio/ontology-eventbus/actor"
	"time"
)

type VerifyActor struct {
	Count int
}

func (s *VerifyActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")

	case *VerifyRequest:
		start := time.Now().UnixNano()
		s.Count++
		crypto.SetAlg("")
		buf := bytes.NewBuffer(msg.PublicKey)
		pubKey := new(crypto.PubKey)
		err := pubKey.DeSerialize(buf)
		if err != nil {
			fmt.Println("DeSerialize failed.", err)
		}
		err = crypto.Verify(*pubKey, msg.Data, msg.Signature)
		//fmt.Println(context.Self().Id, "done verifying...")
		if err != nil {
			fmt.Println("verify error :", err)
			end := time.Now().UnixNano()
			vrftime := end - start
			response := &VerifyResponse{Seq: msg.Seq, Result: false, ErrorMsg: err.Error(), Vrftime: vrftime}
			context.Sender().Tell(response)
		} else {
			end := time.Now().UnixNano()
			vrftime := end - start
			response := &VerifyResponse{Seq: msg.Seq, Result: true, ErrorMsg: "", Vrftime: vrftime}
			context.Sender().Tell(response)
		}
	default:
		fmt.Printf("---unknown message%v\n", msg)
	}
}
