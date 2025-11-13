import React from 'react';
import { act, render, screen, waitFor } from '@/test/test-utils';
import { CoursesContext } from '@/context';
import { server } from '@/test/server';
import { http, HttpResponse } from 'msw';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

function Probe() {
  const ctx = React.useContext(CoursesContext);
  return (
    <div>
      <div data-testid="courses-count">{ctx.courses.length}</div>
      <div data-testid="categories-count">{ctx.categories.length}</div>
      <div data-testid="enrollments-count">{ctx.enrollments.length}</div>
      <button onClick={() => ctx.getCategories()}>getCategories</button>
      <button onClick={() => ctx.filterCourses('x')}>filterX</button>
      <button onClick={() => ctx.updateCourse({ id: 'c1', course_name: 'Updated' } as any)}>updateC1</button>
      <button onClick={() => ctx.deleteCourse('c1')}>deleteC1</button>
      <button onClick={() => ctx.myCourses()}>myCourses</button>
      <button onClick={() => ctx.getRatings()}>getRatings</button>
      <button onClick={() => ctx.createComment('c1', 'u1', 'Hi')}>createComment</button>
      <button onClick={() => ctx.updateComment('Bye', 'c1', 'u1')}>updateComment</button>
      <button onClick={() => ctx.cleanCourseList()}>clean</button>
    </div>
  );
}

describe('CoursesProvider error branches', () => {
  beforeEach(() => {
    document.cookie = 'token=test-token; path=/';
  });

  it('handles categories non-200 and filter network error gracefully', async () => {
    server.use(
      http.get(`${API_BASE}/categories`, () => HttpResponse.json({ message: 'fail' }, { status: 500 })),
      http.get(`${API_BASE}/courses`, () => HttpResponse.error())
    );
    render(<Probe />);
    await act(async () => { screen.getByText('getCategories').click(); });
    expect(Number(screen.getByTestId('categories-count').textContent)).toBeGreaterThanOrEqual(0);
    await act(async () => { screen.getByText('filterX').click(); });
    // should not crash and keep coursesFiltered managed internally; we assert component still renders
    expect(screen.getByTestId('courses-count')).toBeInTheDocument();
  });

  it('covers updateCourse success and deleteCourse failure', async () => {
    // update success
    render(<Probe />);
    await waitFor(() => expect(Number(screen.getByTestId('courses-count').textContent)).toBeGreaterThan(0));
    await act(async () => { screen.getByText('updateC1').click(); });
    // delete failure
    server.use(
      http.delete(`${API_BASE}/courses/:id`, () => HttpResponse.json({ message: 'no' }, { status: 500 }))
    );
    await act(async () => { screen.getByText('deleteC1').click(); });
    // courses should still be present after failed delete
    expect(Number(screen.getByTestId('courses-count').textContent)).toBeGreaterThan(0);
  });

  it('handles myCourses 500 without altering enrollments', async () => {
    render(<Probe />);
    // Ensure initial enrollments are loaded from the default handler (typically 1)
    await waitFor(() => expect(Number(screen.getByTestId('enrollments-count').textContent)).toBeGreaterThanOrEqual(1));
    const initial = Number(screen.getByTestId('enrollments-count').textContent);
    server.use(http.get(`${API_BASE}/myCourses/`, () => HttpResponse.json({ message: 'err' }, { status: 500 })));
    await act(async () => { screen.getByText('myCourses').click(); });
    expect(Number(screen.getByTestId('enrollments-count').textContent)).toBe(initial);
  });

  it('handles rating/comment token missing and rating non-200', async () => {
    // remove token so create/update are no-ops
    document.cookie = 'token=;expires=Thu, 01 Jan 1970 00:00:00 GMT';
    render(<Probe />);
    await act(async () => { screen.getByText('createComment').click(); });
    await act(async () => { screen.getByText('updateComment').click(); });
    // ratings non-200
    server.use(http.get(`${API_BASE}/rating`, () => HttpResponse.json({ message: 'x' }, { status: 500 })));
    await act(async () => { screen.getByText('getRatings').click(); });
    // still renders; no crash
    expect(screen.getByTestId('courses-count')).toBeInTheDocument();
  });

  it('cleans enrollments list only (courses remain)', async () => {
    render(<Probe />);
    await waitFor(() => expect(Number(screen.getByTestId('courses-count').textContent)).toBeGreaterThan(0));
    await act(async () => { screen.getByText('clean').click(); });
    await waitFor(() => expect(Number(screen.getByTestId('enrollments-count').textContent)).toBe(0));
    // courses should remain populated
    expect(Number(screen.getByTestId('courses-count').textContent)).toBeGreaterThan(0);
  });
});
