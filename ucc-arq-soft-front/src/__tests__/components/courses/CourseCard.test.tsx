import { render, screen } from '@/test/test-utils';
import { CourseCard } from '@/components/courses/card';

describe('CourseCard', () => {
		const course = {
			id: 'c1',
			courseName: 'Intro Go',
			description: 'Basics of Go',
			price: 10,
			duration: 5,
			capacity: 30,
			categoryId: 'cat1',
			initDate: new Date().toISOString(),
			state: true,
			image: 'go.png',
			ratingAvg: 4.5,
			categoryName: 'Programming',
		} as any;

	it('renders course name and price', () => {
		render(<CourseCard course={course} onClick={() => {}} />);
			expect(screen.getByText(/Intro Go/i)).toBeInTheDocument();
			expect(screen.getByText(/\$10/i)).toBeInTheDocument();
	});
});
