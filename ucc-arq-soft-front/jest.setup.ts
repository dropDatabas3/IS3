import '@testing-library/jest-dom';
import 'whatwg-fetch';

// Polyfill TextEncoder/TextDecoder BEFORE importing msw/node
// @ts-ignore
import { TextEncoder, TextDecoder } from 'util';
// @ts-ignore
if (!global.TextEncoder) {
	// @ts-ignore
	global.TextEncoder = TextEncoder;
}
// @ts-ignore
if (!global.TextDecoder) {
	// @ts-ignore
	global.TextDecoder = TextDecoder as any;
}

// Polyfill HTMLFormElement.requestSubmit in jsdom to avoid not-implemented errors
// @ts-ignore
if (typeof HTMLFormElement !== 'undefined' && !HTMLFormElement.prototype.requestSubmit) {
	// @ts-ignore
		HTMLFormElement.prototype.requestSubmit = function requestSubmit(this: HTMLFormElement) {
		// Dispatch a submit event that bubbles and is cancelable, similar to real browsers
		const event = new Event('submit', { bubbles: true, cancelable: true });
			(this as HTMLFormElement).dispatchEvent(event);
	} as any;
}

// Polyfill TransformStream required by msw/interceptors in Node
// Node >=18 exposes it in 'stream/web'
// @ts-ignore
try {
	// eslint-disable-next-line @typescript-eslint/no-var-requires
		const { TransformStream, ReadableStream, WritableStream } = require('stream/web');
	// @ts-ignore
	if (!global.TransformStream && TransformStream) {
		// @ts-ignore
		global.TransformStream = TransformStream;
	}
		// @ts-ignore
		if (!global.ReadableStream && ReadableStream) {
			// @ts-ignore
			global.ReadableStream = ReadableStream;
		}
		// @ts-ignore
		if (!global.WritableStream && WritableStream) {
			// @ts-ignore
			global.WritableStream = WritableStream;
		}
} catch {}

// Polyfill BroadcastChannel used by msw
// @ts-ignore
if (typeof global.BroadcastChannel === 'undefined') {
	// @ts-ignore
	global.BroadcastChannel = class {
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		constructor(_name?: string) {}
		// @ts-ignore
		postMessage(_msg: unknown) {}
		close() {}
		addEventListener() {}
		removeEventListener() {}
		onmessage: any = null;
	} as any;
}

// MSW server for Node test environment
import { server } from './src/test/server';

// Establish API mocking before all tests.
beforeAll(() => server.listen({ onUnhandledRequest: 'warn' }));

// Reset any request handlers that are declared as a part of our tests
// (i.e. for testing one-time error scenarios)
afterEach(() => server.resetHandlers());

// Clean up after the tests are finished.
afterAll(() => server.close());
