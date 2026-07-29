+-----------------------------------+
                  |   Incoming Ingress Traffic Path   |
                  +-----------------+-----------------+
                                    |
                        [ Failover Orchestrator ]
                        |  (Atomic Node Pointer) |
                        +---------+-----------+--+
                                  |           |
             (Active Proxy Path)  |           | (Standby Hot-Standby Path)
                                  v           v
                          +-------+---+   +---+-------+
                          |   Node A  |   |   Node B  |
                          | (Active)  |   | (Standby) |
                          +-----+-----+   +-----+-----+
                                |               |
  Passive Diagnostics Track     |               |  Passive Diagnostics Track
  (Latency, 5xx, Conntrack)     |               |  (Latency, 5xx, Conntrack)
                                v               v
                     +----------+---------------+----------+
                     | Lock-Free Sliding Telemetry Buffer  |
                     +------------------+------------------+
                                        |
                             Threshold Evaluator
                        (Immediate Failover Signal)