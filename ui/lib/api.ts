const API_URL =
  process.env.NEXT_PUBLIC_API_URL ??
  "http://localhost:8080";

type HealthResponse = {
  status: string;
  service: string;
};

function joinUrl(base: string, endpoint: string) {
  if (!base) {
    return endpoint;
  }

  return `${base.replace(/\/$/, "")}${endpoint}`;
}

async function request<T>(endpoint: string, init?: RequestInit): Promise<T> {
  const response = await fetch(joinUrl(API_URL, endpoint), {
    cache: "no-store",
    ...init,
  });

  if (!response.ok) {
    const errorBody = await response.text().catch(() => "");

    throw new Error(
      errorBody ? `API Error ${response.status}: ${errorBody}` : `API Error ${response.status}`,
    );
  }

  const contentType = response.headers.get("content-type") ?? "";

  if (contentType.includes("application/json")) {
    return response.json() as Promise<T>;
  }

  return response.text() as Promise<T>;
}

export const api = {
  health() {
    return request<HealthResponse>("/health");
  },

  metrics() {
    return request<Record<string, number>>("/metrics");
  },

  checkpoint(payload: { workflow_id: string; step: number; payload: string }) {
    return request<{ status: string }>("/checkpoint", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });
  },

  recover(payload: { workflow_id: string }) {
    return request<{ status: string }>("/recover", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });
  },
};
