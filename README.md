# Orchestrator Framework of Saga pattern
* Orchestrator's reponsibility is integrating distrubuted transaction in microservice, all upstream service just needs to maintain gRPC endpoint.
* The flows supported by Orchestrator including: Sync Flow, Async Flow(Using MQ), and throttling flow whick locks request, and queue them.
* Just needs to implement framework structure, it will be easy to register and scale transaction.
* We can change driver of Message Queue, Now we support RabbitMQ and NATS.
* A flow is combine by multiple Nodes, Node can be reuse by other flow.



## Contents

- [Orchestrator Framework of Saga pattern](#orchestrator-framework-of-saga-pattern)
  - [Contents](#contents)
  - [Quick Start](#quick-start)
  - [Structure](#structure)
  - [Concept](#concept)
  - [Task](#task)

## Quick Start

````sh
    
    export OCH_PATH=$(pwd)
    export POD_ID=local
    
    cp ./deployment/local/app.yaml app.yaml
    go test ./... 

    # start up related service 
    cd deployment/local
    docker-compose up -d
    cd ../..
    go run main.go

    # Hit it up
    grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "false"}' localhost:10000 pb.BookingService/HandleThrottlingBooking

````

## Structure
````
.
├── facade: Register Sync/Async transaction flows
├── node: Sync/Async transaction nodes
├── handler: Implements protobuf gRPC ServiceServer interface
├── topic: Register Message Queue Topics
````

## Concept
### Node/ Flow
[![flow-node](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggVERcbiAgU3RhcnRcbiAgQVtOb2RlIC0gQ3JlYXRlT3JkZXJdIC0tPiBCKE5vZGUgLSBDcmF0ZVBheW1lbnQpXG4gIEItLSBlcnJvciBoYXBwZW5kIC0tPkMoTm9kZSAtIENhbmNlbGxCb29raW5nKVxuICBTdGFydC0tPkFcbiAgQi0tPkZpbmlzaCIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0IiwidGhlbWVWYXJpYWJsZXMiOnsiYmFja2dyb3VuZCI6IndoaXRlIiwicHJpbWFyeUNvbG9yIjoiI0VDRUNGRiIsInNlY29uZGFyeUNvbG9yIjoiI2ZmZmZkZSIsInRlcnRpYXJ5Q29sb3IiOiJoc2woODAsIDEwMCUsIDk2LjI3NDUwOTgwMzklKSIsInByaW1hcnlCb3JkZXJDb2xvciI6ImhzbCgyNDAsIDYwJSwgODYuMjc0NTA5ODAzOSUpIiwic2Vjb25kYXJ5Qm9yZGVyQ29sb3IiOiJoc2woNjAsIDYwJSwgODMuNTI5NDExNzY0NyUpIiwidGVydGlhcnlCb3JkZXJDb2xvciI6ImhzbCg4MCwgNjAlLCA4Ni4yNzQ1MDk4MDM5JSkiLCJwcmltYXJ5VGV4dENvbG9yIjoiIzEzMTMwMCIsInNlY29uZGFyeVRleHRDb2xvciI6IiMwMDAwMjEiLCJ0ZXJ0aWFyeVRleHRDb2xvciI6InJnYig5LjUwMDAwMDAwMDEsIDkuNTAwMDAwMDAwMSwgOS41MDAwMDAwMDAxKSIsImxpbmVDb2xvciI6IiMzMzMzMzMiLCJ0ZXh0Q29sb3IiOiIjMzMzIiwibWFpbkJrZyI6IiNFQ0VDRkYiLCJzZWNvbmRCa2ciOiIjZmZmZmRlIiwiYm9yZGVyMSI6IiM5MzcwREIiLCJib3JkZXIyIjoiI2FhYWEzMyIsImFycm93aGVhZENvbG9yIjoiIzMzMzMzMyIsImZvbnRGYW1pbHkiOiJcInRyZWJ1Y2hldCBtc1wiLCB2ZXJkYW5hLCBhcmlhbCIsImZvbnRTaXplIjoiMTZweCIsImxhYmVsQmFja2dyb3VuZCI6IiNlOGU4ZTgiLCJub2RlQmtnIjoiI0VDRUNGRiIsIm5vZGVCb3JkZXIiOiIjOTM3MERCIiwiY2x1c3RlckJrZyI6IiNmZmZmZGUiLCJjbHVzdGVyQm9yZGVyIjoiI2FhYWEzMyIsImRlZmF1bHRMaW5rQ29sb3IiOiIjMzMzMzMzIiwidGl0bGVDb2xvciI6IiMzMzMiLCJlZGdlTGFiZWxCYWNrZ3JvdW5kIjoiI2U4ZThlOCIsImFjdG9yQm9yZGVyIjoiaHNsKDI1OS42MjYxNjgyMjQzLCA1OS43NzY1MzYzMTI4JSwgODcuOTAxOTYwNzg0MyUpIiwiYWN0b3JCa2ciOiIjRUNFQ0ZGIiwiYWN0b3JUZXh0Q29sb3IiOiJibGFjayIsImFjdG9yTGluZUNvbG9yIjoiZ3JleSIsInNpZ25hbENvbG9yIjoiIzMzMyIsInNpZ25hbFRleHRDb2xvciI6IiMzMzMiLCJsYWJlbEJveEJrZ0NvbG9yIjoiI0VDRUNGRiIsImxhYmVsQm94Qm9yZGVyQ29sb3IiOiJoc2woMjU5LjYyNjE2ODIyNDMsIDU5Ljc3NjUzNjMxMjglLCA4Ny45MDE5NjA3ODQzJSkiLCJsYWJlbFRleHRDb2xvciI6ImJsYWNrIiwibG9vcFRleHRDb2xvciI6ImJsYWNrIiwibm90ZUJvcmRlckNvbG9yIjoiI2FhYWEzMyIsIm5vdGVCa2dDb2xvciI6IiNmZmY1YWQiLCJub3RlVGV4dENvbG9yIjoiYmxhY2siLCJhY3RpdmF0aW9uQm9yZGVyQ29sb3IiOiIjNjY2IiwiYWN0aXZhdGlvbkJrZ0NvbG9yIjoiI2Y0ZjRmNCIsInNlcXVlbmNlTnVtYmVyQ29sb3IiOiJ3aGl0ZSIsInNlY3Rpb25Ca2dDb2xvciI6InJnYmEoMTAyLCAxMDIsIDI1NSwgMC40OSkiLCJhbHRTZWN0aW9uQmtnQ29sb3IiOiJ3aGl0ZSIsInNlY3Rpb25Ca2dDb2xvcjIiOiIjZmZmNDAwIiwidGFza0JvcmRlckNvbG9yIjoiIzUzNGZiYyIsInRhc2tCa2dDb2xvciI6IiM4YTkwZGQiLCJ0YXNrVGV4dExpZ2h0Q29sb3IiOiJ3aGl0ZSIsInRhc2tUZXh0Q29sb3IiOiJ3aGl0ZSIsInRhc2tUZXh0RGFya0NvbG9yIjoiYmxhY2siLCJ0YXNrVGV4dE91dHNpZGVDb2xvciI6ImJsYWNrIiwidGFza1RleHRDbGlja2FibGVDb2xvciI6IiMwMDMxNjMiLCJhY3RpdmVUYXNrQm9yZGVyQ29sb3IiOiIjNTM0ZmJjIiwiYWN0aXZlVGFza0JrZ0NvbG9yIjoiI2JmYzdmZiIsImdyaWRDb2xvciI6ImxpZ2h0Z3JleSIsImRvbmVUYXNrQmtnQ29sb3IiOiJsaWdodGdyZXkiLCJkb25lVGFza0JvcmRlckNvbG9yIjoiZ3JleSIsImNyaXRCb3JkZXJDb2xvciI6IiNmZjg4ODgiLCJjcml0QmtnQ29sb3IiOiJyZWQiLCJ0b2RheUxpbmVDb2xvciI6InJlZCIsImxhYmVsQ29sb3IiOiJibGFjayIsImVycm9yQmtnQ29sb3IiOiIjNTUyMjIyIiwiZXJyb3JUZXh0Q29sb3IiOiIjNTUyMjIyIiwiY2xhc3NUZXh0IjoiIzEzMTMwMCIsImZpbGxUeXBlMCI6IiNFQ0VDRkYiLCJmaWxsVHlwZTEiOiIjZmZmZmRlIiwiZmlsbFR5cGUyIjoiaHNsKDMwNCwgMTAwJSwgOTYuMjc0NTA5ODAzOSUpIiwiZmlsbFR5cGUzIjoiaHNsKDEyNCwgMTAwJSwgOTMuNTI5NDExNzY0NyUpIiwiZmlsbFR5cGU0IjoiaHNsKDE3NiwgMTAwJSwgOTYuMjc0NTA5ODAzOSUpIiwiZmlsbFR5cGU1IjoiaHNsKC00LCAxMDAlLCA5My41Mjk0MTE3NjQ3JSkiLCJmaWxsVHlwZTYiOiJoc2woOCwgMTAwJSwgOTYuMjc0NTA5ODAzOSUpIiwiZmlsbFR5cGU3IjoiaHNsKDE4OCwgMTAwJSwgOTMuNTI5NDExNzY0NyUpIn19LCJ1cGRhdGVFZGl0b3IiOmZhbHNlfQ)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoiZ3JhcGggVERcbiAgU3RhcnRcbiAgQVtOb2RlIC0gQ3JlYXRlT3JkZXJdIC0tPiBCKE5vZGUgLSBDcmF0ZVBheW1lbnQpXG4gIEItLSBlcnJvciBoYXBwZW5kIC0tPkMoTm9kZSAtIENhbmNlbGxCb29raW5nKVxuICBTdGFydC0tPkFcbiAgQi0tPkZpbmlzaCIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0IiwidGhlbWVWYXJpYWJsZXMiOnsiYmFja2dyb3VuZCI6IndoaXRlIiwicHJpbWFyeUNvbG9yIjoiI0VDRUNGRiIsInNlY29uZGFyeUNvbG9yIjoiI2ZmZmZkZSIsInRlcnRpYXJ5Q29sb3IiOiJoc2woODAsIDEwMCUsIDk2LjI3NDUwOTgwMzklKSIsInByaW1hcnlCb3JkZXJDb2xvciI6ImhzbCgyNDAsIDYwJSwgODYuMjc0NTA5ODAzOSUpIiwic2Vjb25kYXJ5Qm9yZGVyQ29sb3IiOiJoc2woNjAsIDYwJSwgODMuNTI5NDExNzY0NyUpIiwidGVydGlhcnlCb3JkZXJDb2xvciI6ImhzbCg4MCwgNjAlLCA4Ni4yNzQ1MDk4MDM5JSkiLCJwcmltYXJ5VGV4dENvbG9yIjoiIzEzMTMwMCIsInNlY29uZGFyeVRleHRDb2xvciI6IiMwMDAwMjEiLCJ0ZXJ0aWFyeVRleHRDb2xvciI6InJnYig5LjUwMDAwMDAwMDEsIDkuNTAwMDAwMDAwMSwgOS41MDAwMDAwMDAxKSIsImxpbmVDb2xvciI6IiMzMzMzMzMiLCJ0ZXh0Q29sb3IiOiIjMzMzIiwibWFpbkJrZyI6IiNFQ0VDRkYiLCJzZWNvbmRCa2ciOiIjZmZmZmRlIiwiYm9yZGVyMSI6IiM5MzcwREIiLCJib3JkZXIyIjoiI2FhYWEzMyIsImFycm93aGVhZENvbG9yIjoiIzMzMzMzMyIsImZvbnRGYW1pbHkiOiJcInRyZWJ1Y2hldCBtc1wiLCB2ZXJkYW5hLCBhcmlhbCIsImZvbnRTaXplIjoiMTZweCIsImxhYmVsQmFja2dyb3VuZCI6IiNlOGU4ZTgiLCJub2RlQmtnIjoiI0VDRUNGRiIsIm5vZGVCb3JkZXIiOiIjOTM3MERCIiwiY2x1c3RlckJrZyI6IiNmZmZmZGUiLCJjbHVzdGVyQm9yZGVyIjoiI2FhYWEzMyIsImRlZmF1bHRMaW5rQ29sb3IiOiIjMzMzMzMzIiwidGl0bGVDb2xvciI6IiMzMzMiLCJlZGdlTGFiZWxCYWNrZ3JvdW5kIjoiI2U4ZThlOCIsImFjdG9yQm9yZGVyIjoiaHNsKDI1OS42MjYxNjgyMjQzLCA1OS43NzY1MzYzMTI4JSwgODcuOTAxOTYwNzg0MyUpIiwiYWN0b3JCa2ciOiIjRUNFQ0ZGIiwiYWN0b3JUZXh0Q29sb3IiOiJibGFjayIsImFjdG9yTGluZUNvbG9yIjoiZ3JleSIsInNpZ25hbENvbG9yIjoiIzMzMyIsInNpZ25hbFRleHRDb2xvciI6IiMzMzMiLCJsYWJlbEJveEJrZ0NvbG9yIjoiI0VDRUNGRiIsImxhYmVsQm94Qm9yZGVyQ29sb3IiOiJoc2woMjU5LjYyNjE2ODIyNDMsIDU5Ljc3NjUzNjMxMjglLCA4Ny45MDE5NjA3ODQzJSkiLCJsYWJlbFRleHRDb2xvciI6ImJsYWNrIiwibG9vcFRleHRDb2xvciI6ImJsYWNrIiwibm90ZUJvcmRlckNvbG9yIjoiI2FhYWEzMyIsIm5vdGVCa2dDb2xvciI6IiNmZmY1YWQiLCJub3RlVGV4dENvbG9yIjoiYmxhY2siLCJhY3RpdmF0aW9uQm9yZGVyQ29sb3IiOiIjNjY2IiwiYWN0aXZhdGlvbkJrZ0NvbG9yIjoiI2Y0ZjRmNCIsInNlcXVlbmNlTnVtYmVyQ29sb3IiOiJ3aGl0ZSIsInNlY3Rpb25Ca2dDb2xvciI6InJnYmEoMTAyLCAxMDIsIDI1NSwgMC40OSkiLCJhbHRTZWN0aW9uQmtnQ29sb3IiOiJ3aGl0ZSIsInNlY3Rpb25Ca2dDb2xvcjIiOiIjZmZmNDAwIiwidGFza0JvcmRlckNvbG9yIjoiIzUzNGZiYyIsInRhc2tCa2dDb2xvciI6IiM4YTkwZGQiLCJ0YXNrVGV4dExpZ2h0Q29sb3IiOiJ3aGl0ZSIsInRhc2tUZXh0Q29sb3IiOiJ3aGl0ZSIsInRhc2tUZXh0RGFya0NvbG9yIjoiYmxhY2siLCJ0YXNrVGV4dE91dHNpZGVDb2xvciI6ImJsYWNrIiwidGFza1RleHRDbGlja2FibGVDb2xvciI6IiMwMDMxNjMiLCJhY3RpdmVUYXNrQm9yZGVyQ29sb3IiOiIjNTM0ZmJjIiwiYWN0aXZlVGFza0JrZ0NvbG9yIjoiI2JmYzdmZiIsImdyaWRDb2xvciI6ImxpZ2h0Z3JleSIsImRvbmVUYXNrQmtnQ29sb3IiOiJsaWdodGdyZXkiLCJkb25lVGFza0JvcmRlckNvbG9yIjoiZ3JleSIsImNyaXRCb3JkZXJDb2xvciI6IiNmZjg4ODgiLCJjcml0QmtnQ29sb3IiOiJyZWQiLCJ0b2RheUxpbmVDb2xvciI6InJlZCIsImxhYmVsQ29sb3IiOiJibGFjayIsImVycm9yQmtnQ29sb3IiOiIjNTUyMjIyIiwiZXJyb3JUZXh0Q29sb3IiOiIjNTUyMjIyIiwiY2xhc3NUZXh0IjoiIzEzMTMwMCIsImZpbGxUeXBlMCI6IiNFQ0VDRkYiLCJmaWxsVHlwZTEiOiIjZmZmZmRlIiwiZmlsbFR5cGUyIjoiaHNsKDMwNCwgMTAwJSwgOTYuMjc0NTA5ODAzOSUpIiwiZmlsbFR5cGUzIjoiaHNsKDEyNCwgMTAwJSwgOTMuNTI5NDExNzY0NyUpIiwiZmlsbFR5cGU0IjoiaHNsKDE3NiwgMTAwJSwgOTYuMjc0NTA5ODAzOSUpIiwiZmlsbFR5cGU1IjoiaHNsKC00LCAxMDAlLCA5My41Mjk0MTE3NjQ3JSkiLCJmaWxsVHlwZTYiOiJoc2woOCwgMTAwJSwgOTYuMjc0NTA5ODAzOSUpIiwiZmlsbFR5cGU3IjoiaHNsKDE4OCwgMTAwJSwgOTMuNTI5NDExNzY0NyUpIn19LCJ1cGRhdGVFZGl0b3IiOmZhbHNlfQ)
* A Flow is combined by multiple nodes, node an be reuse.
* Nodes will be executed in order.
* Nodes can delivery params.


### Sync Flow
![seq-sync](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgU2FnYVN5bmMKQ2xpZW50LT5FbmRwb2ludDog6KiC5Zau6LOH5paZIChodHRwKQoAFggtPk9yY2hlc3RyYXRvcgAcEGdSUEMpCgAWDC0-T3JkZQAQGWRlABgK5bu656uLAHcGABEKAFsS57eo6JmfAFcWUGF5bWUAgToKAB4NIAoAFwcAHAsAbgbku5jmrL7llq4AGAoAgVgOABsG6LOH6KiKAIFQFgCCNAoAGBQAgjQKAIJoBgBCEACCXgYKCgoK&s=roundgreen)

* Using Context to combine all data from each service

### Rollback
![seq-rollback](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgUm9sbGJhY2sKQ2xpZW50LT5FbmRwb2ludDog6KiC5Zau6LOH5paZIChodHRwKQoAFggtPk9yY2hlc3RyYXRvcgAcEGdSUEMpCgAWDC0-T3JkZQAQGWRlABgK5bu656uLAHcGABEKAFsS57eo6JmfAFcWUGF5bWUAgToKAB4NIAoAFwcAHAsAbgbmlLbmrL7llq7lpLHmlZcAHgoAgV4OZXJyb3IAgU8WTVE6IHRyYWNlSUQgKHB1Ymxpc2gpCk1RAIIfEAAbCWNvbnN1bWUAghYXAEIJAIITFY-W5raIAIIMHE9LAIJ1FgCDWQoAgT4NAINSCgCEBgYAgWEJAIN1BgoK&s=roundgreen)
* If something wrong happend, just call rollback node.
* The important thing in Message is traceID (or x-request-id) of requst chain.
* If rollback trigger, Ochestrator invoke cancel func to each service, service will cancel data by traceID


### Async Flow
![sql-async](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgU2FnYUFzeW5jCgpDbGllbnQtPkVuZHBvaW50OiDoqILllq7os4fmlpkgKGh0dHApCgAWCC0-T3JjaGVzdHJhdG9yABwQZ1JQQykKCgAXDC0-TVEAQxBwdWJsaXNoKQAaDwB4Ck9LAEQIAG4KAIEiBgAWBgCBDgYKTVEAcR5jb25zdW1lAFgQT3JkZQCBJBcAFgUAGQnlu7rnq4sAggkGABMIAIB_CENhbGxiYWNrAIIZCAAxCQCCEgwAgT4KCgoKCgoK&s=roundgreen)
* Data is deliveried in JSON format via MQ. orchestrator(and it's replicas) push message to itself. then send to service via gRPC.
* Response directly when it just receive request. Meanwhile, start preccess in backgroud(MQ)
* If client needs data. send it via http callback.


### Throttling Flow
![seq-throttling](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgVGhyb3R0bGluZwpDbGllbnQtPkVuZHBvaW50OiDoqILllq7os4fmlpkgKGh0dHApCgAWCC0-T3JjaGVzdHJhdG9yABwQZ1JQQykKABYMACAQ55uj6IG95rWB56iL5piv5ZCm5a6M5oiQACYPTVEAdxBwdWJsaXNoKQpNUQBmHmNvbnN1bWUAdRJkZQCBFxlkZQAYCuW7uueriwCBfgYAEQoAgWISSUQAgUwk5L-u5pS5AIFqBueCuuW3sgCBWxUAgm0QAFMKAIJvCgCDIwYAcwwAgxUG&s=roundgreen)
* Using MQ to controll Queue and Concurrency counts, fix that upstream overloading while condurrency.
* It lock request, and do task in background. if task finished, response to client.
* Using topic + pod_id for dispatching to same Pod.

![flow-throttling](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBHYXRld2F5XG4gICAgTVFcbiAgICBHYXRld2F5IC0tIDEuIGdSUEMgTG9hZEJhbGFuY2luZyAtLT5wb2QxXG4gICAgR2F0ZXdheS0tIDEuIGdSUEMgTG9hZEJhbGFuY2luZyAtLT4gcG9kMlxuICAgIHBvZDEtLSAyLiBwdWIgdG9waWMucG9kMSAtLT5NUVxuICAgIHBvZDItLSAyLiBwdWIgdG9waWMucG9kMiAtLT5NUVxuICAgIE1RIC0tIDMuIHN1YiB0b3BpYy5wb2QxIC0tPiBwb2QxXG4gICAgTVEgLS0gMy4gc3ViIHRvcGljLnBvZDIgLS0-IHBvZDJcbiAgICBcbiAgICBzdWJncmFwaCBvcmNoZXN0cmF0b3JcbiAgICBwb2QxXG4gICAgcG9kMlxuICAgIGVuZFxuIiwibWVybWFpZCI6eyJ0aGVtZSI6ImRlZmF1bHQiLCJ0aGVtZVZhcmlhYmxlcyI6eyJiYWNrZ3JvdW5kIjoid2hpdGUiLCJwcmltYXJ5Q29sb3IiOiIjRUNFQ0ZGIiwic2Vjb25kYXJ5Q29sb3IiOiIjZmZmZmRlIiwidGVydGlhcnlDb2xvciI6ImhzbCg4MCwgMTAwJSwgOTYuMjc0NTA5ODAzOSUpIiwicHJpbWFyeUJvcmRlckNvbG9yIjoiaHNsKDI0MCwgNjAlLCA4Ni4yNzQ1MDk4MDM5JSkiLCJzZWNvbmRhcnlCb3JkZXJDb2xvciI6ImhzbCg2MCwgNjAlLCA4My41Mjk0MTE3NjQ3JSkiLCJ0ZXJ0aWFyeUJvcmRlckNvbG9yIjoiaHNsKDgwLCA2MCUsIDg2LjI3NDUwOTgwMzklKSIsInByaW1hcnlUZXh0Q29sb3IiOiIjMTMxMzAwIiwic2Vjb25kYXJ5VGV4dENvbG9yIjoiIzAwMDAyMSIsInRlcnRpYXJ5VGV4dENvbG9yIjoicmdiKDkuNTAwMDAwMDAwMSwgOS41MDAwMDAwMDAxLCA5LjUwMDAwMDAwMDEpIiwibGluZUNvbG9yIjoiIzMzMzMzMyIsInRleHRDb2xvciI6IiMzMzMiLCJtYWluQmtnIjoiI0VDRUNGRiIsInNlY29uZEJrZyI6IiNmZmZmZGUiLCJib3JkZXIxIjoiIzkzNzBEQiIsImJvcmRlcjIiOiIjYWFhYTMzIiwiYXJyb3doZWFkQ29sb3IiOiIjMzMzMzMzIiwiZm9udEZhbWlseSI6IlwidHJlYnVjaGV0IG1zXCIsIHZlcmRhbmEsIGFyaWFsIiwiZm9udFNpemUiOiIxNnB4IiwibGFiZWxCYWNrZ3JvdW5kIjoiI2U4ZThlOCIsIm5vZGVCa2ciOiIjRUNFQ0ZGIiwibm9kZUJvcmRlciI6IiM5MzcwREIiLCJjbHVzdGVyQmtnIjoiI2ZmZmZkZSIsImNsdXN0ZXJCb3JkZXIiOiIjYWFhYTMzIiwiZGVmYXVsdExpbmtDb2xvciI6IiMzMzMzMzMiLCJ0aXRsZUNvbG9yIjoiIzMzMyIsImVkZ2VMYWJlbEJhY2tncm91bmQiOiIjZThlOGU4IiwiYWN0b3JCb3JkZXIiOiJoc2woMjU5LjYyNjE2ODIyNDMsIDU5Ljc3NjUzNjMxMjglLCA4Ny45MDE5NjA3ODQzJSkiLCJhY3RvckJrZyI6IiNFQ0VDRkYiLCJhY3RvclRleHRDb2xvciI6ImJsYWNrIiwiYWN0b3JMaW5lQ29sb3IiOiJncmV5Iiwic2lnbmFsQ29sb3IiOiIjMzMzIiwic2lnbmFsVGV4dENvbG9yIjoiIzMzMyIsImxhYmVsQm94QmtnQ29sb3IiOiIjRUNFQ0ZGIiwibGFiZWxCb3hCb3JkZXJDb2xvciI6ImhzbCgyNTkuNjI2MTY4MjI0MywgNTkuNzc2NTM2MzEyOCUsIDg3LjkwMTk2MDc4NDMlKSIsImxhYmVsVGV4dENvbG9yIjoiYmxhY2siLCJsb29wVGV4dENvbG9yIjoiYmxhY2siLCJub3RlQm9yZGVyQ29sb3IiOiIjYWFhYTMzIiwibm90ZUJrZ0NvbG9yIjoiI2ZmZjVhZCIsIm5vdGVUZXh0Q29sb3IiOiJibGFjayIsImFjdGl2YXRpb25Cb3JkZXJDb2xvciI6IiM2NjYiLCJhY3RpdmF0aW9uQmtnQ29sb3IiOiIjZjRmNGY0Iiwic2VxdWVuY2VOdW1iZXJDb2xvciI6IndoaXRlIiwic2VjdGlvbkJrZ0NvbG9yIjoicmdiYSgxMDIsIDEwMiwgMjU1LCAwLjQ5KSIsImFsdFNlY3Rpb25Ca2dDb2xvciI6IndoaXRlIiwic2VjdGlvbkJrZ0NvbG9yMiI6IiNmZmY0MDAiLCJ0YXNrQm9yZGVyQ29sb3IiOiIjNTM0ZmJjIiwidGFza0JrZ0NvbG9yIjoiIzhhOTBkZCIsInRhc2tUZXh0TGlnaHRDb2xvciI6IndoaXRlIiwidGFza1RleHRDb2xvciI6IndoaXRlIiwidGFza1RleHREYXJrQ29sb3IiOiJibGFjayIsInRhc2tUZXh0T3V0c2lkZUNvbG9yIjoiYmxhY2siLCJ0YXNrVGV4dENsaWNrYWJsZUNvbG9yIjoiIzAwMzE2MyIsImFjdGl2ZVRhc2tCb3JkZXJDb2xvciI6IiM1MzRmYmMiLCJhY3RpdmVUYXNrQmtnQ29sb3IiOiIjYmZjN2ZmIiwiZ3JpZENvbG9yIjoibGlnaHRncmV5IiwiZG9uZVRhc2tCa2dDb2xvciI6ImxpZ2h0Z3JleSIsImRvbmVUYXNrQm9yZGVyQ29sb3IiOiJncmV5IiwiY3JpdEJvcmRlckNvbG9yIjoiI2ZmODg4OCIsImNyaXRCa2dDb2xvciI6InJlZCIsInRvZGF5TGluZUNvbG9yIjoicmVkIiwibGFiZWxDb2xvciI6ImJsYWNrIiwiZXJyb3JCa2dDb2xvciI6IiM1NTIyMjIiLCJlcnJvclRleHRDb2xvciI6IiM1NTIyMjIiLCJjbGFzc1RleHQiOiIjMTMxMzAwIiwiZmlsbFR5cGUwIjoiI0VDRUNGRiIsImZpbGxUeXBlMSI6IiNmZmZmZGUiLCJmaWxsVHlwZTIiOiJoc2woMzA0LCAxMDAlLCA5Ni4yNzQ1MDk4MDM5JSkiLCJmaWxsVHlwZTMiOiJoc2woMTI0LCAxMDAlLCA5My41Mjk0MTE3NjQ3JSkiLCJmaWxsVHlwZTQiOiJoc2woMTc2LCAxMDAlLCA5Ni4yNzQ1MDk4MDM5JSkiLCJmaWxsVHlwZTUiOiJoc2woLTQsIDEwMCUsIDkzLjUyOTQxMTc2NDclKSIsImZpbGxUeXBlNiI6ImhzbCg4LCAxMDAlLCA5Ni4yNzQ1MDk4MDM5JSkiLCJmaWxsVHlwZTciOiJoc2woMTg4LCAxMDAlLCA5My41Mjk0MTE3NjQ3JSkifX19)

#### Performance
* Stress Situation
  * All service has 3 replicas
  * Memory limit is set in k8s, it will cause OOM killed while concurrency.
  * rps 500 / hit 30 second
  * In throttling mode, each pod open 100 Consumer


* Sync Flow：
  ![press-sync-rate500-dur30](doc/press-sync-rate500-dur30.png)
  ````
  Requests      [total, rate, throughput]         15000, 500.04, 250.18
  Duration      [total, attack, wait]             36.878s, 29.998s, 6.88s
  Latencies     [min, mean, 50, 90, 95, 99, max]  20.014ms, 6.019s, 4.786s, 13.826s, 17.791s, 24.55s, 30.274s
  Bytes In      [total, mean]                     793436, 52.90
  Bytes Out     [total, mean]                     0, 0.00
  Success       [ratio]                           61.51%
  Status Codes  [code:count]                      200:9226  500:5774
  Error Set:
  500 Internal Server Error
  ````
  * Upstream were killed
* Throttling mode
  ![press-throttling-rate500-dur30](doc/press-throttling-rate500-dur30.png)
  ````
  Requests      [total, rate, throughput]         10000, 333.36, 212.27
  Duration      [total, attack, wait]             47.11s, 29.998s, 17.113s
  Latencies     [min, mean, 50, 90, 95, 99, max]  256.525ms, 10.65s, 11.651s, 16.731s, 17.258s, 18.084s, 19.187s
  Bytes In      [total, mean]                     880000, 88.00
  Bytes Out     [total, mean]                     0, 0.00
  Success       [ratio]                           100.00%
  Status Codes  [code:count]                      200:10000
  ````

  * All requests were Consumed steady


## Task

### 1. Implements a Node
* Register node under ./node
#### Sync Flow Node
Get Param: Get gRPC Request/Response from context, and convert to protobuf object

```go
func CreateOrderSync() orchestrator.SyncNode {
	return func(requestID string, ctx *ctx.Context) error {
        //  gRPC Request from context
        if req, isExist := ctx.Get("bookingSyncPbReq") ; isExist {
            if request, ok := req.(*pb.BookingRequest); ok {  
```
Push handled data into context, keep delivering.

```go
    resp := &pb.BookingSyncResponse{
        RequestID:            requestID,
        OrderID:              mockOrderID,
        PaymentID:            0,
    }
    ctx.Set("bookingSyncPbResp", resp)
```
#### Async Transaction flow
Declare a struct for delivery data in flow, it needs to extend ``` *orchestrator.AsyncFlowContext ```
```go
type BookingMsgDTO struct {
	//**** Request ****//
	FaultInject bool
	ProductID   int64 `json:"product_id"`
	//**** Delivery ****//
	OrderID int64 `json:"order_id"`
	PaymentID int64 `json:"payment_id"` 
	//**** Context ****//
	*orchestrator.AsyncFlowContext
}
```
AsyncNode will inject following param
* topic topic ID
* data marshaled json data
* next() Going to next node
* rollback() cancel flow while error happened
```go
func CreateOrderAsync() orchestrator.AsyncNode {
	return func(topic orchestrator.Topic, data []byte, next orchestrator.Next, rollback orchestrator.Rollback) {
```
Receive param: convert JSON from context

```go
    d := &BookingMsgDTO{}
    if err := json.Unmarshal(data, d); err != nil {
```
Put Handled data into context, invoke next() for going to next node

```go
d.PaymentID = mockPaymentID
next(d)
```
Invoke rollback() while error happened
````go
    if err := nil {
        rollback(err, d)
        return
    }
````


* Example reference [node](./node/booking.go)

### 2. Register transaction flow
* Register a flow under ./facade
* Define facades，for invoking flow easily.
```go
const (
	AsyncBooking orchestrator.Facade = "AsyncBooking"
	SyncBooking  = "SyncBooking"
	ThrottlingBooking  = "ThrottlingBooking"
)
```

Create Pairs ：One Topic mapping one node, Sync Flow has only one pair : rollback
```go
	createPaymentPair := orchestrator.TopicNodePair{
		Topic:     topic.CreatePayment,
		AsyncNode: node.CreatePaymentAsync(),
	}
```
Create Flow
```go
orchestrator.NewThrottlingFlow(topic.CancelThrottlingBooking)
```

Register Pairs in order
```go
	flow.Use(createOrderPair)
	flow.Use(createPaymentPair)
```

Watch Topic: start up all consumers
```go
	// Start to watch Topic of flow
	flow.Consume()
	// Start to watch rollback topic
	flow.ConsumeRollback(rollbackPair)
```

Register in singlton, it will be invoked in concurrency
```go
orchestrator.GetInstance().SetThrottlingFlows(ThrottlingBooking, flow)
```


* Example refer [facade](./facade/booking.go)

### 3. Invoke transaction flow
In concurrency(or a request), we can get registered flow from singlton. just use run() to exectue flow.

````go
flow := orchestrator.GetInstance().GetThrottlingFlows(facade.ThrottlingBooking)
flow.Run(requestID, reqMsg)
````

* Example refer [handler](./handler/booking.go)