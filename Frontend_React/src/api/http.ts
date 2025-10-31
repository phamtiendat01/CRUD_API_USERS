// Frontend_React/src/api/http.ts

// Lấy base URL từ env (Netlify đã đặt VITE_API_BASE_URL)
const RAW = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.trim();

export const API_BASE =
  RAW && RAW !== ""
    ? RAW.replace(/\/+$/, "") // bỏ dấu / ở cuối cho sạch
    : import.meta.env.DEV
      ? "http://localhost:8080/api/v1" // chỉ fallback khi DEV
      : (() => {
          throw new Error("VITE_API_BASE_URL is not defined for production build");
        })();

/** Authorization header nếu có token trong localStorage */
export function authHeader(): Record<string, string> {
  const token = localStorage.getItem("token");
  const h: Record<string, string> = {};
  if (token) h.Authorization = `Bearer ${token}`;
  return h;
}

/** Wrapper fetch gọn với headers dạng Record<string,string> */
export async function api(
  path: string,
  init: Omit<RequestInit, "headers"> & { headers?: Record<string, string> } = {}
) {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(init.headers || {}),
  };

  // Ghép URL an toàn (nếu path không có / đầu thì tự thêm)
  const url = `${API_BASE}${path.startsWith("/") ? "" : "/"}${path}`;

  const res = await fetch(url, {
    ...init,
    headers,
    credentials: "include",
  });

  if (res.status === 401) {
    localStorage.removeItem("token");
  }
  return res;
}
