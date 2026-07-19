import type {
  Worker,
  Metric,
  TelemetryEvent
  }
  from "./types";
  
  
  export const workers:Worker[]=[
  
  {
  id:"worker-01",
  status:"ACTIVE",
  cpu:42,
  memory:"2GB",
  checkpoint:"v42",
  role:"PRIMARY"
  }
  
  ];
  
  
  export const metrics:Metric[]=[
  
  {
  name:"Latency",
  value:"2.3ms"
  }
  
  ];
  
  
  export const events:TelemetryEvent[]=[
  
  {
  level:"INFO",
  message:"checkpoint created",
  timestamp:"now"
  }
  
  ];