import { commentMapper } from '@/utils/mappers/commentMapper';

describe('commentMapper', () => {
	it('maps raw comment to typed Comment', () => {
		const raw = { comment: 'Nice', user_name: 'Alice', user_avatar: 'a.png', user_id: 'u1', extra: 'ignored' };
		const mapped = commentMapper(raw);
		expect(mapped).toEqual({ comment: 'Nice', user_name: 'Alice', user_avatar: 'a.png', user_id: 'u1' });
	});

	it('handles missing optional fields by passing through undefined', () => {
		const raw = { comment: 'Hello', user_name: 'Bob' } as any;
		const mapped = commentMapper(raw);
		expect(mapped.user_avatar).toBeUndefined();
		expect(mapped.user_id).toBeUndefined();
	});
});
