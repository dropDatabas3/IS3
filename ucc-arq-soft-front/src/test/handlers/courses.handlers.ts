import { http, HttpResponse, delay } from 'msw';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

let sampleCourses = [
  {
    id: 'c1',
    course_name: 'Intro Go',
    description: 'Basics of Go',
    price: 10,
    duration: 5,
    capacity: 30,
    category_id: 'cat1',
    init_date: new Date().toISOString(),
    state: true,
    image: 'go.png',
    ratingavg: 4.5,
    category_name: 'Programming',
  },
];

export const handlers = [
  http.get(`${API_BASE}/courses`, async ({ request }) => {
    const url = new URL(request.url);
    const filter = url.searchParams.get('filter');
    await delay(50);
    if (filter) {
      const filtered = sampleCourses.filter((c) =>
        c.course_name.toLowerCase().includes(filter.toLowerCase())
      );
      return HttpResponse.json(filtered, { status: 200 });
    }
    return HttpResponse.json(sampleCourses, { status: 200 });
  }),

  http.get(`${API_BASE}/categories`, () =>
    HttpResponse.json([{ id: 'cat1', category_name: 'Programming' }], {
      status: 200,
    })
  ),

  // Create course
  http.post(`${API_BASE}/courses/create`, async ({ request }) => {
    const body = await request.json().catch(() => ({} as any));
    const created = {
      id: `c${sampleCourses.length + 1}`,
      course_name: body.course_name ?? 'New',
      description: body.description ?? '',
      price: body.price ?? 0,
      duration: body.duration ?? 0,
      capacity: body.capacity ?? 0,
      category_id: body.category_id ?? 'cat1',
      init_date: body.init_date ?? new Date().toISOString(),
      state: true,
      image: body.image ?? '',
      ratingavg: 0,
      category_name: 'Programming',
    };
    sampleCourses.push(created as any);
    return HttpResponse.json(created, { status: 201 });
  }),

  // Update course
  http.put(`${API_BASE}/courses/update/:id`, async ({ params, request }) => {
    const id = params.id as string;
    const body = await request.json().catch(() => ({} as any));
    const idx = sampleCourses.findIndex(c => c.id === id);
    if (idx === -1) return HttpResponse.json({ message: 'not found' }, { status: 404 });
    sampleCourses[idx] = { ...sampleCourses[idx], ...body } as any;
    return HttpResponse.json(sampleCourses[idx], { status: 200 });
  }),

  // Delete course
  http.delete(`${API_BASE}/courses/:id`, ({ params }) => {
    const id = params.id as string;
    sampleCourses = sampleCourses.filter(c => c.id !== id);
    return HttpResponse.json({ ok: true }, { status: 200 });
  }),

  // My courses
  http.get(`${API_BASE}/myCourses/`, () => {
    // Return first course as enrolled for tests
    return HttpResponse.json([sampleCourses[0]], { status: 200 });
  }),

  // New category
  http.post(`${API_BASE}/category/create`, async ({ request }) => {
    const body = await request.json().catch(() => ({} as any));
    return HttpResponse.json({ category_id: 'cat2', category_name: body.category_name ?? 'NewCat' }, { status: 201 });
  }),

  // Enroll
  http.post(`${API_BASE}/enroll`, () => HttpResponse.json({ ok: true }, { status: 201 })),
];
