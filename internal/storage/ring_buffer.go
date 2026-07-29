+-------------------------------------------------------+
                    |               RingBuffer Array (Capacity: N)           |
                    | [Seq 10: Acked] [Seq 11: Unack] ... [Seq 19: Unack]   |
                    +-------------------------------------------------------+
                                        ^                         ^
                                        |                         |
                                   Tail (Oldest)             Head (Write Pointer)
                                        |                         |
               +------------------------+-------------------------+
               |  Unacknowledged Window Tracker (Tenant Indexing)  |
               |  Tenant A -> { Seq 11, Seq 13 }                   |
               |  Tenant B -> { Seq 12, Seq 14 }                   |
               +--------------------------------------------------+