const API_BASE =
  import.meta.env.VITE_API_BASE ?? "http://localhost:8080/api/v1";

/** Trả về headers Authorization nếu có token */
export function authHeader(): Record<string, string> {
  const token = localStorage.getItem("token");
  const h: Record<string, string> = {};
  if (token) h.Authorization = `Bearer ${token}`;
  return h;
}

/** Wrapper fetch: headers kiểu Record<string,string> cho gọn */
export async function api(
  path: string,
  init: Omit<RequestInit, "headers"> & { headers?: Record<string, string> } = {}
) {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(init.headers || {}),
  };

  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers,
    credentials: "include",
  });

  if (res.status === 401) {
    localStorage.removeItem("token");
  }
  return res;
}
