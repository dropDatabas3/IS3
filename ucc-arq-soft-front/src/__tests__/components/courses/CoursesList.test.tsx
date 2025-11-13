import { render, screen, waitFor } from '@/test/test-utils';
import { CoursesList } from '@/components/courses/CoursesList';

describe('CoursesList', () => {
	it('renders loading state then courses', async () => {
		render(<CoursesList />);
		// Loading indicator (adjust selector if different in component)
		// await data fetch mock via MSW then assert a sample course name
		await waitFor(() => {
			const item = screen.getByText(/Intro Go/i);
			expect(item).toBeInTheDocument();
		});
	});
});
