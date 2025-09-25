export function apiUrl() {
  const internal = process.env.INTERNAL_API || 'http://backend:8000'
  const publicUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000'

  // If running on server (SSR) use internal service name; otherwise use public URL
  if (typeof window === 'undefined') return internal
  return publicUrl
}

export default apiUrl
