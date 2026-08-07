import type {
  Worker,
  Checkpoint,
  TelemetryEvent,
  Metric,
} from "./types";


const API_URL =
  process.env.NEXT_PUBLIC_API_URL ??
  "http://localhost:8080";


async function request<T>(
  endpoint:string
):Promise<T>{

  const response =
    await fetch(
      `${API_URL}${endpoint}`,
      {
        cache:"no-store",
      }
    );


  if(!response.ok){
    throw new Error(
      `API Error ${response.status}`
    );
  }


  return response.json();
}



export const api = {

  workers(){
    return request<Worker[]>(
      "/workers"
    );
  },


  checkpoints(){
    return request<Checkpoint[]>(
      "/checkpoints"
    );
  },


  metrics(){
    return request<Metric[]>(
      "/metrics"
    );
  },


  telemetry(){
    return request<TelemetryEvent[]>(
      "/telemetry"
    );
  }

};