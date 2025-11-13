import React from 'react';
import { fireEvent, render, screen } from '@/test/test-utils';
import { CreateCourseModal } from '@/components/courses/NewCourseForm';
import { AuthContext, CoursesContext, UiContext } from '@/context';

describe('CreateCourseModal', () => {
  function renderWithContexts(onClose = jest.fn(), createCourse = jest.fn()) {
    const authValue = {
      user: {
        id: 'u1',
        email: 'test@example.com',
        username: 'Tester',
        avatar: 'https://example.com/a.png',
        role: 'admin',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      token: 't',
      login: jest.fn(),
      logout: jest.fn(),
      register: jest.fn(),
    } as any;
    const coursesValue = {
      categories: [{ category_id: 'cat1', category_name: 'Programming' }],
      currentCourse: null,
      newCategory: jest.fn(),
      createCourse,
      updateCourse: jest.fn(),
      courses: [],
      coursesFiltered: [],
      enrollments: [],
      comments: [],
      ratings: [],
      createComment: jest.fn(),
      getComments: jest.fn(),
      updateComment: jest.fn(),
      getRatings: jest.fn(),
      createRating: jest.fn(),
      updateRating: jest.fn(),
      deleteCourse: jest.fn(),
      filterCourses: jest.fn(),
      myCourses: jest.fn(),
      setCurrentCourse: jest.fn(),
      getCategories: jest.fn(),
      enroll: jest.fn(),
      cleanCourseList: jest.fn(),
      fetchCourses: jest.fn(),
    } as any;
    const uiValue = {
      isEdit: false,
      isCreateModalOpen: true,
      openCreateModal: jest.fn(),
      closeCreateModal: jest.fn(),
    } as any;
    return render(
      <UiContext.Provider value={uiValue}>
        <AuthContext.Provider value={authValue}>
          <CoursesContext.Provider value={coursesValue}>
            <CreateCourseModal onClose={onClose} />
          </CoursesContext.Provider>
        </AuthContext.Provider>
      </UiContext.Provider>
    );
  }

  it('fills minimal fields and submits createCourse', () => {
    const onClose = jest.fn();
    const createCourse = jest.fn();
    renderWithContexts(onClose, createCourse);

    // Fill description
    fireEvent.change(screen.getByPlaceholderText(/course description/i), {
      target: { value: 'A great course' },
    });
    // Fill basic inputs
    fireEvent.change(screen.getByPlaceholderText(/course name/i), { target: { value: 'Intro Go' } });
    fireEvent.change(screen.getByPlaceholderText(/price/i), { target: { value: '25' } });
    fireEvent.change(screen.getByPlaceholderText(/duration/i), { target: { value: '5' } });
    fireEvent.change(screen.getByPlaceholderText(/capacity/i), { target: { value: '30' } });
    // Select category
    fireEvent.change(screen.getByRole('combobox'), { target: { value: 'cat1' } });
    // Submit
    fireEvent.click(screen.getByRole('button', { name: /create course/i }));

    expect(createCourse).toHaveBeenCalled();
    expect(onClose).toHaveBeenCalled();
  });
});
