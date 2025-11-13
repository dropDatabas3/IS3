import { http, HttpResponse, delay } from 'msw';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

// Simple in-memory store for tests
let ratings = [
  { course_id: 'c1', user_id: 'u1', rating: 5 },
  { course_id: 'c1', user_id: 'u2', rating: 4 },
];

export const handlers = [
  // Load all ratings
  http.get(`${API_BASE}/rating`, async () => {
    await delay(20);
    return HttpResponse.json(ratings, { status: 200 });
  }),

  // Create rating
  http.post(`${API_BASE}/rating`, async ({ request }) => {
    const body = await request.json().catch(() => ({} as any));
    if (!body?.course_id || !body?.user_id || typeof body?.rating !== 'number') {
      return HttpResponse.json({ message: 'invalid rating' }, { status: 400 });
    }
    ratings.push({ course_id: body.course_id, user_id: body.user_id, rating: body.rating });
    return HttpResponse.json({ ok: true }, { status: 201 });
  }),

  // Update rating
  http.put(`${API_BASE}/rating`, async ({ request }) => {
    const body = await request.json().catch(() => ({} as any));
    const idx = ratings.findIndex(r => r.course_id === body?.course_id && r.user_id === body?.user_id);
    if (idx === -1) {
      return HttpResponse.json({ message: 'not found' }, { status: 404 });
    }
    ratings[idx] = { course_id: body.course_id, user_id: body.user_id, rating: body.rating };
    return HttpResponse.json({ ok: true }, { status: 200 });
  }),
];
