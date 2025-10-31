import { api, authHeader } from "./http";

export async function register(email: string, password: string, confirm: string) {
  const r = await api("/auth/register", {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify({ email, password, confirm_password: confirm }),
  });
  if (!r.ok) throw new Error(await r.text());
  return r.json();
}

export async function login(identifier: string, password: string) {
  const r = await api("/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" }, // login chưa cần token
    body: JSON.stringify({ identifier, password }),
  });
  if (!r.ok) throw new Error(await r.text());
  const data = await r.json();
  localStorage.setItem("token", data.access_token);
  return data.user as { id: number; username: string; email: string; role: "user" | "admin" };
}

export async function me() {
  const r = await api("/auth/me", { headers: authHeader() });
  if (!r.ok) throw new Error(await r.text());
  return r.json();
}
