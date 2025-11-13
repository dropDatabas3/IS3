import React from 'react';
import { render, screen, fireEvent } from '@/test/test-utils';
import { CourseCarrusel } from '@/components/courses/CourseCarrusel';

const course = {
  id: 'c1',
  courseName: 'Intro Go',
  description: 'A'.repeat(100),
  price: 10,
  duration: 5,
  capacity: 30,
  categoryId: 'cat1',
  initDate: new Date().toISOString(),
  state: true,
  image: 'go.png',
  ratingAvg: 4.5,
  categoryName: 'Programming',
} as any;

describe('CourseCarrusel', () => {
  it('renders empty state message when no courses', () => {
    render(<CourseCarrusel courses={[]} handlerSelected={jest.fn()} />);
    expect(screen.getByText(/No hay cursos disponibles/i)).toBeInTheDocument();
  });

  it('clicking a slide calls handlerSelected and truncates description', () => {
    const handler = jest.fn();
    render(<CourseCarrusel courses={[course]} handlerSelected={handler} />);
    fireEvent.click(screen.getByText('Intro Go'));
    expect(handler).toHaveBeenCalledWith(expect.objectContaining({ id: 'c1' }));
    // truncated to 60 + ...
    expect(screen.getByText(/\.\.\.$/).textContent?.length).toBeGreaterThan(3);
  });
});
