import { render, screen } from '@/test/test-utils';
import LoginPage from '@/app/auth/login/page';

describe('LoginPage', () => {
	it('renders heading and form', () => {
		render(<LoginPage />);
		expect(screen.getByRole('heading', { name: /login/i })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
	});
});
