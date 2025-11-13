import { cookies } from 'next/headers';
import type { LoginDto, RegisterDto } from '@/types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';
export interface ApiClient {
  get<T = unknown>(url: string, init?: RequestInit): Promise<T>;
  post<T = unknown>(url: string, body?: any, init?: RequestInit): Promise<T>;
  put<T = unknown>(url: string, body?: any, init?: RequestInit): Promise<T>;
  delete<T = unknown>(url: string, init?: RequestInit): Promise<T>;
}

function buildUrl(path: string) {
  return path.startsWith('http') ? path : `${API_URL}${path}`;
}

async function request<T>(url: string, init: RequestInit): Promise<T> {
  const res = await fetch(url, init);
  if (!res.ok) {
    const maybeJson = await res.json().catch(() => ({}));
    const message = (maybeJson && maybeJson.message) || `Request failed: ${res.status}`;
    throw new Error(message);
  }
  return res.json() as Promise<T>;
}

export function createApiClient(): ApiClient {
  return {
    get: (path, init) => request(buildUrl(path), { method: 'GET', ...(init || {}) }),
    post: (path, body, init) => request(buildUrl(path), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...(init?.headers || {}) },
      body: body !== undefined ? JSON.stringify(body) : undefined,
      ...(init || {})
    }),
    put: (path, body, init) => request(buildUrl(path), {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...(init?.headers || {}) },
      body: body !== undefined ? JSON.stringify(body) : undefined,
      ...(init || {})
    }),
    delete: (path, init) => request(buildUrl(path), { method: 'DELETE', ...(init || {}) })
  };
}

// Backwards compatible helpers (used by existing forms)
export async function loginUser(data: LoginDto) {
  return createApiClient().post('/auth/login', data);
}
export async function registerUser(data: RegisterDto) {
  return createApiClient().post('/users/register', data);
}
export function apiUrl() {
  const internal = process.env.INTERNAL_API || 'http://backend:8000'
  // Prefer runtime-injected config in the browser if available
  // @ts-ignore
  const runtimeApi = (typeof window !== 'undefined' && (window as any).__RUNTIME_CONFIG__?.API_URL) || ''
  const publicEnv = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000'
  const publicUrl = runtimeApi || publicEnv

  // If running on server (SSR) use internal service name; otherwise use public URL
  if (typeof window === 'undefined') return internal
  return publicUrl
}

export default apiUrl
