import React from 'react';
import { render, screen, waitFor } from '@/test/test-utils';
import { AuthContext } from '@/context';
import { server } from '@/test/server';
import { http, HttpResponse } from 'msw';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8000';

function Probe() {
  const ctx = React.useContext(AuthContext);
  return (
    <div>
      <div data-testid="has-user">{ctx.user ? 'yes' : 'no'}</div>
      <div data-testid="token">{ctx.token ?? ''}</div>
      <button onClick={() => ctx.login({ Email: 'test@example.com', Password: 'pw' })}>login</button>
      <button onClick={() => ctx.register({ Username: 'User', Email: 'u@e.com', Password: 'pw' })}>register</button>
      <button onClick={() => ctx.register({ Username: '', Email: '', Password: '' } as any)}>registerBad</button>
      <button onClick={() => ctx.login({ Email: '', Password: '' } as any)}>loginBad</button>
    </div>
  );
}

describe('AuthProvider', () => {
  beforeEach(() => {
    document.cookie = '';
  });

  it('refreshes token on mount (success) and sets user', async () => {
    // seed cookie to trigger refresh
    document.cookie = 'token=seed';
    render(<Probe />);
    await waitFor(() => expect(screen.getByTestId('has-user').textContent).toBe('yes'));
    expect(screen.getByTestId('token').textContent).not.toBe('');
  });

  it('refresh token failure logs out', async () => {
    document.cookie = 'token=seed';
    server.use(
      http.post(`${API_BASE}/auth/refresh-token`, () => HttpResponse.json({ message: 'nope' }, { status: 401 }))
    );
    render(<Probe />);
    await waitFor(() => expect(screen.getByTestId('has-user').textContent).toBe('no'));
  });

  it('register success and validation failure', async () => {
    render(<Probe />);
    // missing fields returns false and keeps user null
    await screen.getByText('registerBad').click();
    expect(screen.getByTestId('has-user').textContent).toBe('no');
    // success
    await screen.getByText('register').click();
    await waitFor(() => expect(screen.getByTestId('has-user').textContent).toBe('yes'));
  });

  it('login success and validation failure', async () => {
    render(<Probe />);
    await screen.getByText('loginBad').click();
    expect(screen.getByTestId('has-user').textContent).toBe('no');
    await screen.getByText('login').click();
    await waitFor(() => expect(screen.getByTestId('has-user').textContent).toBe('yes'));
  });
});
