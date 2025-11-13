import React from 'react';
import { act, render, screen, waitFor } from '@/test/test-utils';
import { CoursesContext } from '@/context';

// Helper component to expose context values in DOM for assertions
function Probe() {
  const ctx = React.useContext(CoursesContext);
  return (
    <div>
      <div data-testid="courses-count">{ctx.courses.length}</div>
      <div data-testid="filtered-count">{ctx.coursesFiltered.length}</div>
      <div data-testid="categories-count">{ctx.categories.length}</div>
      <div data-testid="enrollments-count">{ctx.enrollments.length}</div>
      <div data-testid="comments-count">{ctx.comments.length}</div>
      <div data-testid="ratings-count">{ctx.ratings.length}</div>
      <button onClick={() => ctx.getRatings()}>getRatings</button>
      <button onClick={() => ctx.filterCourses('go')}>filterGo</button>
      <button onClick={() => ctx.newCategory('Data')}>newCategory</button>
      <button onClick={() => ctx.createCourse({ course_name: 'X', image: '', description: 'd', price: 1, duration: 1, capacity: 1, category_id: 'cat1', init_date: new Date().toISOString(), state: true })}>createCourse</button>
      <button onClick={() => ctx.enroll('c1')}>enrollC1</button>
      <button onClick={() => ctx.getComments('c1')}>getComments</button>
      <button onClick={() => ctx.createComment('c1', 'u1', 'Nice!')}>createComment</button>
      <button onClick={() => ctx.updateComment('Edited', 'c1', 'u1')}>updateComment</button>
      <button onClick={() => ctx.createRating('c1', 'u3', 3)}>createRating</button>
      <button onClick={() => ctx.updateRating(4, 'c1', 'u3')}>updateRating</button>
    </div>
  );
}

describe('CoursesProvider integration', () => {
  beforeEach(() => {
    // Seed cookie token for endpoints that require auth
    document.cookie = 'token=test-token; path=/';
  });

  it('loads initial data and supports filtering, categories, and course creation', async () => {
    render(<Probe />);

    // fetchCourses and getCategories run on mount
    await waitFor(() => expect(Number(screen.getByTestId('courses-count').textContent)).toBeGreaterThan(0));
    await waitFor(() => expect(Number(screen.getByTestId('categories-count').textContent)).toBeGreaterThan(0));

    // filter courses
    await act(async () => { screen.getByText('filterGo').click(); });
    expect(Number(screen.getByTestId('filtered-count').textContent)).toBeGreaterThan(0);

    // new category
    await act(async () => { screen.getByText('newCategory').click(); });
    await waitFor(() => expect(Number(screen.getByTestId('categories-count').textContent)).toBeGreaterThan(0));

    // create course
    const before = Number(screen.getByTestId('courses-count').textContent);
    await act(async () => { screen.getByText('createCourse').click(); });
    await waitFor(() => expect(Number(screen.getByTestId('courses-count').textContent)).toBeGreaterThan(before));
  });

  it('supports enroll, ratings and comments flows', async () => {
    render(<Probe />);
    // Ensure initial fetch completed
    await waitFor(() => expect(Number(screen.getByTestId('courses-count').textContent)).toBeGreaterThan(0));

    // enroll adds a course into enrollments
    const initialEnrollments = Number(screen.getByTestId('enrollments-count').textContent);
    await act(async () => { screen.getByText('enrollC1').click(); });
    await waitFor(() =>
      expect(Number(screen.getByTestId('enrollments-count').textContent)).toBe(initialEnrollments + 1)
    );

    // ratings
    await act(async () => { screen.getByText('getRatings').click(); });
    await waitFor(() => expect(Number(screen.getByTestId('ratings-count').textContent)).toBeGreaterThan(0));
    await act(async () => { screen.getByText('createRating').click(); });
    await act(async () => { screen.getByText('getRatings').click(); });
    await waitFor(() => expect(Number(screen.getByTestId('ratings-count').textContent)).toBeGreaterThan(0));
    await act(async () => { screen.getByText('updateRating').click(); });

    // comments: 404 path then create & load
    await act(async () => { screen.getByText('getComments').click(); });
    // When 404, provider sets comments to []
    expect(Number(screen.getByTestId('comments-count').textContent)).toBe(0);
    await act(async () => { screen.getByText('createComment').click(); });
    await act(async () => { screen.getByText('getComments').click(); });
    await waitFor(() => expect(Number(screen.getByTestId('comments-count').textContent)).toBe(1));
    await act(async () => { screen.getByText('updateComment').click(); });
  });
});
