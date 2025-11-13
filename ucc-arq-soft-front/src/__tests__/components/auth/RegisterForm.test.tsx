import { render, screen, fireEvent, waitFor } from '@/test/test-utils';
import * as nav from 'next/navigation';
import { RegisterForm } from '@/components/auth/RegisterForm';

const mockNav: any = nav;

describe('RegisterForm', () => {
	it('renders inputs and submit button', () => {
		render(<RegisterForm />);
		// Component uses placeholder 'Name'
		expect(screen.getByPlaceholderText(/name/i)).toBeInTheDocument();
		expect(screen.getByPlaceholderText(/email/i)).toBeInTheDocument();
		expect(screen.getByPlaceholderText(/password/i)).toBeInTheDocument();
		expect(screen.getByRole('button', { name: /register/i })).toBeInTheDocument();
	});

	it('shows validation errors when submitting empty form', async () => {
		render(<RegisterForm />);
		fireEvent.click(screen.getByRole('button', { name: /register/i }));
		await waitFor(() => {
			const errors = screen.getAllByText(/required/i);
			expect(errors.length).toBeGreaterThan(0);
		});
	});
});
