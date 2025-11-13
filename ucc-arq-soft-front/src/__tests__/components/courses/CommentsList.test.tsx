import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { CommentsList } from '@/components/courses/CommentsLists';
import { AuthContext } from '@/context/auth/AuthContext';
import { CoursesContext } from '@/context/courses/CoursesContext';

function wrap(ui: React.ReactElement, ctx?: Partial<React.ComponentProps<typeof CoursesContext.Provider>['value']>) {
  const authValue = { user: { id: 'u1', role: 'user' } as any, token: '', login: jest.fn(), register: jest.fn(), logout: jest.fn() };
  const baseCourses = {
    courses: [],
    coursesFiltered: [],
    categories: [],
    enrollments: [{ id: 'c1' }] as any,
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
  const value = { ...baseCourses, ...(ctx || {}) };
  return (
    <AuthContext.Provider value={authValue}>
      <CoursesContext.Provider value={value}>{ui}</CoursesContext.Provider>
    </AuthContext.Provider>
  );
}

describe('CommentsList', () => {
  it('shows empty state then allows submitting a comment when enrolled', () => {
    const createComment = jest.fn();
    render(
      wrap(<CommentsList courseId="c1" />, {
        enrollments: [{ id: 'c1' }] as any,
        comments: [],
        ratings: [],
        createComment,
        getComments: jest.fn(),
      })
    );
    expect(screen.getByText(/No comments yet/i)).toBeInTheDocument();
    const textarea = screen.getByPlaceholderText(/Write your comment/i);
    fireEvent.change(textarea, { target: { value: 'Nice course' } });
    fireEvent.click(screen.getByRole('button', { name: /Submit Comment/i }));
    expect(createComment).toHaveBeenCalledWith('c1', 'u1', 'Nice course');
  });

  it('renders existing comments with user label', () => {
    render(
      wrap(<CommentsList courseId="c1" />, {
        enrollments: [{ id: 'c1' }] as any,
        comments: [{ user_id: 'u1', user_name: 'Alice', user_avatar: 'a.png', comment: 'Hello' }],
        ratings: [],
        getComments: jest.fn(),
      })
    );
    expect(screen.getByText(/Tu/i)).toBeInTheDocument();
    expect(screen.getAllByText(/Hello/i).length).toBeGreaterThan(0);
  });
});
