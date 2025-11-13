import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { AuthContext, CoursesContext } from '@/context';

// Mock next/navigation with a mutable pathname so tests can change it per-case
let mockPathname = '/courses';
jest.mock('next/navigation', () => ({
  usePathname: () => mockPathname,
}));
import { UserButton } from '@/components/ui/navbar/UserButton';

const user = {
  id: 'u1',
  email: 'e@e.com',
  username: 'John Doe',
  avatar: 'https://example.com/a.png',
  role: 'user',
  createdAt: new Date() as any,
  updatedAt: new Date() as any,
};

function renderWithContexts(pathname: string = '/courses', onLogout = jest.fn(), onClean = jest.fn()) {
  // set mocked pathname for this render
  mockPathname = pathname;

  const authValue: any = { user: null, token: '', login: jest.fn(), register: jest.fn(), logout: onLogout };
  const coursesValue: any = {
    courses: [], currentCourse: null, coursesFiltered: [], categories: [], enrollments: [], comments: [], ratings: [],
    createComment: jest.fn(), getComments: jest.fn(), updateComment: jest.fn(),
    createRating: jest.fn(), updateRating: jest.fn(), getRatings: jest.fn(),
    createCourse: jest.fn(), deleteCourse: jest.fn(), updateCourse: jest.fn(),
    filterCourses: jest.fn(), myCourses: jest.fn(), newCategory: jest.fn(), getCategories: jest.fn(),
    enroll: jest.fn(), cleanCourseList: onClean, fetchCourses: jest.fn(), setCurrentCourse: jest.fn(),
  };
  return render(
    <AuthContext.Provider value={authValue}>
      <CoursesContext.Provider value={coursesValue}>
        <UserButton user={user as any} />
      </CoursesContext.Provider>
    </AuthContext.Provider>
  );
}

describe('UserButton', () => {
  it('opens menu on hover and shows My Courses link when not on /my-courses', () => {
    renderWithContexts('/courses');
    const root = screen.getByText(/john/i).closest('div')!;
    fireEvent.mouseEnter(root);
    expect(screen.getByText(/My Profile/i)).toBeInTheDocument();
    expect(screen.getByText(/My Courses/i)).toBeInTheDocument();
  });

  it('hides My Courses link when on /my-courses path', () => {
    renderWithContexts('/my-courses');
    const root = screen.getByText(/john/i).closest('div')!;
    fireEvent.mouseEnter(root);
    expect(screen.queryByText(/My Courses/i)).toBeNull();
  });

  it('calls cleanCourseList and logout on Signout', () => {
    const onLogout = jest.fn();
    const onClean = jest.fn();
    renderWithContexts('/courses', onLogout, onClean);
    const root = screen.getByText(/john/i).closest('div')!;
    fireEvent.mouseEnter(root);
    fireEvent.click(screen.getByText(/Signout/i));
    expect(onClean).toHaveBeenCalled();
    expect(onLogout).toHaveBeenCalled();
  });
});
