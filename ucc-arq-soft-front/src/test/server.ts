import { setupServer } from 'msw/node';
import { handlers as authHandlers } from './handlers/auth.handlers';
import { handlers as coursesHandlers } from './handlers/courses.handlers';
import { handlers as ratingHandlers } from './handlers/rating.handlers';
import { handlers as commentsHandlers } from './handlers/comments.handlers';

export const server = setupServer(
	...authHandlers,
	...coursesHandlers,
	...ratingHandlers,
	...commentsHandlers,
);
