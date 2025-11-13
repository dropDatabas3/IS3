import React from 'react';
import { render, screen } from '@/test/test-utils';
import MyCoursesPage from '@/app/my-courses/page';

describe('MyCoursesPage', () => {
  it('renders heading even without enrollments', () => {
    render(<MyCoursesPage />);
    expect(screen.getByRole('heading', { name: /My Courses/i })).toBeInTheDocument();
  });
});
