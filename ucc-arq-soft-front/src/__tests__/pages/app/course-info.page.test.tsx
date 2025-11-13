import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import CoursePage from '@/app/course-info/page';
import { AuthContext } from '@/context/auth/AuthContext';
import { CoursesContext } from '@/context/courses/CoursesContext';

function wrap(ui: React.ReactElement, coursesCtx: Partial<React.ComponentProps<typeof CoursesContext.Provider>['value']> = {}) {
  const authValue = { user: { id: 'u1', role: 'user' } as any, token: '', login: jest.fn(), register: jest.fn(), logout: jest.fn() };
  const baseCourses = {
    courses: [],
    coursesFiltered: [],
    categories: [],
    enrollments: [{ id: 'c1' }] as any,
    currentCourse: {
      id: 'c1',
      courseName: 'Intro Go',
      description: 'Basics of Go',
      price: 10,
      duration: 5,
      capacity: 30,
      category_id: 'cat1',
      initDate: new Date().toISOString(),
      state: true,
      image: 'go.png',
      ratingavg: 4,
      categoryName: 'Programming',
    },
    comments: [],
    ratings: [],
    createComment: jest.fn(),
    getComments: jest.fn(),
    updateComment: jest.fn(),
    getRatings: jest.fn(),
    createRating: jest.fn(),
    updateRating: jest.fn(),
    createCourse: jest.fn(),
    deleteCourse: jest.fn(),
    updateCourse: jest.fn(),
    filterCourses: jest.fn(),
    myCourses: jest.fn(),
    newCategory: jest.fn(),
    getCategories: jest.fn(),
    enroll: jest.fn(),
    cleanCourseList: jest.fn(),
    fetchCourses: jest.fn(),
    setCurrentCourse: jest.fn(),
  } as any;
  const value = { ...baseCourses, ...coursesCtx };
  return (
    <AuthContext.Provider value={authValue}>
      <CoursesContext.Provider value={value}>{ui}</CoursesContext.Provider>
    </AuthContext.Provider>
  );
}

describe('Course Info Page', () => {
  it('shows rating controls for enrolled user and triggers update/create', () => {
    const updateRating = jest.fn();
    const createRating = jest.fn();
    render(wrap(<CoursePage />, { updateRating, createRating, ratings: [{ user_id: 'u1', course_id: 'c1', rating: 3 }] as any }));

  // Change to star4 by selecting input with value="4"
    const star4 = screen.getByDisplayValue('4');
    fireEvent.click(star4);
    expect(updateRating).toHaveBeenCalledWith(4, 'c1', 'u1');

    // If user had no rating, createRating should be called
  const { rerender } = render(wrap(<CoursePage />, { createRating, ratings: [] }));
  const stars5 = screen.getAllByDisplayValue('5');
  const star5 = stars5[stars5.length - 1];
    fireEvent.click(star5);
    expect(createRating).toHaveBeenCalledWith('c1', 'u1', 5);
    rerender(null as any);
  });
});
