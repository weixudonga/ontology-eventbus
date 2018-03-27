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
	"bytes"
	"fmt"
	"runtime"
	"time"

	"github.com/Ontology/crypto"
	"github.com/ontio/ontology-eventbus/actor"
	"github.com/ontio/ontology-eventbus/example/testSyncCrypto/commons"
	"github.com/ontio/ontology-eventbus/zmqremote"
)

const loop = 10000

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	zmqremote.Start("127.0.0.1:9081")

	var vrftimeSum int64
	crypto.SetAlg("")
	privKey, pubkey, _ := crypto.GenKeyPair()

	buf := bytes.NewBuffer([]byte(""))
	err := pubkey.Serialize(buf)
	if err != nil {
		fmt.Println("ERROR Serialize publickey: ", err)
	}
	pubKeyBytes := buf.Bytes()

	bb := bytes.NewBuffer([]byte(""))
	for i := 0; i < 100; i++ {
		bb.WriteString("1234567890")
	}

	signature, err := crypto.Sign(privKey, bb.Bytes())
	if err != nil {
		fmt.Println("sign error: ", err)
	}

	vrfactor := actor.NewPID("127.0.0.1:9080", "verify")

	start := time.Now().UnixNano()
	for i := 0; i < loop; i++ {
		vfr := &commons.VerifyRequest{Signature: signature, Data: bb.Bytes(), PublicKey: pubKeyBytes}
		result, err := vrfactor.RequestFuture(vfr, 1*time.Second).Result()
		if err != nil {
			fmt.Println("send vrf error: ", err)
		}
		vrftimeSum = vrftimeSum + result.(*commons.VerifyResponse).Vrftime
	}
	end := time.Now().UnixNano()

	latency := float64(end-start) / loop / 1000000
	vrftime := float64(vrftimeSum) / loop / 1000000

	fmt.Println("latency: ", latency)
	fmt.Println("vrftime: ", vrftime)
}
