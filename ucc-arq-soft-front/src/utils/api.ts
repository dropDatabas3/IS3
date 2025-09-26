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
