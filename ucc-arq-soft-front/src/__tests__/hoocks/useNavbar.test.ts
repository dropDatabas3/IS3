import { renderHook, act } from '@testing-library/react';
import { useNavbar } from '@/utils/hooks/useNavbar';

// Mock next/navigation used by the hook
jest.mock('next/navigation', () => ({
	usePathname: jest.fn(() => '/'),
	useRouter: () => ({ push: jest.fn() })
}));

describe('useNavbar', () => {
	it('sets bgColor based on pathname root and scroll', () => {
		const { result } = renderHook(() => useNavbar());
		// On root path initially bg is transparent
		expect(result.current.bgColor).toBe('bg-transparent');

		// Simulate scrolling past threshold
		Object.defineProperty(window, 'scrollY', { value: 600, writable: true });
		act(() => {
			window.dispatchEvent(new Event('scroll'));
		});
		expect(result.current.bgColor).toBe('bg-navbar');
	});

	it('sets bgColor to bg-navbar for auth pages', () => {
		// Override pathname mock to include auth
		const usePathname = require('next/navigation').usePathname as jest.Mock;
		usePathname.mockReturnValue('/auth/login');

		const { result } = renderHook(() => useNavbar());
		expect(result.current.bgColor).toBe('bg-navbar');
	});
});
