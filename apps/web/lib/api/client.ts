const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "/api";

export async function request<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    cache: "no-store",
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`API error ${res.status}: ${text}`);
  }

  return res.json();
}
