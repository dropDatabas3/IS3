import { coursesReducer } from '@/context/courses/contextReducer';
import type { CoursesState } from '@/context/courses/CoursesProvider';

const baseState: CoursesState = {
  courses: [],
  currentCourse: null,
  coursesFiltered: [],
  categories: [],
  enrollments: [],
  comments: [],
  ratings: [],
};

describe('coursesReducer', () => {
  it('loads all courses and filters mirror', () => {
    const state = coursesReducer(baseState, { type: '[Courses] - Load All', payload: [{ id: 'c1' } as any] });
    expect(state.courses).toHaveLength(1);
    expect(state.coursesFiltered).toHaveLength(1);
  });

  it('sets current course', () => {
    const course = { id: 'c2', courseName: 'X' } as any;
    const state = coursesReducer(baseState, { type: '[Courses] - Set Current', payload: course });
    expect(state.currentCourse).toEqual(course);
  });

  it('enroll adds to enrollments', () => {
    const withCourses = { ...baseState, courses: [{ id: 'c3' } as any] } as CoursesState;
    const state = coursesReducer(withCourses, { type: '[Courses] - Enroll', payload: 'c3' });
    expect(state.enrollments).toHaveLength(1);
    expect(state.enrollments[0]?.id).toBe('c3');
  });

  it('comments and ratings load', () => {
    const s1 = coursesReducer(baseState, { type: '[Comments] - Load All Comments', payload: [{ id: 'cm1' } as any] });
    expect(s1.comments).toHaveLength(1);
    const s2 = coursesReducer(baseState, { type: '[Ratings] - Load All Ratings', payload: [{ id: 'r1' } as any] });
    expect(s2.ratings).toHaveLength(1);
  });

  it('filter updates coursesFiltered only', () => {
    const withData = { ...baseState, courses: [{ id: 'c1' } as any] } as CoursesState;
    const st = coursesReducer(withData, { type: '[Courses] - Filter', payload: [] });
    expect(st.courses).toHaveLength(1);
    expect(st.coursesFiltered).toHaveLength(0);
  });
});
