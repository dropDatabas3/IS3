import { authReducer } from '@/context/auth/authReducer';
import type { User } from '@/types';

describe('authReducer', () => {
  const initial = { user: null, token: '' };

  it('returns state by default', () => {
    // @ts-expect-error unknown action
    expect(authReducer(initial as any, { type: 'UNKNOWN' })).toEqual(initial);
  });

  it('handles login action and maps role', () => {
    const user = { id: 'u1', name: 'Juan', email: 'j@e.com', role: 1 } as unknown as User;
    const next = authReducer(initial as any, {
      type: '[Auth] - Login',
      payload: { userData: user, token: 't' }
    });
    expect(next.user?.role).toBe('admin');
    expect(next.token).toBe('t');
  });

  it('handles logout action', () => {
    const logged = { user: { id: 'u', name: 'n', email: 'e', role: 'user' } as any, token: 't' };
    const next = authReducer(logged as any, { type: '[Auth] - Logout' });
    expect(next.user).toBeNull();
    expect(next.token).toBe('');
  });
});
