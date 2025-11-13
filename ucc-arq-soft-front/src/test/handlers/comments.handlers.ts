import { http, HttpResponse, delay } from 'msw';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

let comments: Array<{ course_id: string; user_id: string; text: string }> = [];

export const handlers = [
  // Get comments by course
  http.get(`${API_BASE}/comment/:courseId`, async ({ params }) => {
    await delay(10);
    const courseId = params.courseId as string;
    const list = comments.filter(c => c.course_id === courseId);
    if (list.length === 0) return HttpResponse.json({ message: 'no comments' }, { status: 404 });
    return HttpResponse.json(list, { status: 200 });
  }),

  // Create comment
  http.post(`${API_BASE}/comment`, async ({ request }) => {
    const body = await request.json().catch(() => ({} as any));
    if (!body?.course_id || !body?.user_id || !body?.text) {
      return HttpResponse.json({ message: 'invalid comment' }, { status: 400 });
    }
    comments.push({ course_id: body.course_id, user_id: body.user_id, text: body.text });
    return HttpResponse.json({ ok: true }, { status: 201 });
  }),

  // Update comment
  http.put(`${API_BASE}/comment`, async ({ request }) => {
    const body = await request.json().catch(() => ({} as any));
    const idx = comments.findIndex(c => c.course_id === body?.course_id && c.user_id === body?.user_id);
    if (idx === -1) return HttpResponse.json({ message: 'not found' }, { status: 404 });
    comments[idx] = { course_id: body.course_id, user_id: body.user_id, text: body.text };
    return HttpResponse.json({ ok: true }, { status: 200 });
  }),
];
