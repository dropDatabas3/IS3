import React from 'react';
import { render, screen, waitFor } from '@/test/test-utils';
import CoursesPage from '@/app/courses/page';

describe('CoursesPage', () => {
  it('renders page sections and headings', async () => {
    render(<CoursesPage />);
  expect(screen.getByRole('heading', { name: /Top Courses/i })).toBeInTheDocument();
  // There are multiple headings containing "Courses", ensure at least one matches
  expect(screen.getAllByRole('heading', { name: /Courses/i }).length).toBeGreaterThan(0);

    // After MSW loads courses, the carousel/list should include sample course
    await waitFor(() => {
      expect(screen.getAllByText(/Intro Go/i).length).toBeGreaterThan(0);
    });
  });
});
