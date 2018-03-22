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
	"github.com/ontio/ontology-eventbus/actor"
	"fmt"
	"github.com/Ontology/crypto"
)

type SignActor struct{
	PrivateKey []byte
}

func (s *SignActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")

	case *SetPrivKey:
		fmt.Println(context.Self().Id," set Privkey")
		s.PrivateKey = msg.PrivKey

	case *SignRequest:
		crypto.SetAlg("")
		fmt.Println(context.Self().Id," is signing")
		signature,err:=crypto.Sign(s.PrivateKey, msg.Data)
		if err!= nil {
			fmt.Println("sign error: ", err)
		}
		response := &SignResponse{Signature:signature,Seq:msg.Seq}
		fmt.Println(context.Self().Id," done signing")
		context.Sender().Request(response,context.Self())

	default:
		fmt.Println("unknown message")
	}
}