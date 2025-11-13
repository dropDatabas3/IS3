import React from 'react';
import { render, screen, fireEvent } from '@/test/test-utils';
import { NavBarMainButtons } from '@/components/ui/navbar/NavBarButtonsMain';
import { AuthContext, CoursesContext, UiContext } from '@/context';

function withContexts(ui: React.ReactElement, { userRole = 'admin', filterCourses = jest.fn(), openCreateModal = jest.fn() } = {}) {
  const authValue = {
    user: {
      id: 'u1',
      email: 'test@example.com',
      username: 'Tester',
      avatar: 'https://example.com/a.png',
      role: userRole,
      createdAt: new Date(),
      updatedAt: new Date(),
    },
    token: 't',
    login: jest.fn(),
    logout: jest.fn(),
    register: jest.fn(),
  } as any;
  const coursesValue = { filterCourses } as any;
  const uiValue = { openCreateModal } as any;
  return render(
    <UiContext.Provider value={uiValue}>
      <AuthContext.Provider value={authValue}>
        <CoursesContext.Provider value={coursesValue}>{ui}</CoursesContext.Provider>
      </AuthContext.Provider>
    </UiContext.Provider>
  );
}

describe('NavBarMainButtons', () => {
  it('calls filterCourses on enter key in search input', () => {
    const filterCourses = jest.fn();
    withContexts(
      <NavBarMainButtons
        pathname="/courses"
        handleNavigateAbout={() => {}}
        handleNavigateCourses={() => {}}
        handleNavigateHome={() => {}}
      />, { filterCourses }
    );
    const input = screen.getByPlaceholderText(/buscar cursos/i);
    fireEvent.change(input, { target: { value: 'go' } });
    fireEvent.keyDown(input, { key: 'Enter', code: 'Enter', charCode: 13 });
    expect(filterCourses).toHaveBeenCalledWith('go');
  });

  it('shows "Create course" button for admin on /courses and calls openCreateModal', () => {
    const openCreateModal = jest.fn();
    withContexts(
      <NavBarMainButtons
        pathname="/courses"
        handleNavigateAbout={() => {}}
        handleNavigateCourses={() => {}}
        handleNavigateHome={() => {}}
      />, { openCreateModal, userRole: 'admin' }
    );
    const createBtn = screen.getByRole('button', { name: /create course/i });
    fireEvent.click(createBtn);
    expect(openCreateModal).toHaveBeenCalled();
  });

  it('does not show "Create course" button for non-admin user', () => {
    withContexts(
      <NavBarMainButtons
        pathname="/courses"
        handleNavigateAbout={() => {}}
        handleNavigateCourses={() => {}}
        handleNavigateHome={() => {}}
      />, { userRole: 'user' }
    );
    expect(screen.queryByRole('button', { name: /create course/i })).not.toBeInTheDocument();
  });
});
