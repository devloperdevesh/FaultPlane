export interface Worker {

  id:string;
  
  status:
  "ACTIVE" |
  "FAILED" |
  "RECOVERING";
  
  cpu:number;
  
  memory:string;
  
  checkpoint:string;
  
  role:
  "PRIMARY" |
  "STANDBY";
  
  }
  
  
  
  export interface Workflow {
  
  id:string;
  
  name:string;
  
  status:
  "RUNNING"|
  "FAILED";
  
  }
  
  
  export interface Checkpoint {
  
  id:string;
  
  createdAt:string;
  
  size:string;
  
  }
  
  
  export interface Metric {
  
  name:string;
  
  value:string;
  
  }
  
  
  export interface TelemetryEvent {
  
  level:
  "INFO"|
  "WARN"|
  "ERROR";
  
  message:string;
  
  timestamp:string;
  
  }