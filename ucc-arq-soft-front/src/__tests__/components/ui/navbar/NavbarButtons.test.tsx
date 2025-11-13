import React from 'react';
import { render, screen } from '@/test/test-utils';
import { NavbarButtons } from '@/components/ui/navbar/NavbarButtons';
import { AuthContext } from '@/context';

// Mock UserButton to avoid pulling extra contexts
jest.mock('@/components/ui/navbar/UserButton', () => ({
  UserButton: () => <div data-testid="user-button">User Button</div>,
}));

describe('NavbarButtons', () => {
  it('shows "Get Started" when no user and on home path', () => {
    const authValue = { user: null, token: '', login: jest.fn(), logout: jest.fn(), register: jest.fn() } as any;
    render(
      <AuthContext.Provider value={authValue}>
        <NavbarButtons pathName="/" />
      </AuthContext.Provider>
    );
    expect(screen.getByText(/get started/i)).toBeInTheDocument();
  });

  it('renders UserButton when user is present', () => {
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
    render(
      <AuthContext.Provider value={authValue}>
        <NavbarButtons pathName="/" />
      </AuthContext.Provider>
    );
    expect(screen.getByTestId('user-button')).toBeInTheDocument();
  });
});
