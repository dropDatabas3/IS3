import { http, HttpResponse } from 'msw';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

export const handlers = [
  http.post(`${API_BASE}/auth/login`, async ({ request }) => {
    const body = await request.json().catch(() => ({} as any));
    if (!body?.Email || !body?.Password) {
      return HttpResponse.json({ message: 'missing fields' }, { status: 400 });
    }
    return HttpResponse.json(
      {
        token: 'test-token',
        user: {
          id: 'u1',
          email: body.Email,
          username: 'Test User',
          avatar: 'https://example.com/avatar.png',
          role: 'admin',
          createdAt: new Date().toISOString() as any,
          updatedAt: new Date().toISOString() as any,
        },
      },
      { status: 200 }
    );
  }),

  http.post(`${API_BASE}/users/register`, async ({ request }) => {
    const body = await request.json().catch(() => ({} as any));
    if (!body?.Username || !body?.Email || !body?.Password) {
      return HttpResponse.json({ message: 'missing fields' }, { status: 400 });
    }
    return HttpResponse.json(
      {
        token: 'test-token',
        user: {
          id: 'u2',
          email: body.Email,
          username: body.Username,
          avatar: 'https://example.com/avatar2.png',
          role: 'user',
          createdAt: new Date().toISOString() as any,
          updatedAt: new Date().toISOString() as any,
        },
      },
      { status: 201 }
    );
  }),

  http.post(`${API_BASE}/auth/refresh-token`, () =>
    HttpResponse.json(
      {
        token: 'refreshed-token',
        user: {
          id: 'u1',
          email: 'test@example.com',
          username: 'Test User',
          avatar: 'https://example.com/avatar.png',
          role: 'admin',
          createdAt: new Date().toISOString() as any,
          updatedAt: new Date().toISOString() as any,
        },
      },
      { status: 200 }
    )
  ),
];
