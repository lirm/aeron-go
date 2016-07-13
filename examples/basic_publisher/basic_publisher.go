package main

import (
	"github.com/lirm/aeron-go/aeron"
	"github.com/lirm/aeron-go/aeron/buffers"
	"github.com/lirm/aeron-go/aeron/term"
	"log"
	"time"
)

func main() {
	ctx := new(aeron.Context).AeronDir("/tmp").MediaDriverTimeout(10000)

	a := aeron.Connect(ctx)

	pubId := a.AddPublication("aeron:udp?endpoint=localhost:40123", 10)
	log.Printf("Publication ID: %d, aeron: %v\n", pubId, a)

	//time.Sleep(time.Second)

	publication := a.FindPublication(pubId)
	for publication == nil {
		publication = a.FindPublication(pubId)
	}

	log.Printf("Publication found %v", publication)

	// t.Logf("+++ Starting to publish +++\n")
	for {
		message := "this is a message"
		srcBuffer := buffers.MakeAtomic(([]byte)(message))
		ret := publication.Offer(srcBuffer, 0, int32(len(message)), term.DEFAULT_RESERVED_VALUE_SUPPLIER)
		switch ret {
		case aeron.NOT_CONNECTED:
			log.Printf("not connected yet")
		case aeron.BACK_PRESSURED:
			log.Printf("back pressured")
		default:
			if ret < 0 {
				log.Printf("Unrecognized code: %d", ret)
			} else {
				log.Printf("success!")
			}
		}
		//log.Printf("Publication.Offer returned %d", ret)
		time.Sleep(time.Second)
	}
}
