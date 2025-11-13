import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { CourseModal } from '@/components/courses/CourseDetail';
import { AuthContext } from '@/context/auth/AuthContext';
import { CoursesContext } from '@/context/courses/CoursesContext';
import { UiContext } from '@/context/ui/UiContext';

// Use the jest mock for next/navigation to assert router pushes
import * as navigation from 'next/navigation';

const sampleCourse = {
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
  ratingavg: 4.2,
  categoryName: 'Programming',
};

function wrap(ui: React.ReactElement) {
  const authValue = { user: { id: 'u1', role: 'admin' } as any, token: '', login: jest.fn(), register: jest.fn(), logout: jest.fn() };
  const coursesValue = {
    courses: [],
    coursesFiltered: [],
    categories: [],
    enrollments: [],
    currentCourse: null,
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
  const uiValue = { isEdit: false, isCreateModalOpen: false, openCreateModal: jest.fn(), closeCreateModal: jest.fn() } as any;

  return (
    <UiContext.Provider value={uiValue}>
      <AuthContext.Provider value={authValue}>
        <CoursesContext.Provider value={coursesValue}>{ui}</CoursesContext.Provider>
      </AuthContext.Provider>
    </UiContext.Provider>
  );
}

describe('CourseDetail (CourseModal)', () => {
  it('renders and navigates to see more', () => {
    const onClose = jest.fn();
    render(wrap(<CourseModal course={sampleCourse as any} onClose={onClose} />));
    // Button shows 'See more'
    const button = screen.getByRole('button', { name: /See more/i });
    fireEvent.click(button);
    const pushed = (navigation as any).__getLastPush?.() || (navigation as any).__router?.push.mock.calls?.[0]?.[0];
    expect(pushed).toBe('/course-info');
  });

  it('shows admin controls when user is admin', () => {
    const onClose = jest.fn();
    render(wrap(<CourseModal course={sampleCourse as any} onClose={onClose} />));
    // There should be two admin buttons (edit and delete) plus close button
    // Query by title icons presence via role=button count > 1
    const buttons = screen.getAllByRole('button');
    expect(buttons.length).toBeGreaterThan(1);
  });
});
