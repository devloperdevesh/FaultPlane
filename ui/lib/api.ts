import type {
  Worker
  }
  from "./types";
  
  
  const API_URL =
  process.env.NEXT_PUBLIC_API_URL ??
  "http://localhost:8000";
  
  
  export async function getWorkers(){
  
  const res =
  await fetch(
  `${API_URL}/workers`
  );
  
  
  if(!res.ok){
  
  throw new Error(
  "Failed fetching workers"
  );
  
  }
  
  
  return res.json() as Promise<Worker[]>;
  
  }