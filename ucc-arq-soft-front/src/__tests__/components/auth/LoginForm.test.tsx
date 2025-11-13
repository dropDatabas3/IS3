import { render, screen, fireEvent, waitFor } from '@/test/test-utils';
// Router mock helpers mapped via jest moduleNameMapper ('^next/navigation$')
// Use star import to avoid TS named export complaints from JS mock
import * as nav from 'next/navigation';
const mockNav: any = nav; // cast for custom mock helpers
import { LoginForm } from '@/components/auth/LoginForm';

describe('LoginForm', () => {
  it('renders email and password inputs and submit button', () => {
    render(<LoginForm />);
    expect(screen.getByPlaceholderText(/email/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
  });

  it('shows validation errors for empty submit', async () => {
    render(<LoginForm />);
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
    await waitFor(() => {
      expect(screen.getAllByText(/required/i).length).toBeGreaterThan(0);
    });
  });

  it('navigates after successful login', async () => {
  mockNav.__reset();
    render(<LoginForm />);
    fireEvent.change(screen.getByPlaceholderText(/email/i), {
      target: { value: 'test@example.com' },
    });
    fireEvent.change(screen.getByPlaceholderText(/password/i), {
      target: { value: '123456' },
    });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
  await waitFor(() => expect(mockNav.__getLastPush()).toBe('/'));
  });
});
