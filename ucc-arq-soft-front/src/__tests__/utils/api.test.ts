import { createApiClient, apiUrl } from '@/utils/api';

// Mock next/headers to avoid Next runtime dependency in tests
jest.mock('next/headers', () => ({ cookies: jest.fn() }));

describe('api client helpers', () => {
  const realFetch = global.fetch;

  afterEach(() => {
    global.fetch = realFetch as any;
    jest.resetModules();
  });

  it('createApiClient performs GET and handles ok response', async () => {
    global.fetch = jest.fn().mockResolvedValue({ ok: true, json: async () => ({ ok: 1 }) });
    const api = createApiClient();
    const res = await api.get('/health');
    expect(res).toEqual({ ok: 1 });
    expect(global.fetch).toHaveBeenCalledWith(expect.stringMatching(/http/), expect.objectContaining({ method: 'GET' }));
  });

  it('createApiClient throws with server error and message', async () => {
    global.fetch = jest.fn().mockResolvedValue({ ok: false, status: 400, json: async () => ({ message: 'bad' }) });
    const api = createApiClient();
    await expect(api.get('/oops')).rejects.toThrow('bad');
  });

  it('apiUrl returns public url in browser by default', () => {
    process.env.NEXT_PUBLIC_API_URL = 'http://public:8000';
    expect(apiUrl()).toBe('http://public:8000');
  });
});
