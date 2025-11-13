import { courseMapper } from '@/utils/mappers/courseMapper';

describe('courseMapper', () => {
  it('maps API course to domain Course with defaults', () => {
    const apiCourse = {
      id: 'c1',
      course_name: 'Intro',
      description: 'desc',
      image: 'img.png',
      category_name: 'Cat',
    } as any;

    const mapped = courseMapper(apiCourse);

    expect(mapped).toEqual(
      expect.objectContaining({
        id: 'c1',
        courseName: 'Intro',
        description: 'desc',
        image: 'img.png',
        categoryName: 'Cat',
        price: 0,
        duration: 0,
        capacity: 15,
        state: true,
      })
    );
  });
});
