import { api, authHeader } from "./http";
import type { User } from "../types";

export type Gender = "male" | "female" | "other";
export type Role = "user" | "admin";
export type Status = "active" | "inactive" | "banned" | string;

export interface BaseUserInput {
  username: string;
  email: string;
  full_name?: string;
  phone?: string;
  gender?: Gender | string;
  date_of_birth?: string;      // yyyy-mm-dd
  avatar_url?: string;
  street?: string;
  city?: string;
  state?: string;
  country?: string;
  postal_code?: string;
  role?: Role;                  // chỉ admin mới đổi
  status?: Status;
}

export interface CreateUserInput extends BaseUserInput {
  password: string;
}

export interface UpdateUserInput extends BaseUserInput {
  password?: string;
}

export async function getUsers(): Promise<User[]> {
  const r = await api("/admin/users", { headers: authHeader() });
  if (!r.ok) throw new Error(await r.text());
  return (await r.json()) as User[];
}

export async function getUser(id: number): Promise<User> {
  const r = await api(`/admin/users/${id}`, { headers: authHeader() });
  if (!r.ok) throw new Error(await r.text());
  return (await r.json()) as User;
}

export async function createUser(input: CreateUserInput): Promise<User> {
  const r = await api("/admin/users", {
    method: "POST",
    headers: authHeader(),
    body: JSON.stringify(input),
  });
  if (!r.ok) throw new Error(await r.text());
  return (await r.json()) as User;
}

export async function updateUser(id: number, input: UpdateUserInput): Promise<User> {
  const r = await api(`/admin/users/${id}`, {
    method: "PUT",
    headers: authHeader(),
    body: JSON.stringify(input),
  });
  if (!r.ok) throw new Error(await r.text());
  return (await r.json()) as User;
}

export async function deleteUser(id: number): Promise<void> {
  const r = await api(`/admin/users/${id}`, { method: "DELETE", headers: authHeader() });
  if (!r.ok) throw new Error(await r.text());
}
